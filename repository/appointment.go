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

func NewAppointmentRepo(ctx context.Context, client *mongo.Client) *AppointmentRepo {
	const repoName = "appointment"
	mongoConfig := config.CONFIG.Repositories[repoName]

	collection := client.Database(mongoConfig.DBName).Collection(mongoConfig.CollName)
	return &AppointmentRepo{
		coll: collection,
		ctx:  ctx,
	}
}

type AppointmentRepo struct {
	coll *mongo.Collection
	ctx  context.Context
}

func (r *AppointmentRepo) Create(payload entity.TAppointment) error {
	oid, err := primitive.ObjectIDFromHex(payload.AppointmentID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	insertResult, err := r.coll.InsertOne(r.ctx, bson.M{"_id": oid})
	if err != nil || insertResult.InsertedID == primitive.NilObjectID {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *AppointmentRepo) Read(id string, result *entity.TAppointment) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = r.coll.FindOne(r.ctx, bson.M{"_id": oid}).Decode(result)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *AppointmentRepo) ReadMany(id string, result *[]entity.TAppointment) error {
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

func (r *AppointmentRepo) Update(appointmentId string, payload entity.TUpdateAppointment) error {
	delResult, err := r.coll.UpdateOne(r.ctx, bson.M{"appointment_id": appointmentId}, bson.M{"$set": payload})
	if err != nil || delResult.MatchedCount == 0 {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *AppointmentRepo) Delete(appointId string) error {
	delResult, err := r.coll.DeleteOne(r.ctx, bson.M{"appointment_id": appointId})
	if err != nil || delResult.DeletedCount == 0 {
		err = mongo.ErrNoDocuments
		logrus.Error(err)
		return err
	}

	return nil
}
