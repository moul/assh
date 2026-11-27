package utils

import "fmt"

// ScalarOrStringSlice represents a YAML scalar or a sequence of strings.
type ScalarOrStringSlice []string

// UnmarshalYAML preserves a single scalar as one entry, including commas and spaces.
func (s *ScalarOrStringSlice) UnmarshalYAML(decode func(any) error) error {
	var single string
	if err := decode(&single); err == nil {
		*s = ScalarOrStringSlice{single}
		return nil
	}

	var items []any
	if err := decode(&items); err != nil {
		return err
	}
	if len(items) == 0 {
		*s = nil
		return nil
	}
	values := make(ScalarOrStringSlice, len(items))
	for i, item := range items {
		value, ok := item.(string)
		if !ok {
			return fmt.Errorf("expected a string at list position %d, got %T", i+1, item)
		}
		values[i] = value
	}
	*s = values
	return nil
}
