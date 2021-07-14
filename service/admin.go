package service

import (
	"context"

	"github.com/agustadewa/hospital-backend/entity"
	"github.com/agustadewa/hospital-backend/repository"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAdminService(client *mongo.Client) *Admin {
	return &Admin{mongoClient: client}
}

type Admin struct {
	mongoClient *mongo.Client
}

func (s *Admin) IsAdmin(ctx context.Context, accountId string) (bool, error) {
	return repository.NewAccountRepo(ctx, s.mongoClient).CheckIsAdmin(accountId)
}

func (s *Admin) GetAccount(ctx context.Context, accountId string, result *entity.TAccount) error {
	if err := repository.NewAccountRepo(ctx, s.mongoClient).Read(accountId, result); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (s *Admin) CheckAccountByEmailAndUsername(ctx context.Context, email, username string) (bool, error) {
	accountRepo := repository.NewAccountRepo(ctx, s.mongoClient)
	emailIsExists, err := accountRepo.CheckAccountByEmail(email)
	if err != nil {
		return false, err
	}
	usernameIsExists, err := accountRepo.CheckAccountByUsername(username)
	if err != nil {
		return false, err
	}

	if emailIsExists || usernameIsExists {
		return true, nil
	}

	return false, nil
}

func (s *Admin) CreateAccount(ctx context.Context, id string, payload entity.TCreateAccountReq) error {
	account := entity.TAccount{
		AccountID: id,
		Role:      entity.ADMINISTRATOR,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Age:       payload.Age,
		Email:     payload.Email,
		Username:  payload.Username,
		Password:  payload.Password,
	}
	if err := repository.NewAccountRepo(ctx, s.mongoClient).Create(account); err != nil {
		logrus.Error("SAdmin.CreateAccount.Create.", err)
		return err
	}

	return nil
}

// +++++++++++++++ DOCTOR ++++++++++++++++

func (s *Admin) CheckDoctorByName(ctx context.Context, firstName, lastName string) (bool, error) {
	doctorRepo := repository.NewDoctorRepo(ctx, s.mongoClient)
	nameIsExists, err := doctorRepo.CheckDoctorByName(firstName, lastName)
	if err != nil {
		return false, err
	}

	if nameIsExists {
		return true, nil
	}

	return false, nil
}

func (s *Admin) GetDoctor(ctx context.Context, doctorId string, result *entity.TDoctor) error {
	if err := repository.NewDoctorRepo(ctx, s.mongoClient).Read(doctorId, result); err != nil {
		logrus.Error("SAdmin.GetDoctor.Read.", err)
		return err
	}

	return nil
}

func (s *Admin) GetManyDoctorByName(ctx context.Context, keyword string, result *[]entity.TDoctor) error {
	if err := repository.NewDoctorRepo(ctx, s.mongoClient).ReadManyByName(keyword, result); err != nil {
		logrus.Error("SAdmin.GetManyDoctor.ReadManyByName.", err)
		return err
	}

	return nil
}

func (s *Admin) CreateDoctor(ctx context.Context, doctorId string, payload entity.TCreateDoctorReq) error {
	doctor := entity.TDoctor{
		DoctorID:  doctorId,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	}
	if err := repository.NewDoctorRepo(ctx, s.mongoClient).Create(doctor); err != nil {
		logrus.Error("SAdmin.CreateDoctor.Create.", err)
		return err
	}

	return nil
}

func (s *Admin) DeleteDoctor(ctx context.Context, doctorId string) error {
	if err := repository.NewDoctorRepo(ctx, s.mongoClient).Delete(doctorId); err != nil {
		logrus.Error("SAdmin.CreateDoctor.Create.", err)
		return err
	}

	return nil
}
