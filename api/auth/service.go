package auth

import (
	"api/config"
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	cfg config.Config
}

func NewService(cfg config.Config) *jwtService {
	return &jwtService{cfg}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(s.cfg.JWTConfig.Secret))
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodeToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodeToken, func(t *jwt.Token) (interface{}, error) {	
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte([]byte(s.cfg.JWTConfig.Secret)),nil
	})
	if err != nil {
		return token, err
	}
	return token, nil
}