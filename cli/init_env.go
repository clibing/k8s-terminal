package cli

import (
	"github.com/k8s-terminal/config"
	"github.com/urfave/cli/v2"
)

func initEnvCommand() *cli.Command {
	return &cli.Command{
		Name:    initEnv.Command,
		Usage:   "环境初始化",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "ip",
				Usage:   "kubernetes dashboard ip",
			},
			&cli.IntFlag{
				Name:    "port",
				Usage:   "kubernetes dashboard port",
			},
			&cli.StringFlag{
				Name:    "token",
				Usage:   "kubernetes dashboard login token",
			},
			&cli.BoolFlag{
				Name:    "force",
				Usage:   "如果配置文件存在，会被覆盖",
			},
		},
		Action: func(c *cli.Context) error {
			len := len(c.FlagNames())
			if len < 3 {
				panic("安装参数失败")
			}
			ip := c.String("ip")
			port := c.Int("port")
			token := c.String("token")
			force := c.Bool("force")
			return initEnvFunc(ip, token, port, force)
		},
	}
}

func initEnvFunc(ip, token string, port int, force bool) error{
	config.Create(ip, token, port, nil, force)
	return nil
}