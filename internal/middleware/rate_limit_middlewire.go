package middleware

import (
	"net/http"

	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/entity"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/infra/controller"
)

type AllowMiddleware struct {
	Repository entity.RepositoryInterface
}

func NewRateLimitMiddleware(db entity.RepositoryInterface) *AllowMiddleware {
	return &AllowMiddleware{Repository: db}
}

func (a *AllowMiddleware) RateLimit(next http.Handler) http.Handler {

	controller := controller.NewAllownTokenOrIpController(a.Repository)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allow := controller.VerifyAllow(w, r)
		if !allow {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
