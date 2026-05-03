package middleware

import (
	"encoding/json"
	"fmt"
	dto "shared/models"
	"shared/utils"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func PanelFilter(ctx *context.Context) {
	token := ctx.Input.Session("token")
	if token == nil || token == "" {
		ctx.Redirect(302, "/auth/login")
		return
	}
	//fmt.Println(token)
	apiUrl, _ := beego.AppConfig.String("api_url")

	res, _ := utils.CallAPI("POST", apiUrl+"user/ValidateToken", token.(string), nil)
	defer res.Body.Close()
	var apiResp dto.UserResponse
	json.NewDecoder(res.Body).Decode(&apiResp)
	if apiResp.Email == "" {
		ctx.Redirect(302, "/auth/login")
		return
	}
	fmt.Println(apiResp)
	ctx.Input.SetData("user", apiResp)
	ctx.Input.SetData("token", token)
}
