package cli

import (
	"github.com/k8s-terminal/kubernetes"
	"github.com/urfave/cli/v2"
)

func secretCommand(req *kubernetes.Request) *cli.Command {
	return &cli.Command{
		Name:    secret.Command,
		Aliases: []string{secret.Abbreviations},
		Usage:   "查看对应的secret详细配置内容",
		UsageText: `
   使用方法：
      空间名字: <namespace> 空间
      保密字典: <secret name> secret 的名字

   简写方式的命令：k8s-terminal secret --ns <namespace> -n <secret name>
   完整方式的命令：k8s-terminal secret --secret-namespace <namespace> --secret-name <pod name>
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "secret-namespace",
				Aliases: []string{"s", "ns"},
				Usage:   "secret 当前所在的namespace",
			},
			&cli.StringFlag{
				Name:    "secret-name",
				Aliases: []string{"n"},
				Usage:   "根据secret的名字进行所说",
			},
			&cli.StringFlag{
				Name:    "secret-edit",
				Aliases: []string{"e"},
				Usage:   "编辑secret内容，并同步到k8s",
			},
		},
		Action: func(c *cli.Context) error {
			ns := ""
			name := ""
			edit := false
			for _, v := range c.FlagNames() {
				if v == "secret-namespace" {
					ns = c.String("secret-namespace")
					continue
				}
				if v == "s" {
					ns = c.String("s")
					continue
				}

				if v == "secret-name" {
					name = c.String("secret-name")
					continue
				}
				if v == "n" {
					name = c.String("n")
					continue
				}

				if v == "secret-edit" {
					edit = c.Bool("secret-edit")
					continue
				}
				if v == "e" {
					edit = c.Bool("e")
					continue
				}
			}
			kubernetes.ShowSecret(ns, name, req, edit)
			return nil
		},
	}
}
