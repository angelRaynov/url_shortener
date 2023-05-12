package database

import (
	"github.com/gocql/gocql"
	"log"
	"time"
	"url_shortener/config"
)

type DB struct {
	Session *gocql.Session
	cfg     *config.Application
}

func InitDB(cfg *config.Application) *DB {
	// create a new cluster configuration object
	cluster := gocql.NewCluster(cfg.DBHost)

	// set the keyspace to use
	cluster.Keyspace = cfg.DBKeyspace

	// set the consistency level to use for queries
	cluster.Consistency = gocql.Quorum

	// set the retry policy to use for failed queries
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
		NumRetries: 3,
		Min:        1 * time.Second,
		Max:        10 * time.Second,
	}

	// create a new session using the cluster configuration
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("create cassandra session:%v", err)
	}

	return &DB{
		Session: session,
	}
}
