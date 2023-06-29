package kubernetes

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/k8s-terminal/util"
	"github.com/olekukonko/tablewriter"
)

type DeploymentApi interface {
	/**
	 * deployment 列表, 支持查询
	 */
	DeploymentList(namespace, filter string, pageNumber, pageSize int) (value DeploymentResponse)
	/**
	 * 获取deployment详细信息
	 */
	DeploymentDetail(namespace, name string) (deployment Deployment)
}

type DeploymentModel struct {
	BaseModel
	deployment map[string]Deployment // deployment 详情 这个玩意没有顺序
}

func (m *DeploymentModel) Init() tea.Cmd {
	data := m.DeploymentList(m.namespace, m.filter, m.pageNumber, m.pageSize)
	m.total = data.Status.Running + data.Status.Failed + data.Status.Pending + data.Status.Succeeded
	m.name = make([]string, 0)
	m.deployment = make(map[string]Deployment)
	m.pageNumber = defaultPageNumber
	m.pageSize = defaultPageSize
	for _, dep := range data.Deployment {
		name := dep.ObjectMeta.Name
		m.deployment[name] = dep
		m.name = append(m.name, name)
	}
	return nil
}

/**
 * 加载数据 支持向上加载，向下加载
 */
func fetchData(m *DeploymentModel, down bool) (err error) {
	pageNumber := m.pageNumber
	if pageNumber*m.pageSize > m.total {
		err = fmt.Errorf("已经拉取最大数据了，不需要再次拉取")
		return
	}

	if down {
		m.pageNumber = pageNumber + 1
	} else {
		m.pageNumber = pageNumber - 1
	}
	data := m.DeploymentList(m.namespace, m.filter, m.pageNumber, m.pageSize)
	for _, dep := range data.Deployment {
		name := dep.ObjectMeta.Name
		m.deployment[name] = dep
		m.name = append(m.name, name)
	}
	return nil
}

func (m *DeploymentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "enter":
			DeploymentValue = m.selected
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.name)-1 {
				m.cursor++
			}

		case " ":
			var value = m.name[m.cursor]
			m.selected = value
		}
	}
	return m, nil
}

