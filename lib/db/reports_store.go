package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/officer47p/addressport/lib/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: Change the coll value
const reportsColl = "address"

type ReportsStore interface {
	Dropper

	GetReportById(context.Context, string) (*types.Report, error)
	GetReportsByAddress(context.Context, string) ([]*types.Report, error)
	GetReports(context.Context) ([]*types.Report, error)
	InsertReport(context.Context, *types.Report) (*types.Report, error)
	DeleteReport(context.Context, string) error
	UpdateReport(context.Context, string, types.UpdateReportParams) error
}

type MongoReportsStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoReportsStore(client *mongo.Client, dbname string) *MongoReportsStore {
	return &MongoReportsStore{
		client: client,
		coll:   client.Database(dbname).Collection(reportsColl),
	}
}

func (s *MongoReportsStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping reports collection")
	return s.coll.Drop(ctx)
}

func (s *MongoReportsStore) GetReports(ctx context.Context) ([]*types.Report, error) {

	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var reports []*types.Report = []*types.Report{}
	if err := cur.All(ctx, &reports); err != nil {
		return nil, err
	}

	return reports, nil
}

func (s *MongoReportsStore) GetReportById(ctx context.Context, id string) (*types.Report, error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var report types.Report
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&report)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (s *MongoReportsStore) GetReportsByAddress(ctx context.Context, address string) ([]*types.Report, error) {

	cur, err := s.coll.Find(ctx, bson.M{"address": address})
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var reports []*types.Report = []*types.Report{}
	if err := cur.All(ctx, &reports); err != nil {
		return nil, err
	}

	return reports, nil
}

func (s *MongoReportsStore) InsertReport(ctx context.Context, report *types.Report) (*types.Report, error) {

	res, err := s.coll.InsertOne(ctx, report)
	if err != nil {
		return nil, err
	}

	report.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return report, nil
}

func (s *MongoReportsStore) DeleteReport(ctx context.Context, id string) error {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no report with the given id was found")
	}

	return nil
}

func (s *MongoReportsStore) UpdateReport(ctx context.Context, id string, params types.UpdateReportParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.D{
		{Key: "$set", Value: params},
	}
	_, err = s.coll.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return err
	}

	return nil
}
