package controllers

import (
	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/server/web/captcha"
)

// Global bir captcha nesnesi oluşturuyoruz
var Cpt *captcha.Captcha

func init() {
	store := cache.NewMemoryCache()
	Cpt = captcha.NewCaptcha("/captcha/", store)
	Cpt.ChallengeNums = 4 // Kaç haneli olacak
	Cpt.StdWidth = 120    // Genişlik
	Cpt.StdHeight = 40    // Yükseklik
}
