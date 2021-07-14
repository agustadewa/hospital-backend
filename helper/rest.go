package helper

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"gitlab.com/yosiaagustadewa/qsl-util/helper"
)

// Error struct
type Error struct {
	Code    string      `json:"code" bson:"code" binding:"required"`
	Message interface{} `json:"refs,omitempty" bson:"refs,omitempty"`
}

// +++++++++++++ RESPONSE ++++++++++++++++

// Response struct
type Response struct {
	Success bool        `json:"success" bson:"success" binding:"required"`
	Kind    string      `json:"kind,omitempty" bson:"kind,omitempty"`
	Values  interface{} `json:"values,omitempty" bson:"values,omitempty"`
	Error   interface{} `json:"error,omitempty" bson:"error,omitempty"`
}

func BadRequest(c *gin.Context, err error) {
	if err == nil {
		err = errors.New("")
	}

	c.JSON(http.StatusOK, Response{
		Error: Error{
			Code:    "400",
			Message: err.Error(),
		},
	})
}

func Unauthorized(c *gin.Context, err error) {
	if err == nil {
		err = errors.New("")
	}

	c.JSON(http.StatusOK, Response{
		Error: Error{
			Code:    "401",
			Message: err.Error(),
		},
	})
}

func Ok(c *gin.Context, value interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Values:  value,
	})
}

// ++++++++++++ REQUEST VALIDATOR ++++++++++++++

var validate *validator.Validate

// BindValidate function
func BindValidate(o interface{}) error {
	ginValidate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		err := binding.Validator.ValidateStruct(o)
		if err != nil {
			return err
		}
		validate = validator.New()
		err = validate.Struct(o)

		return err
	}

	err := ginValidate.Struct(o)
	if err != nil {
		return err
	}

	validate = validator.New()
	err = validate.Struct(o)

	return err
}

// PairValues function
func PairValues(i, o interface{}) error {
	if i == nil {
		return errors.New("error while pair values, values is nil")
	}
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &o)
	if err != nil {
		return err
	}

	// check type of o
	r := reflect.ValueOf(o)
	if r.Kind() == reflect.Ptr && !r.IsNil() {
		r = r.Elem()
	}
	if r.Kind() != reflect.Struct && r.Kind() != reflect.Interface {

		return nil
	}

	// validate struct :
	err = BindValidate(o)
	if err != nil {

		return err
	}

	return nil
}

// Request struct
type Request struct {
	Kind   string      `json:"kind" bson:"kind" binding:"required"`
	Values interface{} `json:"values" bson:"values" binding:"required"`
}

// ParseRequest function
func parseRequest(c *gin.Context) (*Request, error) {
	var request *Request
	err := c.ShouldBindJSON(&request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func parseKind(c *gin.Context, kind string) (*Request, error) {
	req, err := parseRequest(c)
	if err != nil || strings.TrimSpace(req.Kind) != kind {
		return nil, errors.New("bad request")
	}
	return req, err
}

func ParseKindAndBody(c *gin.Context, kind string, result interface{}) error {
	req, err := parseKind(c, kind)
	if err != nil {
		return err
	}

	err = helper.PairValues(req.Values, result)
	if err != nil {
		return err
	}

	return nil
}

func ParseKindAndBodyArray(c *gin.Context, kind string, result *[]interface{}) error {
	req, err := parseKind(c, kind)
	if err != nil {
		return err
	}

	err = helper.PairValuesArray(req.Values, result)
	if err != nil {
		return err
	}

	return nil
}

// +++++++++ Param +++++++++++++

func HandleParam(c *gin.Context, targetParams ...string) (*Param, error) {
	var param Param
	err := param.handle(c, targetParams...)
	return &param, err

}

type Param struct {
	List map[string]string
}

func (P *Param) handle(c *gin.Context, targetParams ...string) error {
	if targetParams == nil {
		return errors.New("nil param")
	}

	P.List = make(map[string]string, len(targetParams))
	var blackList []string
	for _, tParam := range targetParams {
		val := strings.TrimSpace(c.Param(tParam))
		if val == "" {
			blackList = append(blackList, val)
			continue
		}
		P.List[tParam] = val
	}

	if len(blackList) != 0 {
		return errors.New("param " + strings.Join(targetParams, ",") + " required")
	}

	return nil
}

func (P *Param) Get(param string) string {
	return P.List[param]
}
