package yaml_test

import (
	. "github.com/moul/advanced-ssh-config/vendor/gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})
