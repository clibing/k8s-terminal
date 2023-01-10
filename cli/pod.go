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

   重启pod： k8s-terminal pod --namespace <namespace> -n <pod name> --restart
   缩减pod副本数为2个： k8s-terminal pod --namespace <namespace> -n <pod name> --scale --desired 2
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
				Usage:   "根据Pod的名字进行过滤",
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
			&cli.BoolFlag{
				Name:    "download-log",
				Aliases: []string{"d"},
				Usage:   "下载文件到当前工作目录下~/Download/，日志的文件名字使用ContainerName",
			},
			&cli.StringFlag{
				Name:    "download-path",
				Aliases: []string{"dp"},
				Usage:   "保存日志的目录",
			},
			&cli.BoolFlag{
				Name:  "restart",
				Usage: "重启pod(与scale互斥)",
			},
			&cli.BoolFlag{
				Name:  "scale",
				Usage: "是否缩减POD的副本(与restart互斥)",
			},
			&cli.IntFlag{
				Name:  "desired",
				Usage: "Pod的副本数",
			},
		},

		Action: func(c *cli.Context) error {
			ns := ""
			name := ""
			enable := false
			tail := 0
			download := false
			downloadPath := ""
			restartPod := false
			scale := false
			desired := 0
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

				if v == "download-log" {
					download = c.Bool("download-log")
					continue
				}

				if v == "d" {
					download = c.Bool("d")
					continue
				}

				if v == "download-path" {
					downloadPath = c.String("download-path")
					continue
				}

				if v == "dp" {
					downloadPath = c.String("dp")
					continue
				}

				if v == "restart" {
					restartPod = c.Bool("restart")
					continue
				}

				if v == "scale" {
					scale = c.Bool("scale")
					continue
				}

				if v == "desired" {
					desired = c.Int("desired")
					continue
				}

			}
			if restartPod && !scale {
				kubernetes.RestartPod(req, ns, name)
				return nil
			}
			if scale && !restartPod {
				kubernetes.ScalePod(req, ns, name, desired)
				return nil
			}
			kubernetes.ShowPod(req, ns, name, enable, config.GlobalCfg.Log.PageSize, tail, download, downloadPath)
			return nil
		},
	}
}
