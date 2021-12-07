package cli

import (
	"github.com/k8s-terminal/kubernetes"
	"github.com/urfave/cli/v2"
)

func deploymentCommand(req *kubernetes.Request) *cli.Command {
	return &cli.Command{
		Name:    deployment.Command,
		Aliases: []string{deployment.Abbreviations},
		Usage:   "Deployment相关操作, 选择对应的deployment部署信息和service对应的端口",
		UsageText: `
   使用方法：
      空间名字: <namespace> 空间
      部署名字: <deployment name> deployment的名字

   简写方式的命令：k8s-terminal deployment --ns <namespace> -n <deployment name>
   完整方式的命令：k8s-terminal deployment --deployment-namespace <namespace> --deployment-name <deployment name>
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "deployment-namespace",
				Aliases: []string{"ns"},
				Usage:   "deployment当前所在的namespace",
			},
			&cli.StringFlag{
				Name:    "deployment-name",
				Aliases: []string{"n"},
				Usage:   "根据deployment的名字进行所说",
			},
		},
		Action: func(c *cli.Context) error {

			ns := c.String("deployment-namespace")
			if ns == "" {
				ns = c.String("ns")
			}
			if ns == "" {
				panic("请输入指定的namespace")
			}

			name := c.String("deployment-name")
			if name == "" {
				name = c.String("n")
			}
			kubernetes.ShowDeployment(ns, name, req)
			return nil
		},
	}
}
