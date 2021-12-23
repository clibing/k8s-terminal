package kubernetes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type SecretApi interface {
	/**
	 * secret 列表, 支持查询
	 */
	SecretList(namespace, filter string, pageNumber, pageSize int) (value SecretListResponse)
	/**
	 * 获取secret详细信息
	 */
	SecretDetail(namespace, name string) (value SecretResponse)
}

type SecretModel struct {
	BaseModel
	secret map[string]Secret // Secret 详情 这个玩意没有顺序
}

func (m *SecretModel) Init() tea.Cmd {
	data := m.SecretList(m.namespace, m.filter, m.pageNumber, m.pageSize)
	m.name = make([]string, 0)
	m.secret = make(map[string]Secret)
	m.pageNumber = defaultPageNumber
	m.pageSize = defaultPageSize
	for _, se := range data.Secret{
		name := se.ObjectMeta.Name
		m.secret[name] = se
		m.name = append(m.name, name)
	}
	return nil
}

func (m *SecretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "enter":
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

func (m *SecretModel) View() string {
	s := "secret name list(按空格选择secret):\n\n"

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

	s += "\n按Entry, ctrl+c, q 确认? \n展示Secret的详细信息，包括Pod, service.\n"
	return s
}

func (m *SecretModel) SecretList(namespace, filter string, pageNumber, pageSize int) (value SecretListResponse) {

	if filter != "" {
		filter = fmt.Sprintf("name,%s", filter)
	}
	url := fmt.Sprintf("https://%s:%d/api/v1/secret/%s?filterBy=%s&itemsPerPage=%d&name=&page=%d&sortBy=d,creationTimestamp", m.req.Ip, m.req.Port, namespace, filter, pageSize, pageNumber)

	data, err := commonRequest(url, false, nil, true, true, nil)
	if err != nil {
		panic(err)
	}

	var response SecretListResponse
	json.Unmarshal([]byte(data), &response)
	value = response
	return
}

/**
 * 获取secret 详情
 */
func (m *SecretModel) SecretDetail(namespace, name string) (secret SecretResponse) {
	url := fmt.Sprintf("https://%s:%d/api/v1/secret/%s/%s", m.req.Ip, m.req.Port, namespace, name)

	data, err := commonRequest(url, false, nil, false, true, nil)
	if err != nil {
		panic(err)
	}

	var response SecretResponse
	json.Unmarshal([]byte(data), &response)
	secret = response
	return
}

func ShowSecret(namespace, filter string, req *Request) {
	m := SecretModel{
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
		fmt.Println("没有选择Secret，本次结束")
		return
	}
	d := m.SecretDetail(m.namespace, m.selected)

	for k,v := range d.Data {
		fmt.Printf("---------------------------------------[%s] start------------------------------------------------------\n", k)
		result, _ := base64.StdEncoding.DecodeString(v)
		fmt.Println(string(result))
		fmt.Printf("---------------------------------------[%s] end ------------------------------------------------------\n\n", k)
	}
}

