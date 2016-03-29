package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	. "github.com/moul/advanced-ssh-config/pkg/logger"
)

func cmdWrapper(c *cli.Context) {
	/*conf, err := config.Open()
	if err != nil {
		Logger.Fatalf("Cannot open configuration file: %v", err)
	}*/

	Logger.Debugf("Wrapper called with %v", c.Args())
	fmt.Println(c.Args())

	//conf.WriteSshConfigTo(os.Stdout)
}
