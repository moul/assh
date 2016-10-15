package service

import (
	"path/filepath"
	"testing"

	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/lookup"
	"github.com/docker/libcompose/yaml"
	shlex "github.com/flynn/go-shlex"
	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
	exp := []string{"sh", "-c", "exec /opt/bin/flanneld -logtostderr=true -iface=${NODE_IP}"}
	cmd, err := shlex.Split("sh -c 'exec /opt/bin/flanneld -logtostderr=true -iface=${NODE_IP}'")
	assert.Nil(t, err)
	assert.Equal(t, exp, cmd)
}

func TestParseBindsAndVolumes(t *testing.T) {
	ctx := &ctx.Context{}
	ctx.ComposeFiles = []string{"foo/docker-compose.yml"}
	ctx.ResourceLookup = &lookup.FileResourceLookup{}

	abs, err := filepath.Abs(".")
	assert.Nil(t, err)
	cfg, hostCfg, err := Convert(&config.ServiceConfig{
		Volumes: &yaml.Volumes{
			Volumes: []*yaml.Volume{
				{
					Destination: "/foo",
				},
				{
					Source:      "/home",
					Destination: "/home",
				},
				{
					Destination: "/bar/baz",
				},
				{
					Source:      ".",
					Destination: "/home",
				},
				{
					Source:      "/usr/lib",
					Destination: "/usr/lib",
					AccessMode:  "ro",
				},
			},
		},
	}, ctx.Context, nil)
	assert.Nil(t, err)
	assert.Equal(t, map[string]struct{}{"/foo": {}, "/bar/baz": {}}, cfg.Volumes)
	assert.Equal(t, []string{"/home:/home", abs + "/foo:/home", "/usr/lib:/usr/lib:ro"}, hostCfg.Binds)
}

func TestParseLabels(t *testing.T) {
	ctx := &ctx.Context{}
	ctx.ComposeFiles = []string{"foo/docker-compose.yml"}
	ctx.ResourceLookup = &lookup.FileResourceLookup{}
	bashCmd := "bash"
	fooLabel := "foo.label"
	fooLabelValue := "service.config.value"
	sc := &config.ServiceConfig{
		Entrypoint: yaml.Command([]string{bashCmd}),
		Labels:     yaml.SliceorMap{fooLabel: "service.config.value"},
	}
	cfg, _, err := Convert(sc, ctx.Context, nil)
	assert.Nil(t, err)

	cfg.Labels[fooLabel] = "FUN"
	cfg.Entrypoint[0] = "less"

	assert.Equal(t, fooLabelValue, sc.Labels[fooLabel])
	assert.Equal(t, "FUN", cfg.Labels[fooLabel])

	assert.Equal(t, yaml.Command{bashCmd}, sc.Entrypoint)
	assert.Equal(t, []string{"less"}, []string(cfg.Entrypoint))
}
