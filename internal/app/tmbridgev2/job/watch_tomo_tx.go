package job

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"math/big"
// 	"time"

// 	"github.com/anhntbk08/machinery/v1/log"
// 	"github.com/tomochain/tomochain/common"
// 	"github.com/tomochain/tomochain/ethclient"
// )

// // ETH network
// type ETH struct {
// 	*Task
// }

// // WatchTx imp eth task interface
// // should run interval, after end of this one create an other

// func (eth *ETH) delayBeforeRetrying(blockNumber uint64) {
// 	eth.DB.SetCurrentBlock("ETH", blockNumber)
// 	time.Sleep(5 * time.Second)
// }

// // ConfirmTx verify tx
// func (eth *ETH) ConfirmTx(txtype string, to string, txData types.Transaction) bool {
// 	ethclient, err := ethclient.Dial(eth.FullNodeConfig.RPC)
// 	if err != nil {
// 		return false
// 	}
// 	_, isPending, err := ethclient.TransactionByHash(context.Background(), common.HexToHash(txData.Hash))
// 	if err != nil || isPending {
// 		return false
// 	}
// 	// Not needed to verify to != tx.To().Hex() because in internal-transaction, it would not show up
// 	// if to != tx.To().Hex() {
// 	// 	fmt.Println("err 3 ", tx.To().Hex(), " to ", to, " txData.Hash ", txData.Hash)
// 	// 	return false
// 	// }
// 	receip, _, err := ethclient.GetTransactionReceiptResult(context.Background(), common.HexToHash(txData.Hash))
// 	if err != nil {
// 		return false
// 	}
// 	if receip.Status != 1 {
// 		return false
// 	}
// 	currentHeader, err := ethclient.HeaderByNumber(context.Background(), nil)
// 	if err != nil {
// 		return false
// 	}

// 	confirmed := currentHeader.Number.Uint64() - txData.BlockNumber.Uint64()

// 	// broadcast message
// 	// TODO refactor
// 	txData.Confirmations = int64(confirmed)
// 	if confirmed >= uint64(eth.Confirmations) {
// 		if txtype == "in" {
// 			txData.Status = types.StatusDeposited
// 		} else {
// 			txData.Status = types.StatusWithdrawed
// 		}
// 	}
// 	b, _ := json.Marshal(txData)

// 	if txtype == "in" {
// 		eth.GetBroadcaster().NotifyTransaction(txtype, "ETH", string(b), txData.To)
// 	} else {
// 		eth.GetBroadcaster().NotifyWithdrawTx(txData.ScID, "ETH", string(b))
// 	}

// 	if confirmed >= uint64(eth.Confirmations) {
// 		return true
// 	}
// 	return false
// }

// // IsTransTaskactionExist true if transaction existed
// func (eth *ETH) IsTransTaskactionExist(txData types.Transaction) bool {
// 	return false
// }

// // IsTransTaskactionExist true if transaction existed
// func (eth *ETH) GetBestBlock() *big.Int {
// 	ethclient, err := ethclient.Dial(eth.FullNodeConfig.RPC)
// 	currentHeader, err := ethclient.HeaderByNumber(context.Background(), nil)

// 	if err != nil {
// 		return nil
// 	}

// 	return currentHeader.Number
// }

// func (eth *ETH) ComposeTransactionData(receiver string, value *big.Int) (destination common.Address, data []byte, val *big.Int) {
// 	return common.HexToAddress(receiver),
// 		[]byte{},
// 		value
// }

// // for eth contract only
// func (eth *ETH) GetMultisigTx(scAddress string, txID string, isBurning bool) (*types.ScTxType, error) {
// 	client, err := ethclient.Dial(eth.FullNodeConfig.RPC)
// 	if err != nil {
// 		log.ERROR.Println(err)
// 		return nil, err
// 	}

// 	address := common.HexToAddress(scAddress)
// 	instance, err := Multisig.NewMultisigwallet(address, client)

// 	if err != nil {
// 		log.ERROR.Println(err)
// 		return nil, err
// 	}

// 	txNumber := new(big.Int)
// 	txNumber, ok := txNumber.SetString(txID, 10)

// 	if !ok {
// 		log.ERROR.Println(err)
// 		return nil, errors.New("TxId is malform ")
// 	}

// 	// double check the txvalue - in case the tx creator fake it
// 	tx, err := instance.Transactions(nil, txNumber)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &types.ScTxType{
// 		Destination: tx.Destination,
// 		Value:       tx.Value,
// 		Data:        tx.Data,
// 	}, nil
// }
