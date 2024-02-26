package entity

import "context"

type RepositoryInterface interface {
	Add(ctx context.Context, key *IpOrToken) (int, error)
	Block(ctx context.Context, ipOrToken *IpOrToken, qtd int, timeBlock int) error
	Find(ctx context.Context, key *IpOrToken) int
	Register(ctx context.Context, key *IpOrToken) error
	VerifyWithRegister(ctx context.Context, key *IpOrToken) bool
	Delete(ctx context.Context, ipOrToken *IpOrToken) error
}
