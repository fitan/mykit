package conf

type Conf struct {
	App    App    `yaml:"App"`
	Mysql  Mysql  `yaml:"Mysql"`
	Log    Log    `yaml:"Log"`
	Consul Consul `yaml:"Consul"`
}
type App struct {
	Name string `yaml:"Name"`
	Addr string `yaml:"Addr"`
	Port int    `yaml:"Port"`
}
type Mysql struct {
	DSN string `yaml:"DSN"`
}
type Log struct {
	Dir   string `yaml:"Dir"`
	Level int64  `yaml:"Level"`
}

type Consul struct {
	Addr  string `yaml:"Addr"`
	Token string `yaml:"Token"`
}
