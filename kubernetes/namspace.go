package kubernetes

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/k8s-terminal/config"
	"os"
	"sort"
)

/**
 * load namespace  request
 */

type NamespaceApi interface {
	/**
	 * namespace 列表, 支持查询
	 */
	NamespaceList(filter string) (value NamespaceResponse)
}

type NamespaceModel struct {
	BaseModel
	cursor  int
	checked map[int]string
}

func (m *NamespaceModel) NamespaceList(filter string) (value NamespaceResponse) {
	if filter != "" {
		filter = fmt.Sprintf("name,%s", filter)
	}
	url := fmt.Sprintf("https://%s:%d/api/v1/namespace", m.req.Ip, m.req.Port)

	data, err := commonRequest(url, false, nil, false, true, nil)
	if err != nil {
		panic(err)
	}

	var response NamespaceResponse
	json.Unmarshal([]byte(data), &response)
	value = response
	return
}

func (m *NamespaceModel) Init() tea.Cmd {
	data := m.NamespaceList(m.filter)
	for _, ns := range data.Namespaces {
		m.name = append(m.name, ns.ObjectMeta.Name)
	}
	m.checked = make(map[int]string)
	return nil
}

func (m *NamespaceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			_, ok := m.checked[m.cursor]
			if ok {
				delete(m.checked, m.cursor)
			} else {
				m.checked[m.cursor] = ""
			}
		}
	}
	return m, nil
}

func (m *NamespaceModel) View() string {
	s := "Namespace list(按空格选择):\n\n"

	for i, choice := range m.name {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.checked[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\n按Entry, ctrl+c, q 确认? \n展示Namespace的详细信息.\n"
	return s
}

func ShowNamespace(filter string, save bool, req *Request) {
	m := NamespaceModel{
		BaseModel: BaseModel{
			filter: filter,
			req:    req,
		},
	}
	cmd := tea.NewProgram(&m)
	if err := cmd.Start(); err != nil {
		fmt.Println("start failed:", err)
		os.Exit(1)
	}
	if len(m.checked) > 0 && save {
		var t []string
		var key []int
		for k, _ := range m.checked {
			key = append(key, k)
		}
		sort.Ints(key)
		for _, v := range key {
			t = append(t, m.name[v])
		}
		config.GlobalCfg.Namespace.Names = t
		config.GlobalCfg.Namespace.Static = true
		config.PersistenceConfigFile(&config.GlobalCfg, false)
	}
}
