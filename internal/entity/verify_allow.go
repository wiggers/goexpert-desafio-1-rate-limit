package entity

import "github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/configs"

type verifyAllow struct {
	Allow bool
}

func NewVerifyAllow(qtd int, token bool) *verifyAllow {
	if token {
		return &verifyAllow{Allow: qtd <= configs.Cfg.RateLimitToken}
	}

	return &verifyAllow{Allow: qtd <= configs.Cfg.RateLimitIp}
}
