package flexyaml

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

func MakeFlexible(in []byte) ([]byte, error) {
	lines := []string{}
	for _, line := range strings.Split(string(in), "\n") {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			line = fmt.Sprintf("%s:%s", strings.ToLower(parts[0]), parts[1])
		}
		lines = append(lines, line)
	}
	return []byte(strings.Join(lines, "\n")), nil
}

func Unmarshal(in []byte, out interface{}) (err error) {
	flex, err := MakeFlexible(in)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(flex, out)
}
