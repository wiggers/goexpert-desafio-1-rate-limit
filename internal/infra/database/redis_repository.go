package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wiggers/goexpert/desafio-tecnico/1-rate-limit/internal/entity"
)

type Repository struct {
	repository *redis.Client
}

func NewRedisRepository(db *redis.Client) *Repository {
	return &Repository{repository: db}
}

func (r *Repository) Add(ctx context.Context, ipOrToken *entity.IpOrToken) (int, error) {
	key := ipOrToken.IpOrToken
	resul, err := r.repository.Get(ctx, key).Int()
	if err == redis.Nil {
		err = r.repository.Set(ctx, key, 1, 10*time.Second).Err()
		if err != nil {
			return 0, err
		}
		return 1, nil
	}

	err = r.repository.Set(ctx, key, resul+1, 10*time.Second).Err()
	if err != nil {
		return 0, err
	}

	return resul + 1, nil
}

func (r *Repository) Block(ctx context.Context, ipOrToken *entity.IpOrToken, qtd int, timeBlock int) error {
	key := ipOrToken.IpOrToken
	return r.repository.Set(ctx, key, qtd, time.Duration(timeBlock)*time.Second).Err()
}

func (r *Repository) Find(ctx context.Context, ipOrToken *entity.IpOrToken) int {
	key := ipOrToken.IpOrToken
	resul, err := r.repository.Get(ctx, key).Int()

	if err == redis.Nil {
		return 0
	}

	return resul
}

func (r *Repository) Register(ctx context.Context, ipOrToken *entity.IpOrToken) error {
	key := ipOrToken.IpOrToken
	return r.repository.Set(ctx, key+"_REGISTER", true, 0).Err()
}

func (r *Repository) VerifyWithRegister(ctx context.Context, ipOrToken *entity.IpOrToken) bool {
	key := ipOrToken.IpOrToken
	resul, _ := r.repository.Get(ctx, key+"_REGISTER").Bool()
	return resul
}

func (r *Repository) Delete(ctx context.Context, ipOrToken *entity.IpOrToken) error {
	key := ipOrToken.IpOrToken
	resul := r.repository.Del(ctx, key+"_REGISTER").Err()
	if resul != nil {
		return resul
	}

	resul = r.repository.Del(ctx, key).Err()
	if resul != nil {
		return resul
	}

	return nil
}
