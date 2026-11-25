package commands

import (
	"go.uber.org/zap"
	"moul.io/assh/v2/pkg/config"
	"moul.io/assh/v2/pkg/hooks"
)

type configWriteHookArgs struct {
	SSHConfigPath string
}

func rewriteSSHConfig(conf *config.Config) error {
	configHooks := conf.Defaults.Hooks
	args := configWriteHookArgs{SSHConfigPath: conf.SSHConfigPath()}
	if configHooks != nil {
		invokeConfigWriteHooks("BeforeConfigWrite", configHooks.BeforeConfigWrite, args)
	}
	if err := conf.SaveSSHConfig(); err != nil {
		return err
	}
	if configHooks != nil {
		invokeConfigWriteHooks("AfterConfigWrite", configHooks.AfterConfigWrite, args)
	}
	return nil
}

func invokeConfigWriteHooks(name string, list hooks.Hooks, args configWriteHookArgs) {
	drivers, err := list.InvokeAll(args)
	if err != nil {
		logger().Error("Configuration write hook failed", zap.String("hook", name), zap.Error(err))
		return
	}
	for _, err := range drivers.Close() {
		logger().Error("Cannot close configuration write hook", zap.String("hook", name), zap.Error(err))
	}
}
