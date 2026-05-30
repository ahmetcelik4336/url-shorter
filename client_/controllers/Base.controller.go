package controllers

import (
	"log"
	"os"
	"shared/helpers"
	dto "shared/models"
	"shared/utils"
	"strconv"
	"unicode/utf8"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

func init() {
	// app.conf'daki kısa kodlarla (tr, en) eşleşecek şekilde kaydediyoruz
	if err := i18n.SetMessage("tr", "conf/locale_tr-TR.ini"); err != nil {
		log.Println("TR dil dosyası yükleme hatası:", err)
	}
	if err := i18n.SetMessage("en", "conf/locale_en-US.ini"); err != nil {
		log.Println("EN dil dosyası yükleme hatası:", err)
	}

	beego.AddFuncMap("i18n", i18n.Tr)
}

type BaseController struct {
	beego.Controller
	i18n.Locale
}

func (c *BaseController) Prepare() {
	lang := c.Ctx.Input.GetData("lang")

	if lang == nil || lang == "" {
		lang = c.Ctx.Input.Param(":lang")
	}
	defaultLang := os.Getenv("DEFAULTLANG")
	if lang == "" || lang == nil {
		lang = defaultLang
	}

	langStr, ok := lang.(string)
	if !ok || langStr == "" {
		langStr = "en"
	}

	languagesMap, _ := beego.AppConfig.GetSection("languages")
	c.Data["supported_langs"] = languagesMap

	if languagesMap[langStr] == "" && utf8.RuneCountInString(langStr) == 2 && langStr != "qr" {
		c.Redirect(helpers.Baseurl(defaultLang, ""), 302)
	}

	c.Lang = langStr
	c.Data["currentLang"] = c.Tr("lang_name")

	c.Data["currentPath"] = c.Ctx.Request.URL.Path
	c.Data["lang"] = c.Lang
	c.Data["slug"] = c.Lang
	c.Data["IsShowFooter"] = "show"
	c.Data["IsShowHeader"] = "show"

	generalSetting, _ := utils.SendRequest[dto.GeneralSettings](nil, "setting/get/general", "GET", c.Ctx, "")
	c.Data["conf"] = generalSetting
	c.Data["userid"] = 0
	c.Data["headerActive"] = ""
	// Oturum / Token Doğrulama
	token := c.GetSession("token")
	if token != nil {
		if tokenn, ok := token.(string); ok && tokenn != "" {
			user, err := utils.SendRequest[*dto.UserResponse](nil, "user/ValidateToken", "POST", c.Ctx, tokenn)
			if err == nil && user != nil && strconv.Itoa(user.Id) != "" {
				c.Data["userid"] = user.Id
				c.Data["user"] = user
			}
		}
	}
}
