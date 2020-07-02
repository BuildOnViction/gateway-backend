package job

import (
	"errors"
	"sync"

	. "github.com/anhntbk08/gateway/internal/app/tmbridgev2/bus"
	"github.com/anhntbk08/machinery/v1/log"
)

type JobServer struct {
	bus   *Bus
	mutex sync.Mutex
}

func StartServer(bus *Bus) error {
	jobServer := &JobServer{
		bus: bus,
	}
	busServer := jobServer.bus.GetBusServer()

	tasks := map[string]interface{}{
		"new_address": jobServer.ReceiveAddress,
	}
	err := busServer.RegisterTasks(tasks)

	go func() {
		jobWorker := busServer.NewWorker("job_worker", 100)
		err = jobWorker.Launch()
		if err != nil {
			log.ERROR.Println("Can't launch job worker: ", err)
		}
	}()
	// go WatchTx()

	return err
}

// Any way to validate this address
func (js *JobServer) ReceiveAddress(coin string, tomoAddress string, address string, index uint64) error {
	return errors.New("Not implemented yet")
}
