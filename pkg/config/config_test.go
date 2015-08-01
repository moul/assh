package config

import "fmt"

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

func ExampleNew() {
	config := New()
	fmt.Println(config)
	// Output: &{map[] {  0}}
}

func ExampleConfig() {
	config := dummyConfig()
	fmt.Println(config.Hosts["toto"])
	fmt.Println(config.Hosts["titi"])
	fmt.Println(config.Defaults)
	// Output:
	// {1.2.3.4  0}
	// {tata moul 23}
	// { root 22}
}
