package usecase

import (
	"context"

	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/configs"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/entity"
)

type VerifyAllownTokenOrIpInputDto struct {
	IpOrToken string
	Token     bool
}

type VerifyAllownTokenOrIpOuputDto struct {
	Allow bool
}

type VerifyAllownTokenOrIpUseCase struct {
	Repository entity.RepositoryInterface
}

func NewVerifyAllownTokenOrIpUseCase(r entity.RepositoryInterface) *VerifyAllownTokenOrIpUseCase {
	return &VerifyAllownTokenOrIpUseCase{Repository: r}
}

func (r *VerifyAllownTokenOrIpUseCase) Execute(dto VerifyAllownTokenOrIpInputDto) (*VerifyAllownTokenOrIpOuputDto, error) {

	ctx := context.Background()
	ipOrToken, err := entity.NewIpOrToken(dto.IpOrToken)
	if err != nil {
		return nil, err
	}

	exist := r.Repository.VerifyWithRegister(ctx, ipOrToken)
	if !exist {
		return &VerifyAllownTokenOrIpOuputDto{Allow: true}, nil
	}

	qtd, err := r.Repository.Add(ctx, ipOrToken)
	if err != nil {
		return nil, err
	}

	verify := entity.NewVerifyAllow(qtd, dto.Token)

	if !verify.Allow {
		r.Repository.Block(ctx, ipOrToken, qtd, configs.Cfg.RateLimitQtdBlock)
	}

	return &VerifyAllownTokenOrIpOuputDto{Allow: verify.Allow}, nil
}
