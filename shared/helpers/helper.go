package helpers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

func Init() {
	beego.AddFuncMap("baseurl", func(lang interface{}, path string) string {
		l, ok := lang.(string)
		if !ok || l == "" {
			val, err := beego.AppConfig.String("defaultlang")
			if err != nil || val == "" {
				l = "en"
			} else {
				l = val
			}
		}
		baseurl, _ := beego.AppConfig.String("baseurl")

		return fmt.Sprintf("%s/%s/%s", baseurl, l, path)
	})
}
