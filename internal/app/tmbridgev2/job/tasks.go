package job

import (
	"context"
	"math/big"
	"strings"
	"sync"
	"time"

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

	go jobServer.WatchTx()

	return err
}

// Any way to validate this address
func (js *JobServer) UpdateSmartContractWatchList(projectID string, removedAddress []string, newAddresses []string) error {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	// actually we don't need to delete the old address
	// err := js.db.SmartContractDao.BulkRemove(removedAddress)
	// if err != nil {
	// 	fmt.Println("err remove addresses ", err)
	// 	return err
	// }

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

	js.ReadSmartContractTx(projectID, address)
	return err
}

func (js *JobServer) ReadSmartContractTx(projectID, address string) error {
	err := js.scanLogs(address, 0)
	return err
}

func (js *JobServer) scanLogs(address string, toBlock int64) error {
	client, err := ethclient.Dial(js.chainConfig.RPC)

	if err != nil {
		log.FATAL.Println(err)
		return err
	}

	// check current scanned to block's logs
	sc := entity.SmartContract{}
	scannedTo := int64(0)
	err = js.db.SmartContractDao.GetOne(bson.M{
		"address": strings.ToLower(address),
	}, &sc)

	if err == nil {
		scannedTo = sc.ScannedIndex
	}

	// TODO, break query to 10k block each turn
	var query ethereum.FilterQuery

	if toBlock > scannedTo {
		query = ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(scannedTo)),
			ToBlock:   big.NewInt(toBlock),
			Addresses: []common.Address{
				common.HexToAddress(address),
			},
		}
	} else {
		query = ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(scannedTo)),
			Addresses: []common.Address{
				common.HexToAddress(address),
			},
		}
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
			transaction.ID = bson.NewObjectId()
			blockInfo, err := client.BlockByHash(context.Background(), scLog.BlockHash)

			if err == nil {
				transaction.Timestamp = blockInfo.ReceivedAt
				transaction.Nonce = blockInfo.Nonce()
				transaction.Gas = blockInfo.GasLimit()
				transaction.CumulativeGasUsed = blockInfo.GasUsed()
			}

			transaction.From = common.HexToAddress(scLog.Topics[1].Hex()).String()
			transaction.To = common.HexToAddress(scLog.Topics[2].Hex()).String()
			transaction.BlockHash = scLog.BlockHash.Hex()
			transaction.BlockNumber = scLog.BlockNumber
			transaction.SmartContract = scID
			transaction.Hash = scLog.TxHash.Hex()
			transaction.TxIndex = scLog.TxIndex
			transaction.Value = new(big.Int).SetBytes(scLog.Data).String()
			transactions = append(transactions, transaction)

			if scannedTo < int64(transaction.BlockNumber) {
				scannedTo = int64(transaction.BlockNumber)
			}
		}
	}

	err := js.db.SmartContractDao.StopSync(scID, scannedTo)

	log.INFO.Println("Finish syncing  ", scID, scannedTo, " --  ", err)

	go js.db.SmartContractTxDao.InsertBulk(transactions)

	return nil
}

func (js *JobServer) WatchTx() error {
	for {
		ethclient, err := ethclient.Dial(js.chainConfig.RPC)
		if err != nil {
			time.Sleep(js.chainConfig.IntervalRunningTime)
			continue
		}
		currentBlock := js.db.ScannedIndexDao.GetCurrentBlock("TOMO")

		log.INFO.Println("[TOMO] Watch TX of block ", currentBlock)
		if currentBlock == 0 {
			currentBlock = js.chainConfig.StartBlock
		}

		block, err := ethclient.BlockByNumber(context.Background(), big.NewInt(int64(currentBlock)))

		if err == nil {
			js.findSmartContractTransactions(ethclient, block, currentBlock)
		} else {
			time.Sleep(js.chainConfig.IntervalRunningTime)
			continue
		}

		js.db.ScannedIndexDao.SetCurrentBlock("TOMO", currentBlock+1)
		time.Sleep(js.chainConfig.IntervalRunningTime)
	}
}

// TODO check performance
func (js *JobServer) findSmartContractTransactions(ethclient *ethclient.Client, block *ethtypes.Block, currentBlock uint64) {
	txs := block.Transactions()
	scannedSc := map[string]bool{}
	for _, tx := range txs {
		if tx.To() == nil || scannedSc[tx.To().Hex()] == true {
			continue
		}

		sc := js.db.SmartContractDao.GetByAddress(tx.To().Hex())
		if sc != nil {
			scannedSc[tx.To().Hex()] = true
			go js.scanLogs(sc.Address, int64(currentBlock))
		}
	}

}
