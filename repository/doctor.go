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

func NewDoctorRepo(ctx context.Context, client *mongo.Client) *DoctorRepo {
	const repoName = "doctor"
	mongoConfig := config.CONFIG.Repositories[repoName]

	collection := client.Database(mongoConfig.DBName).Collection(mongoConfig.CollName)
	return &DoctorRepo{
		coll: collection,
		ctx:  ctx,
	}
}

type DoctorRepo struct {
	coll *mongo.Collection
	ctx  context.Context
}

func (r *DoctorRepo) Create(payload entity.TDoctor) error {
	insertResult, err := r.coll.InsertOne(r.ctx, payload)
	if err != nil || insertResult.InsertedID == primitive.NilObjectID {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *DoctorRepo) CheckDoctorByName(firstName, lastName string) (bool, error) {
	count, err := r.coll.CountDocuments(r.ctx, bson.M{
		"first_name": firstName,
		"last_name":  lastName,
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

// ------
func (r *DoctorRepo) Read(doctorId string, result *entity.TDoctor) error {
	err := r.coll.FindOne(r.ctx, bson.M{"doctor_id": doctorId}).Decode(result)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *DoctorRepo) ReadMany(id string, result *[]entity.TDoctor) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Error(err)
		return err
	}

	cursor, err := r.coll.Find(r.ctx, bson.M{"doctor_id": oid})
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

func (r *DoctorRepo) ReadManyByName(keyword string, result *[]entity.TDoctor) error {
	reKeywordObj := primitive.Regex{
		Pattern: keyword,
		Options: "i",
	}

	cursor, err := r.coll.Find(r.ctx, bson.M{
		"$or": bson.A{
			bson.M{"first_name": reKeywordObj},
			bson.M{"last_name": reKeywordObj},
		},
	})
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

func (r *DoctorRepo) Delete(doctorId string) error {
	delResult, err := r.coll.DeleteOne(r.ctx, bson.M{"doctor_id": doctorId})
	if err != nil {
		logrus.Error(err)
		return err
	}

	if delResult.DeletedCount == 0 {
		err = mongo.ErrNoDocuments
		logrus.Error(err)
		return err
	}

	return nil
}
