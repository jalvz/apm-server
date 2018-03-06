package client

import (
	"github.com/olivere/elastic"
	"fmt"
	"strings"
	"time"
	"context"
	"net/http"
	"sync"
)

var c *elastic.Client
var once sync.Once

var backPressure sync.WaitGroup

func getClient() *elastic.Client {
	// needs Configs
	once.Do(func() {
		var err error
		c, err = elastic.NewClient(
			elastic.SetURL("http://localhost:9200"),
			elastic.SetSniff(false),
			elastic.SetRetrier(NewRetrier()),
			elastic.SetHealthcheckInterval(time.Second * 2),
			elastic.SetHealthcheckTimeout(time.Second),
			elastic.SetHealthcheckCallback(HealthStatusTracker()),
			)
		// defer c.Stop() needs this somewhere else
		if err != nil {
			panic(err)
		}
	})
	return c
}

func BulkProcessor(processor string, workers int) *elastic.BulkProcessor {
	bulk := elastic.NewBulkProcessorService(getClient()).
		Name(strings.Join([]string{processor, "bulk", "processor"}, "-")).
		//Stats(true).
		//Backoff(elastic.NewConstantBackoff(time.Minute * 1)).
		Workers(workers).
		//BulkActions(1000).BulkSize(1024 * 1024)  // flush every 1000 events or 1 mb accumulated
		FlushInterval(time.Second)

	bulkProcessor, err := bulk.Do(context.Background())
	if err != nil {
		panic(err)
	}
	return bulkProcessor
}


type Retrier struct {
	backoff elastic.Backoff
}

func NewRetrier() *Retrier {
	return &Retrier{
		backoff: elastic.NewExponentialBackoff(5 * time.Second, 5 * time.Minute),
	}
}

func (r *Retrier) Retry(ctx context.Context, retry int, req *http.Request, resp *http.Response, err error) (time.Duration, bool, error) {
	wait, stop := r.backoff.Next(retry)
	// can do eg.: if pressure is insane, return err to abort retrying
	fmt.Println("RETRYING ", retry)
	return wait, stop, nil
}


func SaveToES(data [][]byte, bulk *elastic.BulkProcessor) {
	// needs Configs
	for _, e := range data{
		r := elastic.NewBulkIndexRequest().
			Index("apm-7.0.0-alpha1-2018.03.06").
			Type("doc").
			Doc(e)
		bulk.Add(r)
	}
	backPressure.Wait() // blocks until Done() has been called has many times as Add(1)
	// bulk.Flush() is too much
}


func HealthStatusTracker() func(map[string]bool){
	var m sync.Mutex
	var currentlyHealthy = true

	return func(newState map[string]bool) {
		var healthBalance int
		// a healthy cluster has more than 50% nodes alive
		for _, v := range newState {
			if v {
				healthBalance = healthBalance + 1
			} else {
				healthBalance = healthBalance - 1
			}
		}
		m.Lock()
		if currentlyHealthy && healthBalance < 1 {
			// on transition from healthy to not healthy, tell consumers to not pull the buffer (pause)
			fmt.Println("GOING UNHEALTHY!")
			backPressure.Add(1)
			currentlyHealthy = false
		}
		if !currentlyHealthy && healthBalance > 0 {
			// on transition from not healthy to healthy, tell consumers to resume pulling the buffer
			fmt.Println("GOING HEALTHY!")
			backPressure.Done()
			currentlyHealthy = true
		}
		m.Unlock()
	}

}
