package db

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

var Cluster *gocb.Cluster

func InitCouchbase() error {
	var err error
	Cluster, err = gocb.Connect("couchbase://localhost", gocb.ClusterOptions{
		Username: "Administrator",
		Password: "Aditi17",
	})
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}

	err = Cluster.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		return fmt.Errorf("cluster not ready: %w", err)
	}

	log.Println("✅ Connected to Couchbase cluster")
	return nil
}
