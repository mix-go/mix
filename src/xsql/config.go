package xsql

type Option interface {
	Apply(*Config)
}

type Config struct {
	// 默认为: INSERT INTO
	InsertKey string

	// 默认为: ? , oracle 可配置为 :%d
	Placeholder string
}

func (c *Config) Apply(config *Config) {
	if config != c {
		*config = *c
	}
}
