package test

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/configs"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/entity"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/infra/database"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/infra/usecase"
)

func start() (entity.RepositoryInterface, usecase.RateLimitInputDto) {
	configs.LoadConfig("../")
	client := redis.NewClient(&redis.Options{
		Addr:     configs.Cfg.BdAddress,
		Password: configs.Cfg.BdPassword,
		DB:       0, // use default DB
	})

	repo := database.NewRedisRepository(client)
	register := usecase.NewRegisterLimitUseCase(repo)
	dto := usecase.RateLimitInputDto{IpOrToken: "192.158.1.38"}
	register.Execute(dto)

	return repo, dto
}

func end(repo entity.RepositoryInterface, ipOrToken string) {
	repo.Delete(context.Background(), &entity.IpOrToken{IpOrToken: ipOrToken})
}

type allowResul struct {
	key   int
	resul bool
}

func TestGivenIp_ThenShouldReceivedAllowForAll(t *testing.T) {
	repository, dto := start()
	configs.Cfg.RateLimitIp = 4

	resul := []allowResul{
		{1, true},
		{2, true},
		{3, true},
		{4, true},
	}

	allowUseCase := usecase.NewVerifyAllownTokenOrIpUseCase(repository)
	verifyDto := usecase.VerifyAllownTokenOrIpInputDto{
		IpOrToken: dto.IpOrToken,
		Token:     false,
	}
	for i := 0; i <= len(resul)-1; i++ {
		resulUseCase, _ := allowUseCase.Execute(verifyDto)
		if resulUseCase.Allow != resul[i].resul {
			t.Errorf("Expected %t but got %t", resul[i].resul, resulUseCase.Allow)
		}

	}

	end(repository, dto.IpOrToken)

}

func TestGivenIp_ThenShouldReceivedBlocked(t *testing.T) {
	repository, dto := start()
	configs.Cfg.RateLimitIp = 2

	resul := []allowResul{
		{1, true},
		{2, true},
		{3, false},
		{4, false},
	}

	allowUseCase := usecase.NewVerifyAllownTokenOrIpUseCase(repository)
	verifyDto := usecase.VerifyAllownTokenOrIpInputDto{
		IpOrToken: dto.IpOrToken,
		Token:     false,
	}
	for i := 0; i <= len(resul)-1; i++ {
		resulUseCase, _ := allowUseCase.Execute(verifyDto)
		if resulUseCase.Allow != resul[i].resul {
			t.Errorf("Expected %t but got %t", resul[i].resul, resulUseCase.Allow)
		}

	}

	end(repository, dto.IpOrToken)

}

func TestGivenIpAndToken_ThenShouldReceivedAllowBecauseToken(t *testing.T) {
	repository, dto := start()
	configs.Cfg.RateLimitIp = 2
	configs.Cfg.RateLimitToken = 4

	resul := []allowResul{
		{1, true},
		{2, true},
		{3, true},
		{4, true},
	}

	allowUseCase := usecase.NewVerifyAllownTokenOrIpUseCase(repository)
	verifyDto := usecase.VerifyAllownTokenOrIpInputDto{
		IpOrToken: dto.IpOrToken,
		Token:     true,
	}
	for i := 0; i <= len(resul)-1; i++ {
		resulUseCase, _ := allowUseCase.Execute(verifyDto)
		if resulUseCase.Allow != resul[i].resul {
			t.Errorf("Expected %t but got %t", resul[i].resul, resulUseCase.Allow)
		}

	}

	end(repository, dto.IpOrToken)

}

func TestGivenIpAndToken_ThenShouldReceivedBlockedBecauseToken(t *testing.T) {
	repository, dto := start()
	configs.Cfg.RateLimitIp = 2
	configs.Cfg.RateLimitToken = 4

	resul := []allowResul{
		{1, true},
		{2, true},
		{3, true},
		{4, true},
		{5, false},
	}

	allowUseCase := usecase.NewVerifyAllownTokenOrIpUseCase(repository)
	verifyDto := usecase.VerifyAllownTokenOrIpInputDto{
		IpOrToken: dto.IpOrToken,
		Token:     true,
	}
	for i := 0; i <= len(resul)-1; i++ {
		resulUseCase, _ := allowUseCase.Execute(verifyDto)
		if resulUseCase.Allow != resul[i].resul {
			t.Errorf("Expected %t but got %t", resul[i].resul, resulUseCase.Allow)
		}

	}

	end(repository, dto.IpOrToken)

}

func TestGivenIp_ThenShouldReceivedBlockedAndAllowAfterTime(t *testing.T) {
	repository, dto := start()
	configs.Cfg.RateLimitIp = 2
	configs.Cfg.RateLimitQtdBlock = 2

	resul := []allowResul{
		{1, true},
		{2, true},
		{3, false},
		{4, true},
	}

	allowUseCase := usecase.NewVerifyAllownTokenOrIpUseCase(repository)
	verifyDto := usecase.VerifyAllownTokenOrIpInputDto{
		IpOrToken: dto.IpOrToken,
		Token:     false,
	}
	for i := 0; i <= len(resul)-1; i++ {
		if i == len(resul)-1 {
			time.Sleep(2 * time.Second)
		}
		resulUseCase, _ := allowUseCase.Execute(verifyDto)
		if resulUseCase.Allow != resul[i].resul {
			t.Errorf("%d Expected %t but got %t", i, resul[i].resul, resulUseCase.Allow)
		}

	}

	end(repository, dto.IpOrToken)

}

func TestGivenIp_ThenShouldReceivedBlockedAndBlockedAfterTime(t *testing.T) {
	repository, dto := start()
	configs.Cfg.RateLimitIp = 2
	configs.Cfg.RateLimitQtdBlock = 2

	resul := []allowResul{
		{1, true},
		{2, true},
		{3, false},
		{4, false},
	}

	allowUseCase := usecase.NewVerifyAllownTokenOrIpUseCase(repository)
	verifyDto := usecase.VerifyAllownTokenOrIpInputDto{
		IpOrToken: dto.IpOrToken,
		Token:     false,
	}
	for i := 0; i <= len(resul)-1; i++ {
		if i == len(resul)-1 {
			time.Sleep(1 * time.Second)
		}
		resulUseCase, _ := allowUseCase.Execute(verifyDto)
		if resulUseCase.Allow != resul[i].resul {
			t.Errorf("%d Expected %t but got %t", i, resul[i].resul, resulUseCase.Allow)
		}

	}

	end(repository, dto.IpOrToken)

}
