package cli

import (
	"github.com/k8s-terminal/config"
	"github.com/k8s-terminal/kubernetes"
	"github.com/urfave/cli/v2"
)

func podCommand(req *kubernetes.Request) *cli.Command {
	return &cli.Command{
		Name:    pod.Command,
		Aliases: []string{pod.Abbreviations},
		Usage:   "POD相关操作, 选择对应的Pod,查看POD的配置，支持实时查看Log",
		UsageText: `
   使用方法：
      空间名字: <namespace> 空间
      部署名字: <pod name> pod 的名字

   简写方式的命令：k8s-terminal pod --ns <namespace> -n <pod name>
   完整方式的命令：k8s-terminal pod --pod-namespace <namespace> --pod-name <pod name>

   简写方式的命令开启日志自动模式：k8s-terminal pod --ns <namespace> -n <pod name> -e
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "pod-namespace",
				Aliases: []string{"ns"},
				Usage:   "Pod当前所在的namespace",
			},
			&cli.StringFlag{
				Name:    "pod-name",
				Aliases: []string{"n"},
				Usage:   "根据Pod的名字进行所说",
			},
			&cli.BoolFlag{
				Name:    "enable-log",
				Aliases: []string{"e"},
				Usage:   "查看Pod日志",
			},
			&cli.IntFlag{
				Name:    "tail-line",
				Aliases: []string{"l"},
				Usage:   "查看日志的位置，支持正负数",
			},
		},
		Action: func(c *cli.Context) error {
			ns := ""
			name := ""
			enable := false
			tail := 0
			for _, v := range c.FlagNames() {
				if v == "pod-namespace" {
					ns = c.String("pod-namespace")
					continue
				}
				if v == "ns" {
					ns = c.String("ns")
					continue
				}

				if v == "pod-name" {
					name = c.String("pod-name")
					continue
				}
				if v == "n" {
					name = c.String("n")
					continue
				}

				if v == "enable-log" {
					enable = c.Bool("enable-log")
					continue
				}
				if v == "e" {
					enable = c.Bool("e")
					continue
				}

				if v == "tail-line" {
					tail = c.Int("tail-line")
					continue
				}
				if v == "l" {
					tail = c.Int("l")
					continue
				}
			}
			kubernetes.ShowPod(req, ns, name, enable, config.GlobalCfg.Log.PageSize, tail)
			return nil
		},
	}
}
