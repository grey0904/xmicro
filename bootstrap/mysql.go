package bootstrap

import (
	"strconv"
)

// 将val值转换为php-redis能取的格式
func CovertToTpRedisVal(val []byte) string {
	slen := strconv.Itoa(len(val))
	s := "s:" + slen + ":\"" + string(val) + "\";"
	return s
}

//func InitMysql() error {
//	dbinfo := config.Conf.Mysql
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
//		dbinfo.Username, dbinfo.Password, dbinfo.Host, dbinfo.Port, dbinfo.Db, dbinfo.Charset)
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		log.Logger.Error("Mysql 初始化错误:", err)
//		return err
//	}
//	app.App.DB = db
//	return nil
//}
