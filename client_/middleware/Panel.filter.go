package middleware

import (
	"os"
	dto "shared/models"
	"shared/utils"

	"github.com/beego/beego/v2/server/web/context"
)

func PanelFilter(ctx *context.Context) {
	lang := ctx.Input.Param(":lang")
	if lang == "" {
		lang = os.Getenv("DEFAULTLANG")
	}
	token := ctx.Input.Session("token")
	if token == nil || token == "" {
		ctx.Redirect(302, "/"+lang+"/auth/login")
		return
	}

	tokenn, _ := token.(string)

	apiResp, err := utils.SendRequest[*dto.UserResponse](nil, "user/ValidateToken", "POST", ctx, tokenn)
	if err != nil || apiResp == nil {
		ctx.Redirect(302, "/"+lang+"/auth/login")
		return
	}

	ctx.Input.SetData("user", apiResp)
	ctx.Input.SetData("token", token)
}
