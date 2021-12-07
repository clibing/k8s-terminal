package cli

import (
	"github.com/k8s-terminal/kubernetes"
	"github.com/urfave/cli/v2"
)

func namespaceCommand(req *kubernetes.Request) *cli.Command {
	return &cli.Command{
		Name:    namespace.Command,
		Aliases: []string{namespace.Abbreviations},
		Usage:   "获取当前kubernetes集群的namespace",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "namespace-filter",
				Aliases: []string{"nf", "f"},
				Usage:   "是否为过滤模式",
			},
			&cli.BoolFlag{
				Name:    "save",
				Aliases: []string{"s"},
				Usage:   "将选择的namespace保存到配置文件中",
			},
		},

		Action: func(c *cli.Context) error {
			var key string
			save := false
			for _, v := range c.FlagNames() {
				if v == "namespace-filter" || v == "nf" || v == "f" {
					key = c.String(v)
					break
				}
				if v == "save" {
					save = c.Bool("save")
				}
			}
			kubernetes.ShowNamespace(key, save, req)
			return nil
		},
	}
}
