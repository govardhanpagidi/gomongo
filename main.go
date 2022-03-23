package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var err error
var client *mongo.Client
var collection *mongo.Collection
var ctx = context.Background()

type Persister struct {
	client       *mongo.Client
	db           *mongo.Database
	traceEnabled bool
	txSession    mongo.Session
}

func main() {
	FirstTx()
}

func FirstTx() {

	var id = primitive.NewObjectID()
	var doc = bson.M{"_id": id, "hometown": "Atlanta", "year": int32(2003)}
	var result *mongo.UpdateResult
	var update = bson.D{{Key: "$set", Value: bson.D{{Key: "year", Value: int32(2004)}}}}
	if client, err = getMongoClient(); err != nil {
		log.Fatal(err)
	}
	//defer client.Disconnect(ctx)
	mdb := client.Database("hydra")
	p := Persister{db: mdb, client: client}
	coll := p.db.Collection("test")
	//Inserting a doc in test collection
	if _, err = coll.InsertOne(ctx, doc); err != nil {
		log.Fatal(err)
	}

	log.Println("Doc inserted")

	p.BeginTX(ctx)

	//execute statements
	f := func(sc mongo.SessionContext) error {

		res, err := coll.DeleteOne(sc, bson.M{"_id": id})
		if err != nil {
			return err
		}
		if res.DeletedCount != 1 {
			log.Println("delete failed, expected 1 but got", res.DeletedCount)
			return err
		}

		if result, err = coll.UpdateOne(sc, bson.M{"year": 2003}, update); err != nil {

			log.Println(err)
			return err
		}

		return nil
	}

	//Do tx
	err = p.DoTransaction(ctx, f)

	if err != nil {
		p.Rollback(ctx)
	}
	p.Commit(ctx)

	var v bson.M
	if err = coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&v); err != nil {
		log.Fatal(err)
	}
	if v["year"] != int32(2003) {

		log.Fatal("expected 2003 but got", v["year"])
	}

}

func (p *Persister) BeginTX(ctx context.Context) (context.Context, error) {
	defer p.client.Disconnect(ctx)
	//Start session
	sesson, err := p.client.StartSession()
	if err != nil {
		return ctx, err
	}
	//Start Transaction
	if err = sesson.StartTransaction(); err != nil {
		return ctx, err
	}
	p.txSession = sesson
	return ctx, nil
}

func (p *Persister) Commit(ctx context.Context) error {
	p.txSession.CommitTransaction(ctx)
	p.txSession.EndSession(ctx)
	return nil
}

func (p *Persister) Rollback(ctx context.Context) error {
	p.txSession.AbortTransaction(ctx)
	p.txSession.EndSession(ctx)
	return nil
}

func (p *Persister) DoTransaction(ctx context.Context, f func(sc mongo.SessionContext) error) (err error) {

	// DO transaction
	if err = mongo.WithSession(ctx, p.txSession, f); err != nil {
		p.Rollback(ctx)
		return err
	}
	return nil
}

func TestTransactionCommit() {

	var id = primitive.NewObjectID()
	var doc = bson.M{"_id": id, "hometown": "Atlanta", "year": int32(2003)}
	var result *mongo.UpdateResult
	var session mongo.Session
	var update = bson.D{{Key: "$set", Value: bson.D{{Key: "year", Value: int32(2004)}}}}
	if client, err = getMongoClient(); err != nil {
		log.Fatal(err)
	}
	//defer client.Disconnect(ctx)
	collection = client.Database("hydra").Collection("test")
	if _, err = collection.InsertOne(ctx, doc); err != nil {
		log.Fatal(err)
	}

	log.Println("Doc inserted")

	if session, err = client.StartSession(); err != nil {
		log.Fatal(err)
	}
	if err = session.StartTransaction(); err != nil {
		log.Fatal(err)
	}

	//execute statements
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		res, _ := collection.DeleteOne(ctx, bson.M{"_id": id})
		if res.DeletedCount != 1 {
			log.Println("delete failed, expected 1 but got", res.DeletedCount)
		}

		if result, err = collection.UpdateOne(sc, bson.M{"year": 2003}, update); err != nil {
			log.Println(err)
		}
		if result.MatchedCount != 1 || result.ModifiedCount != 1 || err != nil {
			log.Println("replace failed, expected 1 but got", result.MatchedCount)
			if session.AbortTransaction(sc); err != nil {
				log.Fatal(err)
			}
			log.Println("aborting transaction")
		} else {
			if err = session.CommitTransaction(sc); err != nil {
				log.Fatal(err)
			}
			log.Println("Committed transaction")
		}

		return nil
	}); err != nil {

		log.Fatal(err)
	}
	session.EndSession(ctx)

	var v bson.M
	if err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&v); err != nil {
		log.Fatal(err)
	}
	if v["year"] != int32(2002) {

		log.Fatal("expected 2000 but got", v["year"])
	}

}

// Helpers

func getMongoClient() (*mongo.Client, error) {
	uri := "mongodb+srv://reddy:reddy@cluster0.ih8do.mongodb.net/test?authSource=admin&replicaSet=atlas-dm1d4p-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true"
	if os.Getenv("DATABASE_URL") != "" {
		uri = os.Getenv("DATABASE_URL")
	}
	return getMongoClientByURI(uri)
}

func getMongoClientByURI(uri string) (*mongo.Client, error) {
	var err error
	var client *mongo.Client
	opts := options.Client()
	opts.ApplyURI(uri)
	opts.SetMaxPoolSize(5)
	if client, err = mongo.Connect(context.Background(), opts); err != nil {
		return client, err
	}
	client.Ping(context.Background(), nil)
	return client, err
}
