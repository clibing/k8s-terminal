package kubernetes

import (
	"time"
)

//======================================================================================================================
type Request struct {
	Ip    string
	Port  int
	Token string

}

/**
 * 跨域响应体
 */
type CsrfToken struct {
	Value string `json:"token"`
}

type LoginRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Token      string `json:"token"`
	KubeConfig string `json:"kubeConfig"`
}

/**
 * 认证返回结果
 */
type JweToken struct {
	JweToken string  `json:"jweToken"`
	Errors   []error `json:"errors"`
}

type RJweToken struct {
	JweToken string  `json:"jweToken"`
}
/**
 * jwe token 下的各个属性
 */
type JweTokenAttribute struct {
	Protected    string `json:"protected"`
	Aad          string `json:"aad"`
	EncryptedKey string `json:"encrypted_key"`
	Iv           string `json:"iv"`
	CipherText   string `json:"ciphertext"`
	Tag          string `json:"tag"`
}

//======================================================================================================================
/**
 * deployment 响应体
 */
type DeploymentResponse struct {
	ListMeta   ListMeta     `json:"listMeta"`    // 一共总deployment记录数
	Deployment []Deployment `json:"deployments"` // 对应的deployment list
	Status     Status       `json:"status"`      //

	CumulativeMetrics []interface{} `json:"cumulativeMetrics"` // 不需要考虑这个字段，保留
	Errors            []interface{} `json:"errors"`
}

/**
 * 列表总记录
 */
type ListMeta struct {
	TotalItems int `json:"totalItems"`
}

/**
 * 系统负载情况
 */
type Status struct {
	Running   int `json:"running"` // 当前正在运行的，大部分都是这个状态
	Pending   int `json:"pending"`
	Failed    int `json:"failed"`
	Succeeded int `json:"succeeded"`
}

/**
 * meta信息
 */
