package db

import (
	"context"
	"fmt"

	"github.com/officer47p/addressport/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const addressColl = "address"

type AddressStore interface {
	Dropper

	GetAddressById(context.Context, string) (*types.Address, error)
	GetAddressByAddress(context.Context, string) ([]*types.Address, error)
	GetAddresses(context.Context) ([]*types.Address, error)
	InsertAddress(context.Context, *types.Address) (*types.Address, error)
	// DeleteAddress(context.Context, string) error
	// UpdateAddress(ctx context.Context, filter bson.M, params types.UpdateAddressParams) error
}

type MongoAddressStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoAddressStore(client *mongo.Client, dbname string) *MongoAddressStore {
	return &MongoAddressStore{
		client: client,
		coll:   client.Database(dbname).Collection(addressColl),
	}
}

func (s *MongoAddressStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping user collection")
	return s.coll.Drop(ctx)
}

func (s *MongoAddressStore) GetAddresses(ctx context.Context) ([]*types.Address, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var addresses []*types.Address = []*types.Address{}
	if err := cur.All(ctx, &addresses); err != nil {
		return nil, err
	}
	return addresses, nil
}

func (s *MongoAddressStore) GetAddressById(ctx context.Context, id string) (*types.Address, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var address types.Address
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&address)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (s *MongoAddressStore) GetAddressByAddress(ctx context.Context, address string) ([]*types.Address, error) {
	cur, err := s.coll.Find(ctx, bson.M{"address": address})
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var addresses []*types.Address = []*types.Address{}
	if err := cur.All(ctx, &addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (s *MongoAddressStore) InsertAddress(ctx context.Context, address *types.Address) (*types.Address, error) {
	// var insertedUser types.Address
	res, err := s.coll.InsertOne(ctx, address)
	if err != nil {
		return nil, err
	}

	address.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return address, nil
}

// func (s *MongoAddressStore) DeleteAddress(ctx context.Context, id string) error {
// 	oid, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": oid})
// 	if err != nil {
// 		return err
// 	}

// 	if res.DeletedCount == 0 {
// 		return errors.New("no address with the given id was found")
// 	}

// 	return nil
// }

// func (s *MongoAddressStore) UpdateAddress(ctx context.Context, filter bson.M, params types.UpdateAddressParams) error {
// 	// to prevent changing a field we can do this:
// 	// if values["email"] != nil {
// 	// 	return errors.New("can't change email")
// 	// }
// 	update := bson.D{
// 		{Key: "$set", Value: params},
// 	}

// 	_, err := s.coll.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return err
// 	}
// 	return nil

// }
