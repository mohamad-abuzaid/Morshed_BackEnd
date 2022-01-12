package helpers

import (
	"log"
	"os"

	"github.com/kataras/iris/v12"
)

const debug = true

func Mdebugf(format string, args ...interface{}) {
	if !debug {
		return
	}

	log.Printf(format, args...)
}

func MwriteInternalServerError(ctx iris.Context) {
	ctx.StopWithJSON(iris.StatusInternalServerError, MnewError(iris.StatusInternalServerError, ctx.Request().Method, ctx.Path(), ""))
}

func MwriteEntityNotFound(ctx iris.Context) {
	ctx.StopWithJSON(iris.StatusNotFound, MnewError(iris.StatusNotFound, ctx.Request().Method, ctx.Path(), "entity does not exist"))
}

func Mgetenv(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	return v
}
