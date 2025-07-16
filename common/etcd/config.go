package etcd

type EtcdConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
