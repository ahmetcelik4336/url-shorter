package middleware

import (
	"github.com/beego/beego/v2/server/web/context" // Beego context
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

var limiter *redis_rate.Limiter

func InitRateLimiter(rdb *redis.Client) {
	limiter = redis_rate.NewLimiter(rdb)
}

func CheckRateLimit(ctx *context.Context, action string, duration int) bool {
	// Kullanıcıyı IP adresine göre ayırt ediyoruz
	userIP := ctx.Input.IP()

	// Örn: Dakikada 10 isteğe izin ver
	res, err := limiter.Allow(ctx.Request.Context(), action+"ratelimit:"+userIP, redis_rate.PerMinute(duration))
	if err != nil {
		return true // Redis hatasında sistemi kilitlememek için izin ver (fail-open)
	}

	if res.Allowed == 0 {
		ctx.Output.SetStatus(429)
		ctx.Output.Body([]byte("Çok fazla istek gönderdiniz. Lütfen biraz bekleyin."))
		return false
	}
	return true
}
