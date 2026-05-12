package helpers

import (
	"fmt"
	"os"

	beego "github.com/beego/beego/v2/server/web"
)

func Init() {
	beego.AddFuncMap("baseurl", func(lang interface{}, path string) string {
		l, ok := lang.(string)
		if !ok || l == "" {
			val := os.Getenv("BASEURL")
			if val == "" {
				l = "en"
			} else {
				l = val
			}
		}
		baseurl := os.Getenv("BASEURL")

		return fmt.Sprintf("%s/%s/%s", baseurl, l, path)
	})

}
