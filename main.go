package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/configs"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/infra/controller"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/infra/database"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/middleware"
)

func main() {
	configs.LoadConfig(".")
	client := redis.NewClient(&redis.Options{
		Addr:     configs.Cfg.BdAddress,
		Password: configs.Cfg.BdPassword,
		DB:       0, // use default DB
	})

	repo := database.NewRedisRepository(client)
	register := controller.NewRegisterLimitController(repo)
	rate := middleware.NewRateLimitMiddleware(repo)

	r := chi.NewRouter()
	r.Post("/register", register.Register)
	r.Route("/message", func(r chi.Router) {
		r.Use(rate.RateLimit)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode("Endpoint para teste de rateLimit")
		})
	})

	http.ListenAndServe(":"+configs.Cfg.WebServerPort, r)
}
