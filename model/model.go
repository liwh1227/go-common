package model

import "fmt"

// Gorm 框架的相关配置
type Gorm struct {
	// 日志打印级别
	Debug bool `json:"debug" yaml:"Debug"`
	// 数据库类型：例如mysql
	DBType            string `json:"dbType" yaml:"DBType"`
	MaxLifetime       int    `json:"maxLifetime" yaml:"MaxLifetime"`
	MaxOpenConns      int    `json:"maxOpenConns" yaml:"MaxOpenConns"`
	MaxIdleConns      int    `json:"maxIdleConns" yaml:"MaxIdleConns"`
	EnableAutoMigrate bool   `json:"enableAutoMigrate" yaml:"EnableAutoMigrate"`
	// 是否开启日志打印
	IsLoggerOn bool `json:"isLoggerOn"`
}

// mysql数据库配置
type Mysql struct {
	// ip
	Host string `json:"host" yaml:"Host"`
	// 端口
	Port int `json:"port" yaml:"Port"`
	// mysql cli用户
	User string `json:"user" yaml:"User"`
	// 密码
	Password string `json:"password" yaml:"Password"`
	// 数据库
	DBName string `json:"dbName" yaml:"DBName"`
	// 其他参数
	Parameters string `json:"parameters" yaml:"Parameters"`
}

// DSN 数据库连接串
func (m Mysql) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		m.User, m.Password, m.Host, m.Port, m.DBName, m.Parameters)
}
