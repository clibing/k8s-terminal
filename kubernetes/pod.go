package kubernetes

import (
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type PodApi interface {
	PodList(namespace, filter string, pageNumber, pageSize int) (value PodListResponse)
	PodDetail(namespace, name string) (pod Pod)
}

type PodModel struct {
	BaseModel
	pods    map[string]Pod
	enable  bool
	tail    int
	restart bool
}

func (m *PodModel) Init() tea.Cmd {
	data := m.PodList(m.namespace, m.filter, m.pageNumber, m.pageSize)
	m.total = data.Status.Running + data.Status.Failed + data.Status.Pending + data.Status.Succeeded
	m.name = make([]string, 0)
	m.pods = make(map[string]Pod)
	m.pageNumber = defaultPageNumber
	m.pageSize = defaultPageSize
	for _, dep := range data.Pods {
		name := dep.ObjectMeta.Name
		m.pods[name] = dep
		m.name = append(m.name, name)
	}
	return nil
}

func (m *PodModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *PodModel) View() string {
	s := "Pod name list(按空格选择Pod):\n\n"

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

	s += "\n按Entry, ctrl+c, q 确认? \n展示Pod的详细信息, 支持实时查看日志.\n"
	return s
}

func (m *PodModel) PodList(namespace, filter string, pageNumber, pageSize int) (value PodListResponse) {
	if filter != "" {
		filter = fmt.Sprintf("name,%s", filter)
	}

	// https://domain.../api/v1/deployment/stage-2?filterBy=name,payment&itemsPerPage=50&name=&page=1&sortBy=d,creationTimestamp
	url := fmt.Sprintf("https://%s:%d/api/v1/pod/%s?filterBy=%s&itemsPerPage=%d&name=&page=%d&sortBy=d,creationTimestamp", m.req.Ip, m.req.Port, namespace, filter, pageSize, pageNumber)

	data, err := commonRequest(url, false, nil, true, true, nil)
	if err != nil {
		panic(err)
	}

	var response PodListResponse
	json.Unmarshal([]byte(data), &response)
	value = response
	return
}

func ShowPod(req *Request, namespace, filter string, log bool, pageSize, tail int, download bool, downloadPath string) (err error) {
	p := PodModel{
		BaseModel: BaseModel{
			namespace:  namespace,
			filter:     filter,
			pageNumber: defaultPageNumber,
			pageSize:   defaultPageSize,
			req:        req,
		},
		enable: log,
		tail:   tail,
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

	if !log && !download {
		fmt.Println("没有开启实时日志功能或者下载功能，自动退出")
		return
	}
	tl := LogModel{
		BaseModel: BaseModel{
			namespace: namespace,
			req:       req,
		},
		pod:      p.selected,
		tail:     tail,
		download: download,
	}
	if log {
		tl.TailLog(req, namespace, p.selected, pageSize, tail)
	}
	if download {
		tl.DownloadLog(downloadPath)
	}
	return
}
