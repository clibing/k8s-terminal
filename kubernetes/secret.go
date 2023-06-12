package kubernetes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
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

	/**
	 * 保存密文
	 */
	SecretUpdate(namespace, name, secretFile, secretContentBase64 string) error
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
	for _, se := range data.Secret {
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
		case "ctrl+c", "q":
			m.selected = ""
			return m, tea.Quit

		case "enter":
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

	s += "\n1. 按Entry下一步\n2. Ctrl+c或者 q 退出？\n"
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

/**
 * 更新密文
 */
func (m *SecretConfirm) SecretUpdate(namespace, name, secretFile, secretContentBase64 string) error {

	// get https://dashboard-msvc2.test1.bj.yxops.net/api/v1/_raw/secret/namespace/payment/name/wechat-apicerts
	url := fmt.Sprintf("https://%s:%d/api/v1/_raw/secret/namespace/%s/name/%s", m.req.Ip, m.req.Port, namespace, name)
	data, err := commonRequest(url, false, nil, false, true, nil)
	if err != nil {
		panic(err)
	}

	var secretUpdate SecretUpdate
	json.Unmarshal([]byte(data), &secretUpdate)

	secretUpdate.Data[secretFile] = secretContentBase64

	// v, err := json.Marshal(secretUpdate)
	// if err != nil {
	// 	panic("json格式化数据异常")
	// }
	// fmt.Println("secret: ", string(v))

	// put
	url = fmt.Sprintf("https://%s:%d/api/v1/_raw/secret/namespace/%s/name/%s", m.req.Ip, m.req.Port, namespace, name)
	_, err = commonRequestV2(url, PUT, secretUpdate, true, true, nil)
	if err != nil {
		fmt.Println(err.Error())
		panic("保密字段数据修改失败: " + secretFile)
	}

	fmt.Printf("\n修改成功[%s]!", secretFile)
	// 接收到的数据替换对应的key 对应的value
	// fmt.Printf("namespace: %s, name: %s, secret file: %s, value: %s\n", namespace, name, secretFile, secretContentBase64)
	return nil
}

func ShowSecret(namespace, filter string, req *Request, is_edit bool) {
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
	// fmt.Println("m.selected=", m.selected)

	d := m.SecretDetail(m.namespace, m.selected)

	// fmt.Printf("secret response data: %v\n", d)

	// 非编辑模式 直接退出
	if !is_edit {
		for k, v := range d.Data {
			// 进去编辑模式， 需要将内容存储到指定目录
			fmt.Printf("---------------------------------------[%s] start------------------------------------------------------\n", k)
			result, _ := base64.StdEncoding.DecodeString(v)
			fmt.Println(string(result))
			fmt.Printf("---------------------------------------[%s] end ------------------------------------------------------\n\n", k)
		}
		return
	}

	e := SecretEdit{
		BaseModel: BaseModel{
			namespace:  namespace,
			filter:     filter,
			pageNumber: defaultPageNumber,
			pageSize:   defaultPageSize,
			req:        req,
		},
	}

	secretFiles := make(map[string]string)
	for k, v := range d.Data {
		result, _ := base64.StdEncoding.DecodeString(v)
		e.name = append(e.name, k)
		secretFiles[k] = string(result)
	}

	edit := tea.NewProgram(&e)
	if err := edit.Start(); err != nil {
		fmt.Println("start failed:", err)
		os.Exit(1)
	}
	if e.selected == "" {
		fmt.Println("请选择需要修改的密文")
		return
	}

	// fmt.Println("selected value: ", e.selected)
	content := secretFiles[e.selected]

	file, errs := ioutil.TempFile("", e.selected)
	if errs != nil {
		fmt.Println(errs)
		return
	}
	defer os.Remove(file.Name())

	// Write some text to the file
	_, errs = file.WriteString(content)
	if errs != nil {
		fmt.Println(errs)
		return
	}

	// Close the file
	errs = file.Close()
	if errs != nil {
		fmt.Println(errs)
		return
	}

	// fmt.Println("current secret file: ", file.Name())

	vim := exec.Command("vim", file.Name())
	vim.Stdin = os.Stdin
	vim.Stdout = os.Stdout
	err := vim.Run()
	if err != nil {
		fmt.Println("使用vim编辑密文异常")
		return
	}

	c, err := os.ReadFile(file.Name())
	if err != nil {
		fmt.Println("读取编辑后的文件异常")
		return
	}

	// fmt.Println("读取编辑后的文件: \n", string(c))

	sc := SecretConfirm{
		BaseModel: BaseModel{
			namespace:  namespace,
			filter:     filter,
			pageNumber: defaultPageNumber,
			pageSize:   defaultPageSize,
			req:        req,
		},
	}
	sc.name = append(sc.name, "save")
	sc.name = append(sc.name, "cannel")
	sc.selected = "canel"
	confirm := tea.NewProgram(&sc)
	if err := confirm.Start(); err != nil {
		fmt.Println("confirm failed:", err)
		os.Exit(1)
	}
	if sc.selected != "save" {
		fmt.Println("编辑后暂未确认，退出。")
		return
	}

	value := base64.StdEncoding.EncodeToString(c)
	sc.SecretUpdate(namespace, m.selected, e.selected, value)
}

// ----------------------------------------------------------------
type SecretEdit struct {
	BaseModel
	SecreFiles map[string]string // 加密文件 名字 + 文件路径
}

func (m *SecretEdit) Init() tea.Cmd {

	return nil
}

func (m *SecretEdit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.selected = ""
			return m, tea.Quit

		case "enter":
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

func (m *SecretEdit) View() string {
	s := "\n选择待编辑文件 (按空格选择secret):\n\n"

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

	s += "\n1. 按Entry下一步\n2. Ctrl+c或者 q 退出？\n"
	return s
}

// ----------------------------------------------------------------

type SecretConfirm struct {
	BaseModel
	secretName string // 当前编辑的 密文
}

func (m *SecretConfirm) Init() tea.Cmd {
	return nil
}
func (m *SecretConfirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.selected = ""
			return m, tea.Quit

		case "enter":
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

func (m *SecretConfirm) View() string {
	s := fmt.Sprintf("是否确认更新 (%s):\n\n", m.secretName)

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

	s += "\n1. 按Entry提交本地修改\n2. Ctrl+c或者 q 退出？\n"
	return s
}
