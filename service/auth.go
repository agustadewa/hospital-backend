package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/agustadewa/hospital-backend/config"
	"github.com/agustadewa/hospital-backend/entity"
	"github.com/agustadewa/hospital-backend/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

import (
	"github.com/sirupsen/logrus"
)

func NewAuthService(mongoClient *mongo.Client) *Auth {
	return &Auth{mongoClient: mongoClient}
}

type Auth struct {
	mongoClient *mongo.Client
}

// CreateAccessToken
func (s *Auth) CreateAccessToken(AccountID string, ExpiredAt time.Duration) (string, error) {

	expiredAt := time.Now().Add((time.Second) * ExpiredAt).Unix()

	claims := jwt.MapClaims{
		"exp": expiredAt,
		"aut": true,
		"aid": AccountID,
		"iat": time.Now().Unix(),
	}

	to := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := to.SignedString([]byte(config.JWTSecretKey))
	if err != nil {
		logrus.Error(err.Error())
		return accessToken, err
	}

	return accessToken, nil
}

// MatchAndGetAccount
func (s *Auth) MatchAndGetAccount(ctx context.Context, email, password string, result *entity.TAccount) (bool, error) {
	var account entity.TAccount
	if err := repository.NewAccountRepo(ctx, s.mongoClient).ReadByEmail(email, &account); err != nil {
		if err != mongo.ErrNoDocuments {
			return false, err
		}
		return false, nil
	}

	// decode password
	// ---------------
	if account.Password != password {
		return false, nil
	}

	*result = account

	return true, nil
}

// --------------------------------

func VerifyTokenHeader(c *gin.Context) (*jwt.Token, error) {
	tokenHeader := c.GetHeader("Authorization")
	if len(tokenHeader) == 0 {
		logrus.Error("MAuth.VerifyTokenHeader.GetHeader.BadToken")
		return nil, errors.New("bad auth header")
	}

	accessToken := strings.SplitAfter(tokenHeader, "Bearer")[1]

	token, err := jwt.Parse(strings.Trim(accessToken, " "), func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecretKey), nil
	})

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if err = token.Claims.Valid(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return token, nil
}

func DecodeToken(accessToken *jwt.Token) (entity.AccessToken, error) {
	var token entity.AccessToken
	stringify, err := json.Marshal(&accessToken)
	if err != nil {
		logrus.Error("MDecodeToken.Marshal.", err)
		return entity.AccessToken{}, err
	}

	if err = json.Unmarshal(stringify, &token); err != nil {
		logrus.Error("MDecodeToken.Unmarshal.", err)
		return entity.AccessToken{}, err
	}

	return token, nil
}
