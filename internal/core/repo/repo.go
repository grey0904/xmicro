package repo

import "xmicro/internal/common/database"

type Manager struct {
	Mongo *database.MongoManager
	Redis *database.RedisManager
	Mysql *database.MysqlManager
}

func (m *Manager) Close() {
	if m.Mongo != nil {
		m.Mongo.Close()
	}
	if m.Redis != nil {
		m.Redis.Close()
	}
	if m.Redis != nil {
		m.Mysql.Close()
	}
}

func New() *Manager {
	return &Manager{
		Mongo: database.NewMongo(),
		Redis: database.NewRedis(),
		Mysql: database.NewMysql(),
	}
}
