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

func BaseurlWithoutSlug(path string) string {
	baseurl := os.Getenv("BASEURL")
	return fmt.Sprintf("%s/%s", baseurl, path)
}

func GetShortUrl(alias, password, shortCode, type_ string) string {
	if alias == "" {
		alias = shortCode
	}

	if password != "" {
		password = "/" + password
	} else {
		password = ""
	}

	return BaseurlWithoutSlug(type_ + "/" + alias + password)
}
