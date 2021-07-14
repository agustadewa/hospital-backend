package middleware

import (
	"errors"

	"github.com/agustadewa/hospital-backend/helper"
	"github.com/agustadewa/hospital-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HeaderVerifier(c *gin.Context) {
	token, err := service.VerifyTokenHeader(c)
	if err != nil {
		logrus.Error("MHeaderVerifier.VerifyTokenHeader.", err)
		helper.Unauthorized(c, err)
		c.Abort()
		return
	}

	accessToken, err := service.DecodeToken(token)
	if err != nil {
		logrus.Error("MHeaderVerifier.DecodeToken.", err)
		helper.Unauthorized(c, err)
		c.Abort()
		return
	}

	if accessToken.Claims.AccountID == "" {
		logrus.Error("MHeaderVerifier.AccountID.IsNil")
		helper.Unauthorized(c, errors.New("invalid token"))
		c.Abort()
		return
	}

	c.Request.Header.Set("Account-Id", accessToken.Claims.AccountID)

	c.Next()
}
