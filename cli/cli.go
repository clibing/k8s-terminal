package cli

import (
	"fmt"
	"github.com/k8s-terminal/config"
	"github.com/k8s-terminal/kubernetes"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	initEnv   = Commands{Command: "init"}
	install   = Commands{Command: "install", Abbreviations: "i",}
	uninstall = Commands{Command: "uninstall", Abbreviations: "u",}

	namespace  = Commands{Command: "namespace", Abbreviations: "n",}
	deployment = Commands{Command: "deployment", Abbreviations: "d",}
	service    = Commands{Command: "service", Abbreviations: "s",}
	pod        = Commands{Command: "pod", Abbreviations: "p",}
	secret = Commands{Command: "secret", Abbreviations: "",}

	// todo
	configMap = Commands{Command: "configMap", Abbreviations: "",}

	watchLog     = Commands{Command: "log", Abbreviations: "",}
	tail         = Commands{Command: "tail", Abbreviations: "",}
	filter       = Commands{Command: "filter", Abbreviations: "f",}
	req          *kubernetes.Request
	refreshToken *time.Ticker = time.NewTicker(time.Second * 15)
	start = make(chan bool, 1)
)

func init() {

	config.DefLoad()
	req = &kubernetes.Request{
		Ip:    config.GlobalCfg.Cluster.Ip,
		Port:  config.GlobalCfg.Cluster.Port,
		Token: config.GlobalCfg.Cluster.Auth.Token.Value,
	}

	// 异步刷新token
	go RefreshTokenAsync(req)
}

func RefreshTokenAsync(req *kubernetes.Request) {
	defer refreshToken.Stop()
	_ = <-start
	for {
		_ = <-refreshToken.C
		//fmt.Println("异步刷新token.....")
		req.RefreshCsrtToken()
		req.RefreshJweToken()
	}
}

func CliApp(version, buildDate, commitID string) {
	app := &cli.App{
		Name:  "k8s-terminal",
		Usage: "k8s集群，主要替代kubernetes dashboard token登录，对deployment, pod, service, namespace, configMap的查看",
		Authors:              []*cli.Author{{Name: "clibing", Email: "wmsjhappy@gmail.com"}},
		EnableBashCompletion: true,
		Copyright:            "Copyright (c) " + time.Now().Format("2006") + " clibing, All rights reserved.",
		Version:              fmt.Sprintf("%s / %s / %s", version, buildDate, commitID),
		Before: func(context *cli.Context) error {
			if config.CheckInstalled() {
				req.CsrfToken()
				req.JweToken()
				start <- true
			}
			return nil
		},
		CommandNotFound: func(context *cli.Context, s string) {
			fmt.Println("暂不支持该命令", s)
		},
		Commands: []*cli.Command{
			namespaceCommand(req),
			deploymentCommand(req),
			podCommand(req),
			installCommand(),
			initEnvCommand(),
			secretCommand(req),
		},
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
