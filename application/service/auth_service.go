package service

import (
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/yuktake/todo-webapp/domain/auth"
	"github.com/yuktake/todo-webapp/domain/user"
	"github.com/yuktake/todo-webapp/logger"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
)

type AuthService interface {
	CreateToken(user user.User) (string, error)
}

type authService struct {
	logger logger.Logger
}

// `fx.In` で `UserRepository` を自動 DI
type authServiceParams struct {
	fx.In
	Logger logger.Logger
}

func NewAuthService(params authServiceParams) AuthService {
	return &authService{logger: params.Logger}
}

func (s *authService) CreateToken(user user.User) (string, error) {
	// JWTクレームの設定
	claims := &auth.JwtCustomClaims{
		Email: user.Email,
		Name:  user.Name,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
		},
	}

	// クレームを持つトークンを生成
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	// トークンを署名
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		s.logger.Error("トークンの署名エラー", err)
		return "", err
	}

	// トークンを返す
	return t, nil
}

// 暗号化 (hash)
func PasswordEncrypt(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashPassword), err
}

// 暗号化パスワードと比較
func CheckHashPassword(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
