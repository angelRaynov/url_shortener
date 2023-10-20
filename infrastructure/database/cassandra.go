package database

import (
	"github.com/gocql/gocql"
	"url_shortener/infrastructure/config"
)

type DB struct {
	Session *gocql.Session
	cfg     *config.Application
}

//func InitDB(cfg *config.Application, logger *zap.SugaredLogger) *DB {
//	// create a new cluster configuration object
//	cluster := gocql.NewCluster(cfg.CassHost)
//
//	// set the keyspace to use
//	cluster.Keyspace = cfg.CassKeyspace
//
//	// set the consistency level to use for queries
//	cluster.Consistency = gocql.Quorum
//
//	// set the retry policy to use for failed queries
//	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
//		NumRetries: 3,
//		Min:        1 * time.Second,
//		Max:        10 * time.Second,
//	}
//
//	// create a new session using the cluster configuration
//	session, err := cluster.CreateSession()
//	if err != nil {
//		logger.Fatalf("creating db session:%v", err)
//	}
//
//	logger.Debug("db started")
//
//	return &DB{
//		Session: session,
//	}
//}
