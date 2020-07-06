package job

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	. "github.com/anhntbk08/gateway/internal/app/tmbridgev2/bus"
	store "github.com/anhntbk08/gateway/internal/app/tmbridgev2/store"
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
	appCommon "github.com/anhntbk08/gateway/internal/common"
	"github.com/anhntbk08/machinery/v1/log"
	"github.com/globalsign/mgo/bson"
	ethereum "github.com/tomochain/tomochain"
	"github.com/tomochain/tomochain/common"
	ethtypes "github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/ethclient"
)

type JobServer struct {
	bus         *Bus
	mutex       sync.Mutex
	db          *store.Mongo
	chainConfig appCommon.ChainConfig
}

func StartServer(bus *Bus, db *store.Mongo, chainConfig appCommon.ChainConfig) error {
	jobServer := &JobServer{
		bus:         bus,
		db:          db,
		chainConfig: chainConfig,
	}
	busServer := jobServer.bus.GetBusServer()

	tasks := map[string]interface{}{
		"notify_project_sc_changed":      jobServer.UpdateSmartContractWatchList,
		"sync_smartcontract_transaction": jobServer.SyncSmartContractTransaction,
	}
	err := busServer.RegisterTasks(tasks)

	go func() {
		jobWorker := busServer.NewWorker("job_worker", 100)
		err = jobWorker.Launch()
		if err != nil {
			log.ERROR.Println("Can't launch job worker: ", err)
		}
	}()

	return err
}

// Any way to validate this address
func (js *JobServer) UpdateSmartContractWatchList(projectID string, removedAddress []string, newAddresses []string) error {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	err := js.db.SmartContractDao.BulkRemove(removedAddress)
	if err != nil {
		// TODO project error_log to instructment metrics
		fmt.Println("err remove addresses ", err)
		return err
	}

	// create separate tasks for syncing transactions each address
	for i := 0; i < len(newAddresses); i++ {
		js.bus.CreateSyncingSmartContractTransaction(projectID, newAddresses[i])
	}

	return nil
}

// Any way to validate this address
func (js *JobServer) SyncSmartContractTransaction(projectID, address string) error {
	// Check if is syncing
	js.mutex.Lock()

	isSyncing := js.db.SmartContractDao.IsSyncing(address)
	if isSyncing {
		return nil
	}

	err := js.db.SmartContractDao.StartSync(address)
	js.mutex.Unlock()

	return err
}

func (js *JobServer) ReadSmartContractTx(projectID, address string) error {
	err := js.scanLogs(address)
	return err
}

func (js *JobServer) scanLogs(address string) error {
	client, err := ethclient.Dial(js.chainConfig.RPC)

	if err != nil {
		log.FATAL.Println(err)
		return err
	}

	// check current scanned to block's logs
	var sc entity.SmartContract
	scannedTo := int64(0)
	err = js.db.SmartContractDao.GetOne(bson.M{
		"address": address,
	}, &sc)

	if err == nil {
		scannedTo = sc.ScannedIndex
	}

	// TODO, break query to 10k block each turn
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(scannedTo)),
		Addresses: []common.Address{
			common.HexToAddress(address),
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)

	err = js.saveLogs(logs, client, scannedTo, sc.ID)

	return err
}

func (js *JobServer) saveLogs(logs []ethtypes.Log, client *ethclient.Client, scannedTo int64, scID bson.ObjectId) error {
	execution := []byte("Transfer(address,address,uint256)")
	executionHash := crypto.Keccak256Hash(execution)
	transactions := []entity.SmartContractTransaction{}

	for _, scLog := range logs {
		switch scLog.Topics[0].Hex() {
		case executionHash.Hex():
			transaction := entity.SmartContractTransaction{}
			blockInfo, err := client.BlockByHash(context.Background(), scLog.BlockHash)

			if err == nil {
				transaction.Timestamp = blockInfo.ReceivedAt
				transaction.Nonce = blockInfo.Nonce()
				transaction.Gas = blockInfo.GasLimit()
				transaction.CumulativeGasUsed = blockInfo.GasUsed()
			}

			transaction.From = scLog.Topics[1].Hex()
			transaction.To = scLog.Topics[2].Hex()
			transaction.BlockHash = scLog.BlockHash.Hex()
			transaction.BlockNumber = scLog.BlockNumber
			transaction.SmartContract = scID
			transaction.Hash = scLog.TxHash.Hex()
			transaction.TxIndex = scLog.TxIndex
			transaction.Value = new(big.Int).SetBytes(scLog.Data)

			transactions = append(transactions, transaction)
		}
	}

	err := js.db.SmartContractDao.StopSync(scID, scannedTo)

	log.ERROR.Println("Finalizing sync for ", scID, scannedTo, " got err ", err)

	// go tomo.db.SmartContractTxDao.SaveNewTxs(transactions, tomo.CoinType)
	return nil
}
