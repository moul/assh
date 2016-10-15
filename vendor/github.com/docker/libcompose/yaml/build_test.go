package yaml

import (
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
)

func TestMarshalBuild(t *testing.T) {
	builds := []struct {
		build    Build
		expected string
	}{
		{
			expected: `{}
`,
		},
		{
			build: Build{
				Context: ".",
			},
			expected: `context: .
`,
		},
		{
			build: Build{
				Context:    ".",
				Dockerfile: "alternate",
			},
			expected: `context: .
dockerfile: alternate
`,
		},
		{
			build: Build{
				Context:    ".",
				Dockerfile: "alternate",
				Args: map[string]string{
					"buildno": "1",
					"user":    "vincent",
				},
			},
			expected: `args:
  buildno: "1"
  user: vincent
context: .
dockerfile: alternate
`,
		},
	}
	for _, build := range builds {
		bytes, err := yaml.Marshal(build.build)
		assert.Nil(t, err)
		assert.Equal(t, build.expected, string(bytes), "should be equal")
	}
}

func TestUnmarshalBuild(t *testing.T) {
	builds := []struct {
		yaml     string
		expected *Build
	}{
		{
			yaml: `.`,
			expected: &Build{
				Context: ".",
			},
		},
		{
			yaml: `context: .`,
			expected: &Build{
				Context: ".",
			},
		},
		{
			yaml: `context: .
dockerfile: alternate`,
			expected: &Build{
				Context:    ".",
				Dockerfile: "alternate",
			},
		},
		{
			yaml: `context: .
dockerfile: alternate
args:
  buildno: 1
  user: vincent`,
			expected: &Build{
				Context:    ".",
				Dockerfile: "alternate",
				Args: map[string]string{
					"buildno": "1",
					"user":    "vincent",
				},
			},
		},
		{
			yaml: `context: .
args:
  - buildno
  - user`,
			expected: &Build{
				Context: ".",
				Args: map[string]string{
					"buildno": "\x00",
					"user":    "\x00",
				},
			},
		},
		{
			yaml: `context: .
args:
  - buildno=1
  - user=vincent`,
			expected: &Build{
				Context: ".",
				Args: map[string]string{
					"buildno": "1",
					"user":    "vincent",
				},
			},
		},
	}
	for _, build := range builds {
		actual := &Build{}
		err := yaml.Unmarshal([]byte(build.yaml), actual)
		assert.Nil(t, err)
		assert.Equal(t, build.expected, actual, "should be equal")
	}
}
