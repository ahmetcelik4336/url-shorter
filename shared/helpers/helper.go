package helpers

import (
	"fmt"
	"os"

	beego "github.com/beego/beego/v2/server/web"
)

func Init() {
	beego.AddFuncMap("baseurl", func(lang interface{}, path string) string {
		l, _ := lang.(string)
		baseurl := os.Getenv("BASEURL")
		return fmt.Sprintf("%s/%s/%s", baseurl, l, path)
	})

}

func Baseurl(lang interface{}, path string) string {
	l, _ := lang.(string)
	baseurl := os.Getenv("BASEURL")
	return fmt.Sprintf("%s/%s/%s", baseurl, l, path)
}
