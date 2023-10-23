package utils

import (
	"go-test/pkg/framework"
	"net/http"

	"github.com/kpango/glg"
)

func CheckErrorCtx(ctx *framework.Context, err error) {
	glg.Warnf("%v %v %v", ctx.Request.Method, ctx.Request.URL.Path, err.Error())
	ctx.JSON(http.StatusBadRequest, framework.H{
		"error": err.Error(),
	})
}
