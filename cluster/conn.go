package cluster

import (
	"GoCouchbase/config"
	"GoCouchbase/utils/logutil"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"log"
	"sync"
	"time"
)

var (
	instance *gocb.Cluster
	once     sync.Once
)

func GetCluster() (*gocb.Cluster, error) {
	once.Do(func() {
		var err error
		instance, err = connect()
		if err != nil {
			log.Fatal(err)
		}
	})
	return instance, nil
}

func connect() (*gocb.Cluster, error) {
	logger := logutil.GetLogger()
	cfg := config.GetConfig()

	logger.Println(logutil.Info, false, "Connecting to cluster...")

	cluster, err := gocb.Connect(cfg.CouchbaseConfig.URL, gocb.ClusterOptions{
		Username: cfg.CouchbaseConfig.User,
		Password: cfg.CouchbaseConfig.Password,
	})

	if err != nil {
		logger.Println(logutil.Error, false, "Error connecting to cluster:", err)
		return nil, err
	}

	pingResults, err := cluster.Ping(&gocb.PingOptions{
		ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeQuery},
		Timeout:      2 * time.Second,
	})
	if err != nil {
		logger.Println(logutil.Error, false, "Error pinging the cluster:", err)
		return nil, err
	}

	for service, reports := range pingResults.Services {
		for _, report := range reports {
			if report.State != gocb.PingStateOk {
				logger.Println(logutil.Error, false, "Ping failed for service", service, ":", report.Error)
				return nil, fmt.Errorf("ping failed for service %v: %s", service, report.Error)
			}
		}
	}

	logger.Println(logutil.Info, false, "Connected to cluster.")
	return cluster, nil
}
