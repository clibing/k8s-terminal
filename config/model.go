package config

var ConfigTemplate = `
# namespace 相关配置
auth-config:
  token: "{{ input token value }}"
namespace-config:
  # 是否使用静态模式，默认关闭，如果开启，会使用 names下的配置
  static: false
  names:
    - "{{ input namespace }}"
`

/**
 * k8s-terminal 全局配置
 */
type Config struct {
	Cluster   Cluster   `yaml:"cluster"`
	Namespace Namespace `yaml:"namespace"`
	Log       Log       `yaml:"log"`
}

/**
 * 集群dashboard的管理平台入口
 * 注意 请求的https的证书时私有证书，导致在根验证证书异常，系统默认http请求时tls不进行验证
 */
type Cluster struct {
	Ip   string `yaml:"ip"`
	Port int    `yaml:"port"`
	Auth Auth   `yaml:"auth"`
}

/**
 * 认证方式
 */
type Auth struct {
	Token      Token  `yaml:"token"`
	KubeConfig string `yaml:"kubeConfig"`
}

/**
 * 采用Token认证
 */
type Token struct {
	Value   string `yaml:"value"`
	Refresh int    `yaml:"refresh"` // 每个多久刷新一次token
}

/**
 * k8s 的空间相关配置
 */
type Namespace struct {
	Static bool     `yaml:"static"` // 是否使用静态配置， 默认: false; 如果: true, 会向k8s集群动态请求获取所有的空间并与静态的合并
	Names  []string `yaml:"names"`  // namespace的名字
}

type Log struct {
	PageSize int `yaml:"pageSize"`
}
