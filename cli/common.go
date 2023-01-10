package cli

import "github.com/urfave/cli/v2"

const (
	COMMON_NAMESPACE_NAME = "namespace"
	COMMON_NAME_NAME      = "name"
	COMMON_PATH_NAME      = "path"
)

var commonNamespaceFlag = &cli.StringFlag{
	Name:    COMMON_NAMESPACE_NAME,
	Aliases: []string{"ns", "s"},
	Usage:   "所属namespace",
}

var commonNameFlag = &cli.StringFlag{
	Name:    COMMON_NAME_NAME,
	Aliases: []string{"n"},
	Usage:   "名字",
}

var commonPathFlag = &cli.StringFlag{
	Name:    COMMON_PATH_NAME,
	Aliases: []string{"p"},
	Usage:   "路径",
}
