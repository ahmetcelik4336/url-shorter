package middleware

import (
	"os"

	"github.com/beego/beego/v2/server/web/context"
)

func ApiKeyMiddleware(ctx *context.Context) {
	authHeader := ctx.Input.Header("APIKEY")
	apikey := os.Getenv("APIKEY")
	if authHeader == "" || authHeader != apikey {
		ctx.ResponseWriter.WriteHeader(401)
		ctx.WriteString("missing token")
		return
	}
}
