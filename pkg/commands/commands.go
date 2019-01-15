package commands // import "moul.io/assh/pkg/commands"

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/assh/pkg/config"
	loggerpkg "moul.io/assh/pkg/logger"
	"moul.io/assh/pkg/version"
)

var commands = []*cobra.Command{
	pingCommand,
	proxyCommand,
	infoCommand,
	configCommand,
	socketsCommand,
	wrapperCommand,
	genCompletionsCommand,
}

var RootCmd = &cobra.Command{
	Use:              "assh",
	Short:            "assh - advanced ssh config",
	Version:          version.VERSION + " (" + version.GITCOMMIT + ")",
	TraverseChildren: true,
}

func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	abspath, err := filepath.Abs(ex)
	if err != nil {
		log.Fatal(err)
	}
	config.SetASSHBinaryPath(abspath)

	RootCmd.Flags().BoolP("help", "h", false, "print usage")
	RootCmd.Flags().StringP("config", "c", "~/.ssh/assh.yml", "Location of config file")
	RootCmd.Flags().BoolP("debug", "D", false, "Enable debug mode")
	RootCmd.Flags().BoolP("verbose", "V", false, "Enable verbose mode")

	viper.BindEnv("debug", "ASSH_DEBUG")
	viper.BindEnv("config", "ASSH_CONFIG")
	viper.BindPFlags(RootCmd.Flags())

	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if viper.GetBool("debug") {
			if err := os.Setenv("ASSH_DEBUG", "1"); err != nil {
				return err
			}
		}
		if err := initLogging(viper.GetBool("debug"), viper.GetBool("verbose")); err != nil {
			return err
		}
		return nil
	}

	RootCmd.AddCommand(commands...)
}

func initLogging(debug bool, verbose bool) error {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(loggerpkg.MustLogLevel(debug, verbose))
	if !debug {
		config.DisableStacktrace = true
		config.DisableCaller = true
		config.EncoderConfig.TimeKey = ""
		config.EncoderConfig.NameKey = ""
	}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	l, err := config.Build()
	if err != nil {
		return errors.Wrap(err, "failed to initialize logger")
	}
	zap.ReplaceGlobals(l)
	return nil
}
