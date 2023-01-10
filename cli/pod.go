package cli

import (
	"github.com/k8s-terminal/config"
	"github.com/k8s-terminal/kubernetes"
	"github.com/urfave/cli/v2"
)

var (
	textPodEnableLog = "enableLog"
	textTail         = "tail"
	textDownload     = "download"
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
			commonNamespaceFlag,
			commonNameFlag,
			commonPathFlag,
			&cli.BoolFlag{
				Name:    textPodEnableLog,
				Aliases: []string{"e"},
				Usage:   "查看Pod日志",
			},
			&cli.IntFlag{
				Name:    textTail,
				Aliases: []string{"l"},
				Usage:   "查看日志的位置，支持正负数",
			},
			&cli.BoolFlag{
				Name:    "download",
				Aliases: []string{"d"},
				Usage:   "下载文件到当前工作目录下~/Download/，日志的文件名字使用ContainerName",
			},
		},

		Action: func(c *cli.Context) error {
			namespace := ""
			name := ""
			enable := false
			tail := 0
			download := false
			downloadPath := ""

			for _, v := range c.FlagNames() {
				if v == COMMON_NAMESPACE_NAME {
					namespace = c.String(COMMON_NAMESPACE_NAME)
					continue
				}
				if v == "ns" {
					namespace = c.String("ns")
					continue
				}

				if v == COMMON_NAME_NAME {
					name = c.String(COMMON_NAME_NAME)
					continue
				}
				if v == "n" {
					name = c.String("n")
					continue
				}

				if v == textPodEnableLog {
					enable = c.Bool(textPodEnableLog)
					continue
				}
				if v == "e" {
					enable = c.Bool("e")
					continue
				}

				if v == textTail{
					tail = c.Int(textTail)
					continue
				}
				if v == "l" {
					tail = c.Int("l")
					continue
				}

				if v == textDownload{
					download = c.Bool(textDownload)
					continue
				}

				if v == "d" {
					download = c.Bool("d")
					continue
				}

				if v == COMMON_PATH_NAME {
					downloadPath = c.String(COMMON_PATH_NAME)
					continue
				}

				if v == "p" {
					downloadPath = c.String("p")
					continue
				}
			}
			kubernetes.ShowPod(req, namespace, name, enable, config.GlobalCfg.Log.PageSize, tail, download, downloadPath)
			return nil
		},
	}
}
