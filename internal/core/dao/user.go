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
	table := d.repo.Mongo.Db.Collection("account")
	_, err := table.InsertOne(ctx, ac)
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
