package cli

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func installCommand() *cli.Command {
	return &cli.Command{
		Name:    install.Command,
		Aliases: []string{install.Abbreviations},
		Usage:   "安装",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "安装到指定目录",
			},
		},
		Action: func(c *cli.Context) error {
			len := len(c.FlagNames())
			path := ""
			if len > 0 {
				for _, v := range c.FlagNames() {
					if v == "path" || v == "p" {
						path = c.String(v)
						break
					}
				}
			}
			Install(path)
			return nil
		},
	}
}

func Install(path string)  {
	goos := runtime.GOOS
	if path == "" {
		if goos == "window" {
			// todo
			path = "C://Windows/System32"
		} else if goos == "linux" || goos == "darwin" {
			path = "/usr/local/bin"
		}
	}

	binPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		fmt.Errorf("failed to get bin file info: %s: %s", os.Args[0], err)
		return
	}

	currentFile, err := os.Open(binPath)
	if err != nil {
		fmt.Errorf("failed to get bin file info: %s: %s", binPath, err)
		return
	}
	defer func() { _ = currentFile.Close() }()

	installFile, err := os.OpenFile(filepath.Join(path, "k8s-terminal"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Errorf("failed to create bin file: %s: %s", filepath.Join(path, "knife"), err)
		return
	}
	defer func() { _ = installFile.Close() }()

	_, err = io.Copy(installFile, currentFile)
	if err != nil {
		fmt.Errorf("failed to copy file: %s: %s", filepath.Join(path, "knife"), err)
		return
	}
	fmt.Println("install success")
	return
}