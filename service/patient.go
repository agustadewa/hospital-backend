package service

import (
	"context"

	"github.com/agustadewa/hospital-backend/entity"
	"github.com/agustadewa/hospital-backend/repository"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewPatientService(client *mongo.Client) *Patient {
	return &Patient{mongoClient: client}
}

type Patient struct {
	mongoClient *mongo.Client
}

func (s *Patient) GetAccount(ctx context.Context, id string, result *entity.TAccount) error {
	accountSvc := repository.NewAccountRepo(ctx, s.mongoClient)

	if err := accountSvc.Read(id, result); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (s *Patient) CreateAccount(ctx context.Context, id string, payload entity.TCreateAccountReq) error {
	account := entity.TAccount{
		AccountID: id,
		Role:      entity.PATIENT,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Age:       payload.Age,
		Email:     payload.Email,
		Username:  payload.Username,
		Password:  payload.Password,
	}
	if err := repository.NewAccountRepo(ctx, s.mongoClient).Create(account); err != nil {
		logrus.Error("SPatient.CreateAccount.Create.", err)
		return err
	}

	return nil
}
