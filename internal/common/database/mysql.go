package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"xmicro/internal/common/config"
	"xmicro/internal/common/logs"
)

type MysqlManager struct {
	DB *sql.DB
}

func NewMysql() *MysqlManager {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s",
		config.Conf.Database.MysqlConf.Username,
		config.Conf.Database.MysqlConf.Password,
		config.Conf.Database.MysqlConf.Host,
		config.Conf.Database.MysqlConf.Database))
	if err != nil {
		logs.Fatal("Error opening MySQL connection: %v\n", err)
		return nil
	}

	// 设置最大连接数和最大空闲连接数
	db.SetMaxOpenConns(config.Conf.Database.MysqlConf.MaxOpenConns)
	db.SetMaxIdleConns(config.Conf.Database.MysqlConf.MaxIdleConns)

	// 检查连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		logs.Fatal("Error pinging MySQL server: %v\n", err)
		return nil
	}

	return &MysqlManager{DB: db}
}

func (m *MysqlManager) Close() {
	if m.DB != nil {
		if err := m.DB.Close(); err != nil {
			logs.Error("Mysql close err: %v", err)
		}
	}
}

func (m *MysqlManager) Set(ctx context.Context, key, value string, expire time.Duration) error {
	// 这里填写你的 Mysql SET 操作逻辑，根据你的需求进行修改
	// 例如：_, err := m.DB.ExecContext(ctx, "INSERT INTO my_table (key, value) VALUES (?, ?) ON DUPLICATE KEY UPDATE value = ?", key, value, value)
	// 这只是一个示例，请根据你的实际情况进行修改
	return nil
}
