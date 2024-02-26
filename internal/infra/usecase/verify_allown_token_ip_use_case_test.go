package usecase

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/configs"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/entity"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/infra/database"
)

func start() entity.RepositoryInterface {
	//configs.LoadConfig("../../..")
	client := redis.NewClient(&redis.Options{
		Addr:     configs.Cfg.BdAddress,
		Password: configs.Cfg.BdPassword,
		DB:       0, // use default DB
	})

	repo := database.NewRedisRepository(client)
	return repo
}

type allowResul struct {
	key   int
	resul bool
}

func TestGivenIpWithSomeRequest_ThenShouldReceivedAllowForAll(t *testing.T) {
	configs.LoadConfig("../..")
	repository := start()
	register := NewRegisterLimitUseCase(repository)
	dto := RateLimitInputDto{IpOrToken: "192.158.1.38"}
	register.Execute(dto)

	t.Setenv("RATE_LIMIT_IP", "2")

	resul := []allowResul{
		{1, true},
		{2, true},
	}

	allowUseCase := NewVerifyAllownTokenOrIpUseCase(repository)
	for i := 0; i < 2; i++ {
		verifyDto := VerifyAllownTokenOrIpInputDto{
			IpOrToken: dto.IpOrToken,
			Token:     false,
		}
		resulUseCase, _ := allowUseCase.Execute(verifyDto)
		if resulUseCase.Allow != resul[i].resul {
			t.Errorf("Expected %t but got %t", resul[i].resul, resulUseCase.Allow)
		}

	}

}
