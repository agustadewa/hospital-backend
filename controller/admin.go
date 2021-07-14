package controller

import (
	"errors"

	"github.com/agustadewa/hospital-backend/entity"
	"github.com/agustadewa/hospital-backend/helper"
	"github.com/agustadewa/hospital-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAdminController(client *mongo.Client) *adminController {
	return &adminController{
		MongoClient:    client,
		AdminService:   service.NewAdminService(client),
		PatientService: service.NewPatientService(client),
	}
}

type adminController struct {
	MongoClient    *mongo.Client
	AdminService   *service.Admin
	PatientService *service.Patient
}

func (ctrl *adminController) GetAccount(c *gin.Context) {
	ctx := c.Request.Context()

	paramObj, err := helper.HandleParam(c, "account_id")
	if err != nil {
		logrus.Error("CAdmin.GetAccount.", err)
		helper.BadRequest(c, err)
		return
	}

	accountId := paramObj.Get("account_id")

	selfAccountID := c.GetHeader("Account-Id")
	isAdmin, err := ctrl.AdminService.IsAdmin(ctx, selfAccountID)
	if err != nil || !isAdmin {
		logrus.Error(err)
		helper.Unauthorized(c, errors.New("not admin"))
		return
	}

	var result entity.TAccount
	if err := ctrl.AdminService.GetAccount(ctx, accountId, &result); err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.Error("CGetAccount.GetAccount.NoDocuments.", err)
			helper.BadRequest(c, errors.New("not found"))
			return
		} else {
			logrus.Error(err)
			helper.BadRequest(c, err)
			return
		}
	}

	helper.Ok(c, result)
}

func (ctrl *adminController) CreateAccount(c *gin.Context) {
	ctx := c.Request.Context()

	var reqBody entity.TCreateAccountReq
	if err := helper.ParseKindAndBody(c, "account#create", &reqBody); err != nil {
		logrus.Error(err)
		helper.BadRequest(c, err)
		return
	}

	selfAccountID := c.GetHeader("Account-Id")
	isAdmin, err := ctrl.AdminService.IsAdmin(ctx, selfAccountID)
	if err != nil || !isAdmin {
		logrus.Error(err)
		helper.Unauthorized(c, errors.New("not admin"))
		return
	}

	isExists, err := ctrl.AdminService.CheckAccountByEmailAndUsername(ctx, reqBody.Email, reqBody.Username)
	if err != nil {
		logrus.Error("CCreateAccount.CheckAccountByEmailAndUsername", err)
		helper.BadRequest(c, err)
		return
	}

	if isExists {
		logrus.Info("CCreateAccount.CheckAccountByEmailAndUsername.Exists")
		helper.BadRequest(c, errors.New("already exists"))
		return
	}

	newOID := primitive.NewObjectID().Hex()
	if err := ctrl.AdminService.CreateAccount(ctx, newOID, reqBody); err != nil {
		logrus.Info("CCreateAccount.CreateAccount.")
		helper.BadRequest(c, errors.New("can't create account"))
		return
	}

	helper.Ok(c, entity.TCreateAccountRes{AccountID: newOID})
}

// ++++++++++++++++++ DOCTOR +++++++++++++++++++

func (ctrl *adminController) GetDoctor(c *gin.Context) {
	ctx := c.Request.Context()

	paramObj, err := helper.HandleParam(c, "doctor_id")
	if err != nil {
		logrus.Error("CAdmin.GetDoctor.", err)
		helper.BadRequest(c, err)
		return
	}

	doctorId := paramObj.Get("doctor_id")

	selfAccountID := c.GetHeader("Account-Id")
	isAdmin, err := ctrl.AdminService.IsAdmin(ctx, selfAccountID)
	if err != nil || !isAdmin {
		logrus.Error(err)
		helper.Unauthorized(c, errors.New("not admin"))
		return
	}

	var result entity.TDoctor
	if err := ctrl.AdminService.GetDoctor(ctx, doctorId, &result); err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.Error("CGetDoctor.GetDoctor.NoDocuments.", err)
			helper.BadRequest(c, errors.New("not found"))
			return
		} else {
			logrus.Error(err)
			helper.BadRequest(c, err)
			return
		}
	}

	helper.Ok(c, result)
}

func (ctrl *adminController) GetManyDoctor(c *gin.Context) {
	ctx := c.Request.Context()

	var reqBody entity.TGetManyDoctorReq
	if err := helper.ParseKindAndBody(c, "doctor#getmany", &reqBody); err != nil {
		logrus.Error(err)
		helper.BadRequest(c, err)
		return
	}

	var result []entity.TDoctor
	if err := ctrl.AdminService.GetManyDoctorByName(ctx, reqBody.Keyword, &result); err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.Error("CGetDoctor.GetDoctor.NoDocuments.", err)
			helper.BadRequest(c, errors.New("not found"))
			return
		} else {
			logrus.Error(err)
			helper.BadRequest(c, err)
			return
		}
	}

	helper.Ok(c, result)
}

func (ctrl *adminController) CreateDoctor(c *gin.Context) {
	ctx := c.Request.Context()

	var reqBody entity.TCreateDoctorReq
	if err := helper.ParseKindAndBody(c, "doctor#create", &reqBody); err != nil {
		logrus.Error(err)
		helper.BadRequest(c, err)
		return
	}

	selfAccountID := c.GetHeader("Account-Id")
	isAdmin, err := ctrl.AdminService.IsAdmin(ctx, selfAccountID)
	if err != nil || !isAdmin {
		logrus.Error(err)
		helper.Unauthorized(c, errors.New("not admin"))
		return
	}

	isExists, err := ctrl.AdminService.CheckDoctorByName(ctx, reqBody.FirstName, reqBody.LastName)
	if err != nil {
		logrus.Error("CCreateDoctor.CheckDoctorName", err)
		helper.BadRequest(c, err)
		return
	}

	if isExists {
		logrus.Info("CCreateDoctor.CheckDoctorName.Exists")
		helper.BadRequest(c, errors.New("already exists"))
		return
	}

	newOID := primitive.NewObjectID().Hex()
	if err := ctrl.AdminService.CreateDoctor(ctx, newOID, reqBody); err != nil {
		logrus.Info("CCreateDoctor.CreateDoctor.")
		helper.BadRequest(c, errors.New("can't create doctor"))
		return
	}

	helper.Ok(c, entity.TCreateDoctorRes{DoctorID: newOID})
}

func (ctrl *adminController) DeleteDoctor(c *gin.Context) {
	ctx := c.Request.Context()

	paramObj, err := helper.HandleParam(c, "doctor_id")
	if err != nil {
		logrus.Error("CAdmin.DeleteDoctor.", err)
		helper.BadRequest(c, err)
		return
	}

	doctorId := paramObj.Get("doctor_id")

	selfAccountID := c.GetHeader("Account-Id")
	isAdmin, err := ctrl.AdminService.IsAdmin(ctx, selfAccountID)
	if err != nil || !isAdmin {
		logrus.Error(err)
		helper.Unauthorized(c, errors.New("not admin"))
		return
	}

	if err := ctrl.AdminService.DeleteDoctor(ctx, doctorId); err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.Error("CDeleteDoctor.DeleteDoctor.NotFound.", err)
			helper.BadRequest(c, errors.New("not found"))
			return

		} else {
			logrus.Error(err)
			helper.BadRequest(c, err)
			return
		}
	}

	helper.Ok(c, nil)
}
