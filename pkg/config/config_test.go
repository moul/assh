package config

import "testing"

func dummyConfig() *Config {
	config := New()
	config.Hosts["toto"] = Host{
		Host: "1.2.3.4",
	}
	config.Hosts["titi"] = Host{
		Host: "tata",
		Port: 23,
		User: "moul",
	}
	config.Defaults = Host{
		Port: 22,
		User: "root",
	}
	return config
}

func TestNew(t *testing.T) {
	config := New()

	if len(config.Hosts) != 0 {
		t.Fatalf("Expected len(config.Hosts)=0 got %d", len(config.Hosts))
	}

	if config.Defaults.Port != 0 {
		t.Fatalf("Expected config.Defaults.Port=0 got %d", config.Defaults.Port)
	}
}

func TestConfig(t *testing.T) {
	config := dummyConfig()

	if len(config.Hosts) != 2 {
		t.Fatalf("Expected len(config.Hosts)=2 got %d", len(config.Hosts))
	}

	if config.Hosts["toto"].Host != "1.2.3.4" {
		t.Fatalf("Expected config.Hosts[\"toto\"].Host=1.2.3.4 got %s", config.Hosts["toto"].Host)
	}

	if config.Defaults.Port != 22 {
		t.Fatalf("Expected config.Defaults.Port=22 got %d", config.Defaults.Port)
	}
}
