package usecase

import (
	"context"

	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/entity"
)

type RateLimitInputDto struct {
	IpOrToken string `json:"ipOrToken"`
}

type RateLimitOutputDto struct {
	Message string `json:"message"`
}

type registerLimitUseCase struct {
	Repository entity.RepositoryInterface
}

func NewRegisterLimitUseCase(r entity.RepositoryInterface) *registerLimitUseCase {
	return &registerLimitUseCase{Repository: r}
}

func (r *registerLimitUseCase) Execute(dto RateLimitInputDto) (*RateLimitOutputDto, error) {

	ipOrToken, err := entity.NewIpOrToken(dto.IpOrToken)
	if err != nil {
		return nil, err
	}

	err = r.Repository.Register(context.Background(), ipOrToken)
	if err != nil {
		return nil, err
	}

	return &RateLimitOutputDto{Message: "successful"}, nil
}
