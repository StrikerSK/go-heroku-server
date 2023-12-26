package userPorts

import userDomains "go-heroku-server/api/user/domain"

type ITokenService interface {
	CreateToken(user userDomains.User) (string, error)
	ParseToken(signedToken string) (*userDomains.UserClaims, error)
}
