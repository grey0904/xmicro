package repo

import "xmicro/internal/common/database"

type Manager struct {
	Redis *database.RedisManager
	Mysql *database.MysqlManager
}

func (m *Manager) Close() {
	if m.Redis != nil {
		m.Redis.Close()
	}
	if m.Redis != nil {
		m.Mysql.Close()
	}
}

func New() *Manager {
	return &Manager{
		Redis: database.NewRedis(),
		Mysql: database.NewMysql(),
	}
}
