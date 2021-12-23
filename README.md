kubernetes dashboard "替代品"

在terminal发起http请求查看kubernetes相关的信息

> 注意 
>
> 1. 本程序以查看当前namespace、deployment、pod等信息为主， 面向业务开发人员，主要原因是权限受限。如果是运维人员绕行即可
> 2. 当前只是使用与token登录方法，暂不支持kubeConfig方式 

#### 1. 安装

下载二进制安装包

以MacOS为例

```shell script
chmod a+x k8s-terminal_darwin_amd64
./k8s-terminal_darwin_amd64 install  # 默认安装到 /usr/local/bin/ 目录下
```

或者 下载源码直接make编译

```shell script
cd k8s-terminal
make
```

#### 2. 使用

##### 2.1 初始化

```shell script
k8s-terminal init \ 
    --ip <k8s dashboard web ip>  \
    --port <k8s dashboard web port>  \
    --token <k8s dashboard login with token>
    --force # 如果别你存在配置文件会进行备份后覆盖修改
````

##### 2.2 查看当前namespace

查看全部namesapce

````shell script
k8s-terminal namespace 
k8s-terminal n
k8s-terminal n -h 
````

##### 2.3 查看deployment 

功能描述: 查看namespace下的某个deployment的具体信息, 使用场景例如查看某个deployment的端口，部署信息等

```shell script
k8s-terminal deployment --ns <namespace> -n <deployment name>
k8s-terminal deployment --deployment-namespace <namespace> --deployment-name <deployment name>
```

##### 2.4 查看实时日志信息

````shell script
k8s-terminal pod --ns <namespace> -n <pod name>
k8s-terminal pod --pod-namespace <namespace> --pod-name <pod name>
k8s-terminal pod --ns <namespace> -n <pod name> -e
````

##### 2.5 更多

查看 secret 中的内容

更多的使用细节查看帮助文档

````shell script
k8s-terminal -h
k8s-terminal <command> -h
````

#### 附录

1. bubbletea