type ObjectMeta struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Labels            Labels            `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
}

type Labels struct {
	App             string `json:"app"`               // deployment replicaSet 使用
	Name            string `json:"name"`              // namespace 使用
	PodTemplateHash string `json:"pod-template-hash"` // pod replicaSet 使用
}

/**
 * 请求的类型
 * 基本用不到
 */
type TypeMeta struct {
	Kind string `json:"kind"` // service, secret, replicaSets, pod, namespace, deployment, configMap
}

/**
 * pod 运行状态 也没啥大用 ，当前开发环境只有1个， 需要展示Running(Pending)个数即可
 */
type PodStatus struct {
	// deployment 列表使用
	Current   int           `json:"current"`
	Desired   int           `json:"desired"`
	Running   int           `json:"running"`
	Pending   int           `json:"pending"`
	Failed    int           `json:"failed"`
	Succeeded int           `json:"succeeded"`
	Warnings  []interface{} `json:"warnings"`

	// deployment 详情使用
	Status          string            `json:"status"`
	PodPhase        string            `json:"podPhase"`
	ContainerStates []ContainerStates `json:"containerStates"`
}

type Running struct {
	StartedAt time.Time `json:"startedAt"`
}
type ContainerStates struct {
	Running Running `json:"running"`
}

type Pod struct {
	ObjectMeta   ObjectMeta    `json:"objectMeta"`
	TypeMeta     TypeMeta      `json:"typeMeta"`
	PodStatus    PodStatus     `json:"podStatus"`
	RestartCount int           `json:"restartCount"`
	Metrics      interface{}   `json:"metrics"`
	Warnings     []interface{} `json:"warnings"`
	NodeName     string        `json:"nodeName"`
}

/**
 * deployment 结构
 */
type Deployment struct {
	ObjectMeta      ObjectMeta `json:"objectMeta"`
	TypeMeta        TypeMeta   `json:"typeMeta"`
	PodStatus       PodStatus  `json:"pods"`
	ContainerImages []string   `json:"containerImages"`

	InitContainerImages interface{} `json:"initContainerImages"` // 保留

	PodList                     PodList                     `json:"podList"`
	Labels                      Labels                      `json:"selector"`
	StatusInfo                  StatusInfo                  `json:"statusInfo"`
	Strategy                    string                      `json:"strategy"`
	MinReadySeconds             int                         `json:"minReadySeconds"`
	RollingUpdateStrategy       RollingUpdateStrategy       `json:"rollingUpdateStrategy"`
	OldReplicaSetList           OldReplicaSetList           `json:"oldReplicaSetList"`
	NewReplicaSet               NewReplicaSet               `json:"newReplicaSet"`
	RevisionHistoryLimit        int                         `json:"revisionHistoryLimit"`
	EventList                   EventList                   `json:"eventList"`
	HorizontalPodAutoscalerList HorizontalPodAutoscalerList `json:"horizontalPodAutoscalerList"`
	Errors                      []interface{}               `json:"errors"`
}

//======================================================================================================================
type PodList struct {
	ListMeta          ListMeta      `json:"listMeta"`
	CumulativeMetrics []interface{} `json:"cumulativeMetrics"`
	Status            Status        `json:"status"`
	Pods              []Pod        `json:"pods"`
	Errors            []interface{} `json:"errors"`
}

type StatusInfo struct {
	Replicas    int `json:"replicas"`
	Updated     int `json:"updated"`
	Available   int `json:"available"`
	Unavailable int `json:"unavailable"`
}
type RollingUpdateStrategy struct {
	MaxSurge       int `json:"maxSurge"`
	MaxUnavailable int `json:"maxUnavailable"`
}
type OldReplicaSetList struct {
	ListMeta          ListMeta      `json:"listMeta"`
	CumulativeMetrics []interface{} `json:"cumulativeMetrics"`
	Status            Status        `json:"status"`
	ReplicaSets       []interface{} `json:"replicaSets"`
	Errors            []interface{} `json:"errors"`
}

type NewReplicaSet struct {
	ObjectMeta          ObjectMeta  `json:"objectMeta"`
	TypeMeta            TypeMeta    `json:"typeMeta"`
	Pods                Pod        `json:"pods"`
	ContainerImages     []string    `json:"containerImages"`
	InitContainerImages interface{} `json:"initContainerImages"`
}

type Events struct {
	ObjectMeta      ObjectMeta `json:"objectMeta"`
	TypeMeta        TypeMeta   `json:"typeMeta"`
	Message         string     `json:"message"`
	SourceComponent string     `json:"sourceComponent"`
	SourceHost      string     `json:"sourceHost"`
	Object          string     `json:"object"`
	Count           int        `json:"count"`
	FirstSeen       time.Time  `json:"firstSeen"`
	LastSeen        time.Time  `json:"lastSeen"`
	Reason          string     `json:"reason"`
	Type            string     `json:"type"`
}
type EventList struct {
	ListMeta ListMeta `json:"listMeta"`
	Events   []Events `json:"events"`
}
type HorizontalPodAutoscalerList struct {
	ListMeta                 ListMeta      `json:"listMeta"`
	Horizontalpodautoscalers []interface{} `json:"horizontalpodautoscalers"`
	Errors                   []interface{} `json:"errors"`
}
//=service=====================================================================================================================

type Service struct {
	ObjectMeta ObjectMeta `json:"objectMeta"`
	TypeMeta TypeMeta `json:"typeMeta"`
	InternalEndpoint InternalEndpoint `json:"internalEndpoint"`
	ExternalEndpoints interface{} `json:"externalEndpoints"`
	EndpointList EndpointList `json:"endpointList"`
	Selector Selector `json:"selector"`
	Type string `json:"type"`
	ClusterIP string `json:"clusterIP"`
	EventList EventList `json:"eventList"`
	PodList PodList `json:"podList"`
	SessionAffinity string `json:"sessionAffinity"`
	Errors []interface{} `json:"errors"`
}

type Ports struct {
	Name string `json:"name"`
	Port int `json:"port"`
	Protocol string `json:"protocol"`
	NodePort int `json:"nodePort"`
}
type InternalEndpoint struct {
	Host string `json:"host"`
	Ports []Ports `json:"ports"`
}

type Endpoints struct {
	ObjectMeta ObjectMeta `json:"objectMeta"`
	TypeMeta TypeMeta `json:"typeMeta"`
	Host string `json:"host"`
	NodeName string `json:"nodeName"`
	Ready bool `json:"ready"`
	Ports []Ports `json:"ports"`
}
type EndpointList struct {
	ListMeta ListMeta `json:"listMeta"`
	Endpoints []Endpoints `json:"endpoints"`
}
type Selector struct {
	App string `json:"app"`
}
//=service=====================================================================================================================

type PodListResponse struct {
	ListMeta ListMeta `json:"listMeta"`
	CumulativeMetrics []interface{} `json:"cumulativeMetrics"`
	Status Status `json:"status"`
	Pods []Pod `json:"pods"`
	Errors []interface{} `json:"errors"`
}
//=pod=====================================================================================================================
type LogResponse struct {
	Info Info           `json:"info"`
	Selection Selection `json:"selection"`
	Logs []Log          `json:"logs"`
}
type Info struct {
	PodName string `json:"podName"`
	ContainerName string `json:"containerName"`
	InitContainerName string `json:"initContainerName"`
	FromDate time.Time `json:"fromDate"`
	ToDate time.Time `json:"toDate"`
	Truncated bool `json:"truncated"`
}
type ReferencePoint struct {
	Timestamp time.Time `json:"timestamp"`
	LineNum int `json:"lineNum"`
}
type Selection struct {
	ReferencePoint ReferencePoint `json:"referencePoint"`
	OffsetFrom int `json:"offsetFrom"`
	OffsetTo int `json:"offsetTo"`
	LogFilePosition string `json:"logFilePosition"`
}
type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Content string `json:"content"`
}
//=pod=====================================================================================================================
type NamespaceResponse struct {
	ListMeta ListMeta `json:"listMeta"`
	Namespaces []Namespace `json:"namespaces"`
	Errors []interface{} `json:"errors"`
}

type Namespace struct {
	ObjectMeta ObjectMeta `json:"objectMeta,omitempty"`
	TypeMeta TypeMeta `json:"typeMeta"`
	Phase string `json:"phase"`
}

// = secret =================================
type SecretListResponse struct {
	ListMeta ListMeta `json:"listMeta"`
	Secret []Secret `json:"secrets"`
	Errors []interface{} `json:"errors"`
}

type Secret struct {
	ObjectMeta ObjectMeta `json:"objectMeta"`
	TypeMeta TypeMeta `json:"typeMeta"`
	Type string `json:"type"`
}

type SecretResponse struct {
	ObjectMeta ObjectMeta `json:"objectMeta"`
	TypeMeta TypeMeta `json:"typeMeta"`
	Data map[string]string `json:"data"`
	Type string `json:"type"`
}

