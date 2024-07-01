package db

import "context"

const (
	DBNAME     = "addressport"
	TestDBNAME = "addressport"
	DBURI      = "mongodb://root:root@localhost:27017/"
)

type Dropper interface {
	Drop(context.Context) error
}
