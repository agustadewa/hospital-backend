package helper

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func NewMongoConnection(ctx context.Context, uri string) *mongo.Client {
	Client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	err = Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return Client
}

//  Passing session inside transactionCallback by value will referring to another session address.
//  Session must be passed by reference instead of passed by value
func UseMongoTransaction(ctx context.Context, mongoClient *mongo.Client, transactionCallback func(sc *mongo.SessionContext) error, tryCount int) error {

	// Starting session
	txSession, err := mongoClient.StartSession()
	if err != nil {
		return err
	}

	// Doing transaction with txSession
	err = mongo.WithSession(ctx, txSession, func(sc mongo.SessionContext) error {

		// Setup transaction options
		mc := 5 * time.Second
		wc := writeconcern.New(writeconcern.WMajority())
		rc := readconcern.Snapshot()
		opt := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc).SetMaxCommitTime(&mc)

		// Starting transaction
		if err = txSession.StartTransaction(opt); err != nil {
			return err
		}

		// Passing internal session context reference as "sc" to transactionCallback
		err = transactionCallback(&sc)
		if err != nil {

			// Retry aborting transaction if there is a network problem
			if errIterate := tryIterate(txSession.AbortTransaction, sc, tryCount); errIterate != nil {
				return errIterate
			}
			return err
		}

		// Retry committing transaction if there is a network problem
		if errIterate := tryIterate(txSession.CommitTransaction, sc, tryCount); errIterate != nil {
			return errIterate
		}

		return nil
	})
	if err != nil {
		// End session if there is an error
		txSession.EndSession(ctx)
		return err
	}

	// End session after committing or aborting transaction
	txSession.EndSession(ctx)

	return nil
}

func tryIterate(callback func(ctx context.Context) error, sc context.Context, tryCount int) error {

	// tryCount value should be higher than 0
	if tryCount < 1 {
		tryCount = 1
	}

	// selecting between these cases:
	// 1. callback not returning an error, then break the loop and return nil
	// 2. callback returning an error as many as tryCount value, then return that error
	for tryCountLeft := tryCount + 1; tryCountLeft != 0; tryCountLeft-- {
		if err := callback(sc); err != nil {

			// close the loop before reaching 0, then return that error
			if tryCountLeft == 1 {
				return err
			}

			// continue the loop if still returning an error
			continue
		} else {

			// break the loop if success
			break
		}
	}

	return nil
}
