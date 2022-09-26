package conf

type Conf struct {
	App   App   `yaml:"App"`
	Mysql Mysql `yaml:"Mysql"`
	Log   Log   `yaml:"Log"`
}
type App struct {
	Name string `yaml:"Name"`
	Addr string `yaml:"Addr"`
}
type Mysql struct {
	DSN string `yaml:"DSN"`
}
type Log struct {
	Dir   string `yaml:"Dir"`
	Level int64
}
