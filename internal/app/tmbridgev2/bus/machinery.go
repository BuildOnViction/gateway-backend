package bus

import (
	"context"
	"math/big"

	"cloud.google.com/go/pubsub"
	"github.com/anhntbk08/gateway/internal/common"
	"github.com/anhntbk08/machinery/v1"
	"github.com/anhntbk08/machinery/v1/backends/result"
	"github.com/anhntbk08/machinery/v1/config"
	"github.com/anhntbk08/machinery/v1/tasks"
	"google.golang.org/api/option"
)

type Bus struct {
	server *machinery.Server
}

func NewBus(conf common.JobqueueConfig) (*Bus, error) {
	pubsubClient, err := pubsub.NewClient(
		context.Background(),
		conf.ProjectID,
		option.WithServiceAccountFile(conf.GoogleAuth),
	)
	if err != nil {
		return nil, err
	}

	cnf := &config.Config{
		Broker:        conf.Broker,
		DefaultQueue:  conf.DefaultQueue,
		ResultBackend: conf.ResultBackend,
		GCPPubSub: &config.GCPPubSubConfig{
			Client: pubsubClient,
		},
		MongoDB: &config.MongoDBConfig{
			Database: conf.Mongodb,
		},
	}

	machinery, err := machinery.NewServer(cnf)

	if err != nil {
		return nil, err
	}

	return &Bus{
		machinery,
	}, nil

}

func (bc *Bus) GetBusServer() *machinery.Server {
	return bc.server
}

func (bc *Bus) ProposeMintingTx(tx string, to string, blockNumber *big.Int, txId *big.Int) (*result.AsyncResult, error) {
	verifyingTask := &tasks.Signature{
		Name: "verify_minting_tx",
		Args: []tasks.Arg{
			{
				Type:  "uint64",
				Value: blockNumber.Int64(),
			},
			{
				Type:  "string",
				Value: to,
			},
			{
				Type:  "string",
				Value: tx,
			},
			{
				Type:  "string",
				Value: txId.String(),
			},
		},
		IgnoreWhenTaskNotRegistered: true,
	}

	verifyingTask.RetryCount = 10

	asyncResult, err := bc.server.SendTask(verifyingTask)

	return asyncResult, err
}

func (bc *Bus) ProposeWithdrawingTx(tomoScID, to, ethScID, coin string) (*result.AsyncResult, error) {
	verifyingTask := &tasks.Signature{
		Name: "verify_withdrawing_tx",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: tomoScID,
			},
			{
				Type:  "string",
				Value: to,
			},
			{
				Type:  "string",
				Value: ethScID,
			},
			{
				Type:  "string",
				Value: coin,
			},
		},
		IgnoreWhenTaskNotRegistered: true,
	}

	verifyingTask.RetryCount = 10

	asyncResult, err := bc.server.SendTask(verifyingTask)

	return asyncResult, err
}

func (bc *Bus) NotifyExecutedMintedTx(tx string, to string, scTx string, confirmedTx string) (*result.AsyncResult, error) {
	verifyingTask := &tasks.Signature{
		Name: "new_minted_tx",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: tx,
			},
			{
				Type:  "string",
				Value: to,
			},
			{
				Type:  "string",
				Value: scTx,
			},
			{
				Type:  "string",
				Value: confirmedTx,
			},
		},
		IgnoreWhenTaskNotRegistered: true,
	}

	verifyingTask.RetryCount = 2

	asyncResult, err := bc.server.SendTask(verifyingTask)

	return asyncResult, err
}

func (bc *Bus) NotifyWithdrawTx(scTx string, coin, confirmedTx string) (*result.AsyncResult, error) {
	verifyingTask := &tasks.Signature{
		Name: "notify_withdraw_transaction",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: scTx,
			},
			{
				Type:  "string",
				Value: coin,
			},
			{
				Type:  "string",
				Value: confirmedTx,
			},
		},
		IgnoreWhenTaskNotRegistered: true,
	}

	verifyingTask.RetryCount = 2

	asyncResult, err := bc.server.SendTask(verifyingTask)

	return asyncResult, err
}
