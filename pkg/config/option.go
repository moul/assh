package config

import "fmt"

type Option struct {
	Name  string
	Value string
}

type OptionsList []Option

func (o *Option) String() string {
	return fmt.Sprintf("%s=%s", o.Name, o.Value)
}

func (ol *OptionsList) ToStringList() []string {
	list := []string{}
	for _, opt := range *ol {
		list = append(list, opt.String())
	}
	return list
}

func (ol *OptionsList) Remove(key string) {
	for i, opt := range *ol {
		if opt.Name == key {
			*ol = append((*ol)[:i], (*ol)[i+1:]...)
		}
	}
}
