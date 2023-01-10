package cli

import (
	"github.com/k8s-terminal/kubernetes"
	"github.com/urfave/cli/v2"
)

var (
	textRestart = "restart"
	textDesired = "desired"
	textScale   = "scale"
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

   重启： k8s-terminal deployment --namespace <namespace> -n <pod name> --restart
   缩减副本数为2个： k8s-terminal deployment --namespace <namespace> -n <pod name> --scale --desired 2
`,
		Flags: []cli.Flag{
			commonNamespaceFlag,
			commonNameFlag,
			&cli.BoolFlag{
				Name:  textRestart,
				Usage: "重启pod(与scale互斥)",
			},
			&cli.BoolFlag{
				Name:  textScale,
				Usage: "是否缩减POD的副本(与restart互斥)",
			},
			&cli.IntFlag{
				Name:  textDesired,
				Value: 1,
				Usage: "Pod的副本数, 默认1",
			},
		},
		Action: func(c *cli.Context) error {
			restart := c.Bool(textRestart)
			scale := c.Bool(textScale)
			desired := c.Int(textDesired)
			ns := c.String(COMMON_NAMESPACE_NAME)
			if ns == "" {
				ns = c.String("ns")
			}
			if ns == "" {
				panic("请输入指定的namespace")
			}

			name := c.String(COMMON_NAME_NAME)
			if name == "" {
				name = c.String("n")
			}

			if restart && !scale {
				kubernetes.RestartPod(req, ns, name)
				return nil
			}
			if scale && !restart {
				kubernetes.ScalePod(req, ns, name, desired)
				return nil
			}

			kubernetes.ShowDeployment(ns, name, req)
			return nil
		},
	}
}
