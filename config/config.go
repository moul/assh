package config

type Host struct {
}

type Config struct {
	Hosts []Host
}

func Open() (*Config, error) {
	return &Config{}, nil
}