func (m *DeploymentModel) View() string {
	s := "deployment name list(按空格选择deployment):\n\n"

	for i, choice := range m.name {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if choice == m.selected {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\n按Entry, ctrl+c, q 确认? \n展示Deployment的详细信息，包括Pod, service.\n"
	return s
}

/**
 * ip:port/api/v1/deployment/<namespace>?filterBy=name,audio&itemsPerPage=10&name=&page=1&sortBy=d,creationTimestamp
 */
func (m *DeploymentModel) DeploymentList(namespace, filter string, pageNumber, pageSize int) (value DeploymentResponse) {

	if filter != "" {
		filter = fmt.Sprintf("name,%s", filter)
	}
	url := fmt.Sprintf("https://%s:%d/api/v1/deployment/%s?filterBy=%s&itemsPerPage=%d&name=&page=%d&sortBy=d,creationTimestamp", m.req.Ip, m.req.Port, namespace, filter, pageSize, pageNumber)

	data, err := commonRequest(url, false, nil, true, true, nil)
	if err != nil {
		panic(err)
	}

	var response DeploymentResponse
	json.Unmarshal([]byte(data), &response)
	value = response
	return
}

/**
 * 获取deployment 详情
 */
func (m *DeploymentModel) DeploymentDetail(namespace, name string) (deployment Deployment) {
	url := fmt.Sprintf("https://%s:%d/api/v1/deployment/%s/%s", m.req.Ip, m.req.Port, namespace, name)

	data, err := commonRequest(url, false, nil, false, true, nil)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(data))
	var response Deployment
	json.Unmarshal([]byte(data), &response)
	deployment = response
	return
}

func ShowDeployment(namespace, filter string, req *Request) {
	m := DeploymentModel{
		BaseModel: BaseModel{
			namespace:  namespace,
			filter:     filter,
			pageNumber: defaultPageNumber,
			pageSize:   defaultPageSize,
			req:        req,
		},
	}
	cmd := tea.NewProgram(&m)
	if err := cmd.Start(); err != nil {
		fmt.Println("start failed:", err)
		os.Exit(1)
	}
	if m.selected == "" {
		fmt.Println("没有选择deployment，本次结束")
		return
	}
	d := m.DeploymentDetail(m.namespace, m.selected)
	showDeploymentDetailWithTable(d)

	s := ServiceModel{
		BaseModel{
			namespace: namespace,
			req:       req,
		},
	}
	ShowServiceDetails(&s, d)
}

/*
*
deployment
deployment-name	deployment create timestamp	image	status	pod	pod create timestamp
discovery-audio-live	2021-11-05T08:20:01Z	ccr.ccs.tencentyun.com/en-testing/discovery-audio-live-service:cba9eaa-5829-4a12-97a6-581370e9174c	running(1)/pending(0)/failed(0)/succeened(0)	discovery-audio-live-7b8b6779bf-xqt59	2021-11-25T09:05:14Z
discovery-audio-live-7b8b6779bf-xqt59	2021-11-25T09:05:14Z
*/
func showDeploymentDetailWithTable(deployment Deployment) {
	size := len(deployment.PodList.Pods)
	data := make([][]string, size)
	if size == 0 {
		return
	} else if size >= 1 {
		for i := 0; i < size; i++ {
			column := make([]string, 5)
			if i == 0 {
				column[0] = "1"
				column[1] = deployment.ObjectMeta.Name
				column[2] = strings.Join(deployment.NewReplicaSet.ContainerImages, "\n")
				//status := deployment.PodList.Status
				//column[3] = fmt.Sprintf("running(%d)/pending(%d)/failed(%d)/succeeded(%d)", status.Running, status.Pending, status.Failed, status.Succeeded)

			} else {
				column[0] = ""
				column[1] = ""
				column[2] = ""
			}
			column[3] = deployment.PodList.Pods[i].ObjectMeta.Name
			column[4] = deployment.PodList.Pods[i].ObjectMeta.CreationTimestamp.Format("2006-01-02 15:04:05")
			data[i] = column
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NUMBER", "DEPLOYMENT_NAME", "IMAGE", "POD", "POD_CREATE_TIME"})
	table.SetFooter([]string{"", "", "", "Total", strconv.Itoa(size)})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
}

// 重启pod
func RestartPod(req *Request, namespace, filter string) (err error) {
	p := DeploymentModel{
		BaseModel: BaseModel{
			namespace: namespace,
			filter:    filter,
			req:       req,
		},
	}

	cmd := tea.NewProgram(&p)
	if err := cmd.Start(); err != nil {
		fmt.Println("start failed:", err)
		os.Exit(1)
	}
	if p.selected == "" {
		fmt.Println("没有选择Pod，本次结束")
		return
	}
	// scale Desired to zero, and zero to Desired
	d := p.deployment[p.selected]
	p.RestartPod(namespace, d.ObjectMeta.Labels.App, d.PodStatus.Desired, d.PodStatus.Current)
	return
}

func (m *DeploymentModel) RestartPod(namespace, name string, desited, running int) {
	// fetch current pod size
	fmt.Printf("当前Pod共: %d, 正在运行的pod: %d\n", desited, running)
	// 缩减 pod 的副本数为 0
	m.scalePod(m.namespace, name, 0)
	util.WaitAndShowMessagef(3, time.Second*1, "缩小Pod副本数,等待系统%d秒.")
	// 还原pod 的副本数
	m.scalePod(m.namespace, name, desited)
	util.WaitAndShowMessagef(3, time.Second*1, "扩容Pod副本数,等待系统%d秒.")
}

// 重启pod
func ScalePod(req *Request, namespace, name string, scale int) (err error) {
	p := DeploymentModel{
		BaseModel: BaseModel{
			namespace: namespace,
			filter:    name,
			req:       req,
		},
	}

	// 不是重启，进行pod副本scale
	fmt.Println("当前POD的name:  ", name)
	p.scalePod(namespace, name, scale)
	return
}

// https://domain.../api/v1/scale/deployment/stage-2/payment?scaleBy=2
func (m *DeploymentModel) scalePod(namespace, name string, scale int) (error, string) {
	url := fmt.Sprintf("https://%s:%d/api/v1/scale/deployment/%s/%s?scaleBy=%d", m.req.Ip, m.req.Port, namespace, name, scale)
	data, err := commonRequestV2(url, PUT, nil, true, true, nil)
	if err != nil {
		return err, "请求错误"
	}
	var scaleResult ScaleResult
	json.Unmarshal([]byte(data), &scaleResult)
	fmt.Printf("Pod的最终副本数: %d, 当前已经存在pod副本数: %d\n", scaleResult.DesiredReplicas, scaleResult.ActualReplicas)
	return nil, ""
}
