package db

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	DBNAME     = "addressport"
	TestDBNAME = "addressport"
	DBURI      = "mongodb://root:root@localhost:27017/"
)

type Neo4jGraphDatabase struct {
	Driver  *neo4j.DriverWithContext
	Context *context.Context
}

func (gdb *Neo4jGraphDatabase) Disconnect() error {
	d := *gdb.Driver
	return d.Close(*gdb.Context)
}

func NewNeo4jGraphDatabase(dbUri string, dbUser string, dbPassword string) (*Neo4jGraphDatabase, error) {
	ctx := context.Background()

	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		return nil, err
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, err
	}
	return &Neo4jGraphDatabase{Driver: &driver, Context: &ctx}, nil
}

type Dropper interface {
	Drop(context.Context) error
}
