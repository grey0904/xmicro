package dao

import (
	"context"
	"xmicro/internal/app/user/dto"
	"xmicro/internal/core/repo"
)

type UserDao struct {
	repo *repo.Manager
}

func (d *UserDao) SaveAccount(ctx context.Context, ac *dto.Account) error {
	var err error
	if err != nil {
		return err
	}
	return nil
}

func NewAccountDao(m *repo.Manager) *UserDao {
	return &UserDao{
		repo: m,
	}
}
