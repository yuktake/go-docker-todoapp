package user

import (
	"context"

	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type UserRepository interface {
	CreateUser(user *User) (User, error)
	GetUsers() ([]User, error)
	GetUserByID(id string) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUserByID(id string) error
	GetUserByEmail(email string) (User, error)
}

// 小文字始まりの構造体は非公開
type userRepository struct {
	DB *bun.DB
}

// `fx.In` で `bun.DB` を自動 DI
type userRepositoryParams struct {
	fx.In
	DB *bun.DB
}

func NewUserRepository(params userRepositoryParams) UserRepository {
	return &userRepository{DB: params.DB}
}

// User作成
func (r *userRepository) CreateUser(user *User) (User, error) {
	ctx := context.Background()

	_, err := r.DB.NewInsert().Model(user).Returning("*").Exec(ctx)
	if err != nil {
		return User{}, err
	}

	return *user, nil
}

// User取得
func (r *userRepository) GetUserByID(id string) (User, error) {
	var user User
	ctx := context.Background()

	err := r.DB.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *userRepository) GetUsers() ([]User, error) {
	// リポジトリからデータを取得
	var users []User
	ctx := context.Background()

	err := r.DB.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) UpdateUser(user User) (User, error) {
	ctx := context.Background()

	_, err := r.DB.NewUpdate().Model(&user).Where("id = ?", user.ID).Exec(ctx)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *userRepository) DeleteUserByID(id string) error {
	ctx := context.Background()

	_, err := r.DB.NewDelete().Model((*User)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserByEmail(email string) (User, error) {
	var user User
	ctx := context.Background()

	err := r.DB.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
