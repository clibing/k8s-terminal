package kubernetes

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type LogModel struct {
	BaseModel
	pod                                  string
	offsetFrom, offsetTo, pageSize, tail int
	first, download                      bool
	containerName                        string
}

const (
	PRINT  int8 = 0
	SEARCH int8 = 1
	MANUAL int8 = 2
)

var (
	logData                       []Log // 实时日志列表
	LM                            *LogModel
	LogQueue                      = make(chan Log, 10000)
	mode                          = PRINT
	keyValue                      []string
	lock                          *sync.Mutex = new(sync.Mutex)
	allowNextRequest              *sync.Mutex = new(sync.Mutex)
	autoLock                                  = new(sync.RWMutex)
	input                         bool
	auto                          = true
	start, end, searchSnapshotEnd int         // 当开启搜索的时候 需要记录 当前切片的最长 防止数据溢出
	searchData                    []SearchLog // 存储搜索后的结果
	logFilePosition               = "end"
	line                          = 5
	current                       = 0
)

type SearchLog struct {
	index int
	log   Log
}

func (m *LogModel) Init() tea.Cmd {
	go showLog(LogQueue)
	go fetLog()
	return nil
}

func showLog(ch chan Log) {
	for log := range ch {
		if isAuto() {
			fmt.Println(log.Content)
		}
	}
}

func isAuto() (value bool) {
	autoLock.RLock()
	value = auto
	autoLock.RUnlock()
	return
}

func SetAuto(value bool) {
	autoLock.Lock()
	auto = value
	autoLock.Unlock()
}

func fetLog() {
	logMap := make(map[string]interface{}, 1000) // log的map结构 保证记录只有一条

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		fmt.Println("d")
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容
		}
		fmt.Println("e")
	}()
	loadLogFirst := true
	for {
		data := LM.Log(LM.namespace, LM.pod, LM.offsetFrom, LM.offsetTo, LM.pageSize, loadLogFirst)
		if loadLogFirst {
			LM.containerName = data.Info.ContainerName
		}
		for _, log := range data.Logs {
			_, ok := logMap[log.Timestamp.String()]
			if ok {
				continue
			}
			logMap[log.Timestamp.String()] = &log
			LogQueue <- log
			time.Sleep(time.Millisecond * 100)
			logData = append(logData, log)
		}
		time.Sleep(time.Second * 5)
		loadLogFirst = false
	}
}

/**
 * 接收搜索的关键字
 */
func appendSearchValue(value string) {
	bytes := []byte(value)
	if len(bytes) > 1 {
		return
	}
	v := int(bytes[0])
	// https://blog.csdn.net/itzyjr/article/details/102867808
	if v >= 33 && v <= 126 {
		fmt.Print(value)
		keyValue = append(keyValue, value)
	}
}

/**
 * 搜索结果 展示目标的上下行 行数 (增加或者减少)
 * 5 10 20 50
 */
func incrementLine(isAdd bool) {
	if isAdd {
		switch line {
		case 0:
			line = 5
		case 5:
			line = 10
		case 10:
			line = 20
		case 20:
			line = 50
		case 50:
			fmt.Println("最大支持显示当前目标行的上下50行")
		}
	} else {
		switch line {
		case 5:
			line = 0
		case 10:
			line = 5
		case 20:
			line = 10
		case 50:
			line = 20
		}
	}
}

/**
 * 展示 是否上一行还是下一行
 */
func showSearchNext(next bool) {
	// 搜索结果集列表总和
	searchLen := len(searchData) - 1
	maxLen := len(logData) - 1
	if searchLen <= 0 || maxLen <= 0 {
		fmt.Println("暂无数据")
		return
	}

	if current > searchLen {
		current = searchLen
	}
	if current < 0 {
		current = 0
	}

	tmp := searchData[current]
	i := tmp.index
	start = i - line
	end = i + line
	if start < 0 {
		start = 0
	}
	if end > maxLen {
		end = maxLen
	}

	if start > end {
		return
	}
	d := logData[start:end]
	fmt.Printf("(%d/%d)>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>--------\n", current, searchLen)
	for _, v := range d {
		fmt.Println(v.Content)
		time.Sleep(time.Millisecond * 100)
	}
	fmt.Println("----------------------------------------------------------------------------<<<<<<<<<<<<<<<<<")

	// 向下
	if next {
		current = current + 1
	} else {
		current = current - 1
	}

	if current > searchLen {
		fmt.Println("暂无数据")
		return
	}

	for _, s := range searchData[current:searchLen] {
		if end > s.index {
			current = s.index
		}
	}

}

