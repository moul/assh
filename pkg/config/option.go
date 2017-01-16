package config

import "fmt"

// Option is an host option
type Option struct {
	Name  string
	Value string
}

// OptionsList is a list of options
type OptionsList []Option

func (o *Option) String() string {
	return fmt.Sprintf("%s=%s", o.Name, o.Value)
}

// Get returns the option value matching the name or "" if the key is not found
func (ol *OptionsList) Get(name string) string {
	for _, opt := range *ol {
		if opt.Name == name {
			return opt.Value
		}
	}
	return ""
}

// ToStringList returns a list of string with the following format: `key=value`
func (ol *OptionsList) ToStringList() []string {
	list := []string{}
	for _, opt := range *ol {
		list = append(list, opt.String())
	}
	return list
}

// Remove removes an option from the list based on its key
func (ol *OptionsList) Remove(key string) {
	for i, opt := range *ol {
		if opt.Name == key {
			*ol = append((*ol)[:i], (*ol)[i+1:]...)
		}
	}
}
