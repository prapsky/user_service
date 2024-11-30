package service_auth

import (
	"crypto/rsa"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/prapsky/user_service/common/errors"
	"github.com/prapsky/user_service/entity"
)

type JwtAuthService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type JwtAuthServiceOptions struct {
	PrivateKey string
	PublicKey  string
}

type AuthService interface {
	CreateToken(user *entity.User) (string, error)
	ValidateToken(token string) (uint64, error)
}

func NewJwtAuthService(opts JwtAuthServiceOptions) *JwtAuthService {
	pem, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(opts.PrivateKey))
	if err != nil {
		panic(err)
	}

	cert, err := jwt.ParseRSAPublicKeyFromPEM([]byte(opts.PublicKey))
	if err != nil {
		panic(err)
	}

	return &JwtAuthService{
		privateKey: pem,
		publicKey:  cert,
	}
}

func (s *JwtAuthService) CreateToken(user *entity.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  strconv.FormatInt(int64(user.ID), 10),
		"name": user.Username,
		"iat":  now.Unix(),
		"exp":  now.Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *JwtAuthService) ValidateToken(token string) (uint64, error) {
	jt, err := s.parseToken(token)
	if err != nil {
		return 0, err
	}

	if jt == nil {
		return 0, errors.ErrInvalidToken
	}

	if !jt.Valid {
		return 0, errors.ErrInvalidToken
	}

	return s.getClaim(jt)
}

func (s *JwtAuthService) parseToken(token string) (*jwt.Token, error) {
	jt, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.ErrUnexpectedSigning
		}

		return s.publicKey, nil
	})

	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	return jt, nil
}

func (s *JwtAuthService) getClaim(jt *jwt.Token) (uint64, error) {
	id := uint64(0)
	claims, ok := jt.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.ErrInvalidToken
	}

	id, err := strconv.ParseUint(sub, 10, 64)
	if err != nil {
		return 0, errors.ErrInvalidToken
	}

	return id, nil
}
