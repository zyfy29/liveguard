package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	CustomCodeOk          = 0
	CustomCodeClientError = 1
	CustomCodeServerError = 2
)

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`  // error message, empty for success
	Data any    `json:"data"` // could be marshaled to json
}

func ResponseOk(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Resp{
		Code: CustomCodeOk,
		Msg:  "",
		Data: data,
	})
}

func ResponseClientError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, Resp{
		Code: CustomCodeClientError,
		Msg:  err.Error(),
		Data: nil,
	})
}

func ResponseServerError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, Resp{
		Code: CustomCodeServerError,
		Msg:  err.Error(),
		Data: nil,
	})
}
