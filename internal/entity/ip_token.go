package entity

import "errors"

type IpOrToken struct {
	IpOrToken string
}

func NewIpOrToken(ipOrToken string) (*IpOrToken, error) {
	if len(ipOrToken) > 50 {
		return &IpOrToken{}, errors.New("max Length for ip or token is 50")
	}

	return &IpOrToken{IpOrToken: ipOrToken}, nil
}
