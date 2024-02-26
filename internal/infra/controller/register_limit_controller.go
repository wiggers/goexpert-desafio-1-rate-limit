package controller

import (
	"encoding/json"
	"net/http"

	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/entity"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/infra/usecase"
)

type RegisterLimitController struct {
	Repository entity.RepositoryInterface
}

func NewRegisterLimitController(db entity.RepositoryInterface) *RegisterLimitController {
	return &RegisterLimitController{Repository: db}
}

func (c *RegisterLimitController) Register(w http.ResponseWriter, r *http.Request) {
	var dto usecase.RateLimitInputDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usecase := usecase.NewRegisterLimitUseCase(c.Repository)
	output, err := usecase.Execute(dto)

	json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