func help() {
	v := `help:
公共命令
  * ctrl+c或者q: 退出当前Log
  * h: 帮助文档
自动模式
  * c: 将模式切换为自动模式，实时接收微服务的日志信息
搜索模式
  * /: 开启搜索模式 接收输入的内容
  * enter: 回车，当开启搜素模式，将输入的最为关键字
  * n: 搜索到的结果，下翻
  * N: 搜索到的结果，上翻
  * ->: 当搜索到关键字的位置后，可以展开上下的行数, 增加 [5 10 20 50 100 200] 是一个循环递增的过程
  * <-: 当搜索到关键字的位置后，可以展开上下的行数, 减少 [5 10 20 50 100 200] 是一个循环递增的过程
`
	fmt.Println(v)
}
func (m *LogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// 搜索模式并且为输入文本模式 输入的字符为q不处理
			switch mode {
			case SEARCH:
				if input == true {
					if msg.String() == "q" {
						appendSearchValue(msg.String())
					}
				}
			case PRINT:
				return m, tea.Quit
			case MANUAL:
				SetAuto(true)
				input = false
				mode = PRINT
				fmt.Println("从手动模式切换到打印模式")
			}
		case "n", "N":
			switch mode {
			case SEARCH:
				if input == true {
					appendSearchValue(msg.String())
				} else {

				}
			case MANUAL:
				showSearchNext(msg.String() == "n")
			}
		case "right":
			switch mode {
			case MANUAL:
				//fmt.Println("当前为搜索模式， 根据查询到的点 展示点的上下20行, 增加 5 10 20 50 100 200")
				incrementLine(true)
			}

		case "left":
			switch mode {
			case MANUAL:
				//fmt.Println("当前为搜索模式， 根据查询到的点 展示点的上下20行, 减少")
				incrementLine(false)
			}

		case "/":
			if input == false {
				mode = SEARCH
				input = true
				fmt.Println("当前为所搜模式, 输入关键字回车即可！")
				SetAuto(false)
				searchData = searchData[0:0]
				searchSnapshotEnd = 0
			}
		case "esc":

		case "c":
			switch mode {
			case SEARCH:
				if input == true {
					appendSearchValue(msg.String())
				}
			case MANUAL:
				mode = PRINT
				fmt.Println("当前为数据日志自动输出模式")
			}

		case "enter":
			switch mode {
			case SEARCH:
				if input == true {
					input = false
					mode = MANUAL
					search := strings.Join(keyValue, "")
					fmt.Printf("\n从搜索模式切换到手动模式\n当前搜索的内容: [%s]\n", search)
					keyValue = keyValue[0:0]
					searchSnapshotEnd = len(logData) - 1
					for i := 0; i < searchSnapshotEnd; i++ {
						tmp := logData[i]
						if strings.Index(tmp.Content, search) != -1 {
							searchData = append(searchData, SearchLog{
								index: i,
								log:   tmp,
							})
						}
					}
				}
			case PRINT:
				fmt.Println("")
			}
		case "h":
			switch mode {
			case SEARCH:
				if input == true {
					appendSearchValue(msg.String())
				} else {
					help()
				}
			case PRINT:
				fallthrough
			case MANUAL:
				help()
			}
		default:
			// 关键位置 支持搜索输入的结果
			if mode == SEARCH && input == true {
				appendSearchValue(msg.String())
			} else {
				//fmt.Println("没有找到对应的处理 当前 按键是: ", msg.String(), " mode: ", mode, " input: ", input)
			}
		}
	}
	return m, nil
}

func (m *LogModel) View() string {
	//s := "Log日志选项:\n\n"
	//
	//for i, choice := range m.name {
	//	cursor := " "
	//	if m.cursor == i {
	//		cursor = ">"
	//	}
	//
	//	checked := " "
	//	if choice == m.selected {
	//		checked = "x"
	//	}
	//
	//	s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	//}
	//
	//s += "\n按Entry, ctrl+c, q 确认? \n展示Pod的详细信息, 支持实时查看日志.\n"
	//return s
	return ""
}

func (m *LogModel) Log(namespace, name string, offsetFrom, offsetTo, pageSize int, first bool) (value LogResponse) {
	url := ""
	var response LogResponse
	if first {
		url = fmt.Sprintf("https://%s:%d/api/v1/log/%s/%s", m.req.Ip, m.req.Port, namespace, name)
	} else {
		if isAuto() {
			url = fmt.Sprintf("https://%s:%d/api/v1/log/%s/%s/%s?logFilePosition=%s&offsetFrom=%d&offsetTo=%d&previous=false&referenceLineNum=0&referenceTimestamp=newest", m.req.Ip, m.req.Port, namespace, name, m.containerName, logFilePosition, offsetFrom, offsetTo)
		} else {
			//fmt.Println("不是自动模式，需要提供正确的URL")
			//url = fmt.Sprintf("https://%s:%d/api/v1/log/%s/%s/grpc-server?logFilePosition=%s&offsetFrom=%d&offsetTo=%d&previous=false&referenceLineNum=0&referenceTimestamp=newest", m.req.Ip, m.req.Port, namespace, name, logFilePosition, offsetFrom, offsetTo)
			return
		}
	}

	data, err := commonRequest(url, false, nil, false, true, nil)
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(data), &response)
	value = response
	return
}

//
func (m *LogModel) TailLog(req *Request, namespace, name string, pageSize, tail int) {
	LM = &LogModel{
		BaseModel: BaseModel{
			namespace: namespace,
			req:       req,
		},
		pod:        name,
		offsetFrom: 2000000000,
		offsetTo:   2000000100,
		first:      false,
		pageSize:   pageSize,
		tail:       tail,
	}
	fmt.Printf("查看%s/%s的日志\n加载中....\n", namespace, name)
	cmd := tea.NewProgram(LM)
	if err := cmd.Start(); err != nil {
		fmt.Println("start failed:", err)
		os.Exit(1)
	}
}

func (m *LogModel) DownloadLog(path string) {
	pod := m.Log(m.namespace, m.pod, m.offsetFrom, m.offsetTo, m.pageSize, true)
	m.containerName = pod.Info.ContainerName

	url := fmt.Sprintf("https://%s:%d/api/v1/log/file/%s/%s/%s?previous=false", m.req.Ip, m.req.Port, m.namespace, m.pod, m.containerName)
	data, err := commonRequest(url, false, nil, false, true, nil)
	if err != nil {
		panic(err)
	}

	homedir, err := homedir.Dir()
	if err != nil {
		log.Panicln("home dir err")
		return
	}
	if path == "" {
		path = homedir
	}
	fileName := m.pod + "-" + m.containerName + ".log"
	value := filepath.Join(path, fileName)
	ioutil.WriteFile(value, []byte(data), 0755)
	fmt.Println("日志保存成功！", value)
}
