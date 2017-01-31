package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/moul/advanced-ssh-config/pkg/config"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bind-address",
			Value: ":8080",
		},
	}
	app.Action = server
	app.Run(os.Args)
}

func server(c *cli.Context) error {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
	router.POST("/assh-to-ssh", func(c *gin.Context) {
		var (
			err    error
			cfg    = config.New()
			buffer bytes.Buffer
			json   struct {
				AsshConfig string `form:"assh_config" json:"assh_config"`
			}
		)

		if err = c.BindJSON(&json); err != nil {
			goto serverEnd
		}

		if json.AsshConfig == "" {
			err = fmt.Errorf("invalid input")
			goto serverEnd
		}

		if err = cfg.LoadConfig(strings.NewReader(json.AsshConfig)); err != nil {
			goto serverEnd
		}

		if err = cfg.WriteSSHConfigTo(&buffer); err != nil {
			goto serverEnd
		}

	serverEnd:
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				//"assh_config": json.AsshConfig,
				"ssh_config": buffer.String(),
			})
		}
	})
	router.Run(c.String("bind-address"))
	return nil
}
