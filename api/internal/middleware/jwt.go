package middleware

import (
	"shared/utils"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

func JWTMiddleware(ctx *context.Context) {
	authHeader := ctx.Input.Header("Authorization")

	if authHeader == "" {
		ctx.ResponseWriter.WriteHeader(401)
		ctx.WriteString("missing token")
		return
	}

	tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)

	claims, err := utils.ValidateToken(tokenStr)
	if err != nil {
		ctx.ResponseWriter.WriteHeader(401)
		ctx.WriteString("invalid token")
		return
	}

	ctx.Input.SetData("userID", claims.UserID)
}
