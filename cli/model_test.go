package cli

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
)

func TestModel(t *testing.T) {
	goarch := runtime.GOARCH
	goos := runtime.GOOS
	fmt.Println("返回当前的系统架构: ", goarch, "返回当前的操作系统: ", goos)

	var name []string
	name = make([]string, 100)
	var i int
	for i=0; i<100;i++{
		name[i] = strconv.Itoa(i)
	}
	fmt.Println(name[58])
	fmt.Println(name[59])
}
