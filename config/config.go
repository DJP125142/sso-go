package config

type ServerConfig struct {
	Name        string      `mapstructure:"appName"`
	Port        int         `mapstructure:"port"`
	Mysqlinfo   MysqlConfig `mapstructure:"mysql"`
	RedisInfo   RedisConfig `mapstructure:"redis"`
	LogsAddress string      `mapstructure:"logsAddress"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
