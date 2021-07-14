package repository

import (
	"context"

	"github.com/agustadewa/hospital-backend/config"
	"github.com/agustadewa/hospital-backend/entity"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAccountRepo(ctx context.Context, client *mongo.Client) *AccountRepo {
	const repoName = "account"
	mongoConfig := config.CONFIG.Repositories[repoName]

	collection := client.Database(mongoConfig.DBName).Collection(mongoConfig.CollName)
	return &AccountRepo{
		coll: collection,
		ctx:  ctx,
	}
}

type AccountRepo struct {
	coll *mongo.Collection
	ctx  context.Context
}

func (r *AccountRepo) Create(payload entity.TAccount) error {
	insertResult, err := r.coll.InsertOne(r.ctx, payload)
	if err != nil || insertResult.InsertedID == primitive.NilObjectID {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *AccountRepo) Read(accountId string, result *entity.TAccount) error {
	err := r.coll.FindOne(r.ctx, bson.M{"account_id": accountId}).Decode(result)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *AccountRepo) ReadByEmail(email string, result *entity.TAccount) error {
	err := r.coll.FindOne(r.ctx, bson.M{"email": email}).Decode(result)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *AccountRepo) ReadMany(id string, result *[]entity.TAccount) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Error(err)
		return err
	}

	cursor, err := r.coll.Find(r.ctx, bson.M{"_id": oid})
	if err != nil {
		logrus.Error(err)
		return err
	}

	if err = cursor.All(r.ctx, result); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (r *AccountRepo) CheckIsAdmin(accountId string) (bool, error) {
	count, err := r.coll.CountDocuments(r.ctx, bson.M{
		"account_id": accountId,
		"role":       entity.ADMINISTRATOR,
	})
	if err != nil && err != mongo.ErrNoDocuments {
		logrus.Error(err)
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (r *AccountRepo) CheckAccountByEmail(email string) (bool, error) {
	count, err := r.coll.CountDocuments(r.ctx, bson.M{"email": email})
	if err != nil && err != mongo.ErrNoDocuments {
		logrus.Error(err)
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (r *AccountRepo) CheckAccountByUsername(username string) (bool, error) {
	count, err := r.coll.CountDocuments(r.ctx, bson.M{"username": username})
	if err != nil && err != mongo.ErrNoDocuments {
		logrus.Error(err)
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
