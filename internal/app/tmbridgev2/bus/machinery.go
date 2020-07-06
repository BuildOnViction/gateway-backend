package bus

import (
	"context"

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

func (bc *Bus) NotifySmartContractAddressesChange(projectID string, removedAddresses []string, newAddresses []string) (*result.AsyncResult, error) {
	notifyingTask := &tasks.Signature{
		Name: "notify_project_sc_changed",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: projectID,
			},
			{
				Type:  "[]string",
				Value: removedAddresses,
			},
			{
				Type:  "[]string",
				Value: newAddresses,
			},
		},
		IgnoreWhenTaskNotRegistered: true,
	}

	notifyingTask.RetryCount = 10

	asyncResult, err := bc.server.SendTask(notifyingTask)

	return asyncResult, err
}

func (bc *Bus) CreateSyncingSmartContractTransaction(projectID string, newAddress string) (*result.AsyncResult, error) {
	syncingTask := &tasks.Signature{
		Name: "sync_smartcontract_transaction",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: projectID,
			},
		},
		IgnoreWhenTaskNotRegistered: true,
	}

	syncingTask.RetryCount = 1

	asyncResult, err := bc.server.SendTask(syncingTask)

	return asyncResult, err
}
