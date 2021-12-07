package kubernetes

import (
	"github.com/k8s-terminal/config"
	"testing"
)

func TestShowDeployment(t *testing.T) {
	config.DefLoad()
	req := &Request{
		Ip:    config.GlobalCfg.Cluster.Ip,
		Port:  config.GlobalCfg.Cluster.Port,
		Token: config.GlobalCfg.Cluster.Auth.Token.Value,
	}
	req.CsrfToken()
	req.JweToken()
	ShowDeployment("demo", "demo", req)
}
