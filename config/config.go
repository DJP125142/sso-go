package config

type ServerConfig struct {
	Name        string      `mapstructure:"appName"`
	Port        int         `mapstructure:"port"`
	MysqlInfo   MysqlConfig `mapstructure:"mysql"`
	RedisInfo   RedisConfig `mapstructure:"redis"`
	EmailInfo   EmailConfig `mapstructure:"email"`
	LogsAddress string      `mapstructure:"logsAddress"`
	JWTKey      JWTConfig   `mapstructure:"jwt"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type EmailConfig struct {
	Address   string `mapstructure:"address"`
	Host      string `mapstructure:"host"`
	SendName  string `mapstructure:"sendName"`
	SendEmail string `mapstructure:"sendEmail"`
	Password  string `mapstructure:"password"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}
