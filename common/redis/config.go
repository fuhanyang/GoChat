package redis

type RedisConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Password  string `mapstructure:"password"`
	Db        int    `mapstructure:"db"`
	Network   string `mapstructure:"network"`
	MaxIdle   int    `mapstructure:"max_idle"`
	MaxActive int    `mapstructure:"max_active"`
}
