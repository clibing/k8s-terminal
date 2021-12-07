package config

import (
	"fmt"
	"github.com/k8s-terminal/util"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const execHome string = ".k8s-terminal"
const fileName string = "config.yaml"

var (
	GlobalCfg Config
)

func CheckInstalled()(exist bool){
	homedir, err := homedir.Dir()
	if err != nil {
		log.Panicln("home dir err")
		return
	}
	path := filepath.Join(homedir, execHome, fileName)
	exist = util.FileExist(path)
	return
}

/**
 * 加载系统默认的配置文件
 */
func DefLoad() {
	homedir, err := homedir.Dir()
	if err != nil {
		log.Panicln("home dir err")
		return
	}
	join := filepath.Join(homedir, execHome)
	if !util.FileExist(join) {
		err = os.Mkdir(join, 0755)
		if err != nil {
			log.Panicln("create workspace err: ", execHome)
		}
	}
	path := filepath.Join(homedir, execHome, fileName)
	Load(path)
}


/**
 * 加载指定的配置文件
 */
func Load(path string) {
	exist := util.FileExist(path)
	if !exist {
		return
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("read yml file err %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &GlobalCfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func Create(ip, token string, port int, namespace []string, force bool) {
	homedir, err := homedir.Dir()
	if err != nil {
		log.Panicln("home dir err")
		return
	}

	join := filepath.Join(homedir, execHome)
	if !util.FileExist(join) {
		err = os.Mkdir(join, 0755)
		if err != nil {
			log.Panicln("create workspace err: ", execHome)
		}
	}

	cFile := filepath.Join(homedir, execHome, fileName)
	if util.FileExist(cFile) {
		if !force {
			panic("文件已经存在，请删除: " + join)
		} else {
			input, err := ioutil.ReadFile(cFile)
			if err != nil {
				fmt.Println(err)
				return
			}

			bakFile := filepath.Join(homedir, execHome, fmt.Sprintf("%s.%d", fileName, time.Now().Unix()))
			err = ioutil.WriteFile(bakFile, input, 0644)
			if err != nil {
				fmt.Println("备份文件创建失败", bakFile)
				fmt.Println(err)
				return
			}
			fmt.Println("配置文件存在， 备份文件为", bakFile)
		}
	}

	v := Config{
		Cluster: Cluster{
			Ip:   ip,
			Port: port,
			Auth: Auth{
				Token: Token{
					Value: token,
				},
				KubeConfig: "",
			},
		},
		Namespace: Namespace{
			Static: false,
		},
		Log: Log{
			PageSize: 100,
		},
	}
	if len(namespace) > 0 {
		v.Namespace.Static = true
		v.Namespace.Names = namespace
	}

	marshal, err := yaml.Marshal(&v)
	ioutil.WriteFile(cFile, marshal, 0644)
}

func PersistenceConfigFile(cfg *Config, force bool) {
	homedir, err := homedir.Dir()
	if err != nil {
		log.Panicln("home dir err")
		return
	}

	cFile := filepath.Join(homedir, execHome, fileName)
	if util.FileExist(cFile) {
		if !force {
			return
		} else {
			input, err := ioutil.ReadFile(cFile)
			if err != nil {
				fmt.Println(err)
				return
			}

			bakFile := filepath.Join(homedir, execHome, fmt.Sprintf("%s.%d", fileName, time.Now().Unix()))
			err = ioutil.WriteFile(bakFile, input, 0644)
			if err != nil {
				fmt.Println("备份文件创建失败", bakFile)
				fmt.Println(err)
				return
			}
			fmt.Println("配置文件存在， 备份文件为", bakFile)
		}
	}
	marshal, err := yaml.Marshal(cfg)
	ioutil.WriteFile(cFile, marshal, 0644)
}
