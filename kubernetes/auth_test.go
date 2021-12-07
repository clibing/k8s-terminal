package kubernetes

import (
	"fmt"
	"github.com/k8s-terminal/config"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	config.DefLoad()
	req := &Request{
		Ip:    config.GlobalCfg.Cluster.Ip,
		Port:  config.GlobalCfg.Cluster.Port,
		Token: config.GlobalCfg.Cluster.Auth.Token.Value,
	}

	req.CsrfToken()
	fmt.Println(CsrfTokenModel.Value)
	req.JweToken()
	fmt.Println(JweTokenModel.JweToken)

	req.RefreshCsrtToken()
	fmt.Println(CsrfTokenModel.Value)

	req.RefreshJweToken()
	fmt.Println(JweTokenModel.JweToken)
}
