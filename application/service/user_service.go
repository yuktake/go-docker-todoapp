package service

import (
	"github.com/yuktake/todo-webapp/domain/user"
	"github.com/yuktake/todo-webapp/logger"

	"go.uber.org/fx"
)

type User = user.User

type UserService interface {
	CreateUser(user User) (User, error)
	GetUserByID(id string) (User, error)
	GetUsers() ([]User, error)
	UpdateUser(user User) (User, error)
	DeleteUserByID(id string) error
	GetUserByEmail(email string) (User, error)
}

type userService struct {
	repo   user.UserRepository
	logger logger.Logger
}

// `fx.In` で `UserRepository` を自動 DI
type UserServiceParams struct {
	fx.In
	Repo   user.UserRepository
	Logger logger.Logger
}

func NewUserService(params UserServiceParams) UserService {
	return &userService{repo: params.Repo, logger: params.Logger}
}

func (s *userService) CreateUser(user User) (User, error) {
	user, err := s.repo.CreateUser(&user)

	if err != nil {
		s.logger.Error("failed to create user", err)
		return User{}, err
	}

	return user, nil
}

func (s *userService) GetUserByID(id string) (User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		s.logger.Error("failed to get user", err)
		return User{}, err
	}

	return user, nil
}

func (s *userService) GetUsers() ([]User, error) {
	users, err := s.repo.GetUsers()
	if err != nil {
		s.logger.Error("failed to get users", err)
		return nil, err
	}

	return users, nil
}

func (s *userService) UpdateUser(user User) (User, error) {
	user, err := s.repo.UpdateUser(user)

	if err != nil {
		s.logger.Error("failed to update user", err)
		return User{}, err
	}

	return user, nil
}

func (s *userService) DeleteUserByID(id string) error {
	err := s.repo.DeleteUserByID(id)

	if err != nil {
		s.logger.Error("failed to delete user", err)
		return err
	}

	return nil
}

func (s *userService) GetUserByEmail(email string) (User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		s.logger.Error("メールアドレスまたはパスワードが無効です", err)
		return User{}, err
	}

	return user, nil
}
