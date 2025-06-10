package initialize

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"server/global"
)

// GormMysql 初始化Mysql数据库
func GormMysql() *gorm.DB {
	m := global.CONFIG.Mysql
	if m.DbName == "" {
		return nil
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		m.Username,
		m.Password,
		m.Path,
		m.DbName,
		m.Config,
	)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{}); err != nil {
		return nil
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
