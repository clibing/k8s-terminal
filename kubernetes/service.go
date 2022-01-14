package kubernetes

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

type ServiceApi interface {
	/**
	 * service 列表, 支持查询
	 */
	ServiceList(namespace, filter string, pageNumber, pageSize int) (value DeploymentResponse)
	/**
	 * 获取 service 详细信息
	 */
	ServiceDetail(namespace, name string) (service Service)
}

type ServiceModel struct {
	BaseModel
}

/**
* 获取deployment的副本集
 */
func (m *ServiceModel) replicaSetDetailGetServiceName(namespace, podName string)(name string){
	url := fmt.Sprintf("https://%s:%d/api/v1/replicaset/%s/%s?filterBy=&itemsPerPage=50&page=1&sortBy=d,creationTimestamp", m.req.Ip, m.req.Port, namespace, podName)

	data, err := commonRequest(url, false, nil, false, true, nil)
	if err != nil {
		panic(err)
	}

	var response ReplicaSetResponse
	json.Unmarshal([]byte(data), &response)
	name = response.ServiceList.Services[0].ObjectMeta.Name;
	return
}


func (m *ServiceModel) ServiceDetail(namespace, podName string) (service Service) {

	name := m.replicaSetDetailGetServiceName(namespace, podName)

	url := fmt.Sprintf("https://%s:%d/api/v1/service/%s/%s", m.req.Ip, m.req.Port, namespace, name)

	data, err := commonRequest(url, false, nil, false, true, nil)
	if err != nil {
		panic(err)
	}

	var response Service
	json.Unmarshal([]byte(data), &response)
	service = response
	return
}

/**
 * 展示 service 的详情
 */
func ShowServiceDetails(m *ServiceModel, deployment Deployment) {
	size := len(deployment.PodList.Pods)
	data := make([][]string, 1)

	if size == 0 {
		return
	} else if size >= 1 {
		for i := 0; i < size; i++ {
			sd := m.ServiceDetail(deployment.ObjectMeta.Namespace, deployment.NewReplicaSet.ObjectMeta.Name)
			endPointLen := len(sd.InternalEndpoint.Ports)
			for j := 0; j < endPointLen; j++ {
				column := make([]string, 6)
				if j == 0 {
					column[0] = strconv.Itoa(i*endPointLen + j + 1)
					column[1] = sd.ObjectMeta.Name
					column[2] = sd.ObjectMeta.CreationTimestamp.Format("2006-01-02 15:04:05")
				} else {
					column[0] = " "
					column[1] = " "
					column[2] = " "
				}
				ports := sd.InternalEndpoint.Ports[j]
				column[3] = fmt.Sprintf("%d(%s)", ports.Port, ports.Protocol)
				column[4] = strconv.Itoa(ports.NodePort)
				for _, v := range sd.EndpointList.Endpoints {
					for _, w := range v.Ports {
						if w.Port == ports.Port {
							column[5] = w.Name
						}
					}
				}
				data = append(data, column)
			}
		}
		showServiceTable(data, size)
	}
}

func showServiceTable(data [][]string, size int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NO", "Service_Name", "Create_Time", "End_Point", "Node_Port", "DESCRIPTION"})
	table.SetFooter([]string{"", "", "", "", "Total", strconv.Itoa(size)})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
}
