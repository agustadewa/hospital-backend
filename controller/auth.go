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

func NewAuthController(mongoClient *mongo.Client) *auth {
	return &auth{
		MongoClient:    mongoClient,
		AdminService:   service.NewAdminService(mongoClient),
		PatientService: service.NewPatientService(mongoClient),
		AuthService:    service.NewAuthService(mongoClient),
	}
}

type auth struct {
	MongoClient    *mongo.Client
	AdminService   *service.Admin
	PatientService *service.Patient
	AuthService    *service.Auth
}

func (ctrl *auth) CheckAuthentication(c *gin.Context) {
	ctx := c.Request.Context()

	accessToken, err := service.VerifyTokenHeader(c)
	if err != nil {
		logrus.Error("CCheckAuthentication.VerifyTokenHeader.", err)
		helper.Unauthorized(c, err)
		return
	}

	token, err := service.DecodeToken(accessToken)
	if err != nil {
		logrus.Error("CCheckAuthentication.DecodeToken.", err)
		helper.Unauthorized(c, err)
		return
	}

	if token.Claims.AccountID == "" {
		logrus.Error("CCheckAuthentication.AccountIDIsMissing.")
		helper.Unauthorized(c, errors.New("bad token"))
		return
	}

	if err = ctrl.AdminService.GetAccount(ctx, token.Claims.AccountID, &entity.TAccount{}); err != nil {
		logrus.Error("CCheckAuthentication.GetAccount.", err)
		helper.Unauthorized(c, err)
		return
	}

	helper.Ok(c, nil)
}

func (ctrl *auth) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var reqBody entity.LoginReq
	if err := helper.ParseKindAndBody(c, "auth#login", &reqBody); err != nil {
		logrus.Error("CLogin.ParseKindAndBody.", err)
		helper.BadRequest(c, err)
		return
	}

	var account entity.TAccount
	isMatch, err := ctrl.AuthService.MatchAndGetAccount(ctx, reqBody.Email, reqBody.Password, &account)
	if err != nil {
		logrus.Error("CLogin.MatchAndGetAccount.", err)
		helper.BadRequest(c, err)
		return
	}

	if !isMatch {
		logrus.Info("CLogin.INVALID")
		helper.BadRequest(c, errors.New("invalid"))
		return
	}

	accessToken, err := ctrl.AuthService.CreateAccessToken(account.AccountID, 3600)
	if err != nil {
		logrus.Error("CLogin.CreateAccessToken.", err)
		helper.BadRequest(c, err)
		return
	}

	helper.Ok(c, entity.LoginRes{
		Authorization: accessToken,
	})
}

func (ctrl *auth) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var reqBody entity.TCreateAccountReq
	if err := helper.ParseKindAndBody(c, "auth#register", &reqBody); err != nil {
		logrus.Error("CLogin.ParseKindAndBody.", err)
		helper.BadRequest(c, err)
		return
	}

	// Check existing account
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

	// Create patient account
	newOID := primitive.NewObjectID().Hex()
	if err := ctrl.PatientService.CreateAccount(ctx, newOID, reqBody); err != nil {
		logrus.Info("CCreateAccount.CreateAccount.")
		helper.BadRequest(c, errors.New("can't create account"))
		return
	}

	// Create Token
	accessToken, err := ctrl.AuthService.CreateAccessToken(newOID, 3600)
	if err != nil {
		logrus.Error("CAuth.Register.CreateAccessToken.", err)
		helper.BadRequest(c, err)
		return
	}

	helper.Ok(c, entity.LoginRes{
		Authorization: accessToken,
	})
}
