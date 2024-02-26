package controller

import (
	"net/http"

	rip "github.com/vikram1565/request-ip"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/entity"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/infra/usecase"
)

type AllowTokenOrIpController struct {
	Repository entity.RepositoryInterface
}

func NewAllownTokenOrIpController(db entity.RepositoryInterface) *RegisterLimitController {
	return &RegisterLimitController{Repository: db}
}

func (a *RegisterLimitController) VerifyAllow(w http.ResponseWriter, r *http.Request) bool {

	var dto usecase.VerifyAllownTokenOrIpInputDto
	dto.IpOrToken = rip.GetClientIP(r)
	dto.Token = false

	//tratativa para localhost
	if dto.IpOrToken == "::1" {
		dto.IpOrToken = "127.0.0.1"
	}

	token := r.Header.Get("API_KEY")
	if token != "" {
		dto.IpOrToken = token
		dto.Token = true
	}

	usecase := usecase.NewVerifyAllownTokenOrIpUseCase(a.Repository)
	output, err := usecase.Execute(dto)

	if err != nil {
		return false
	}

	return output.Allow
}
