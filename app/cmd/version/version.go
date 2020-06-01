package version

import (
	"fmt"
	"github.com/urfave/cli"
	"runtime"
	"strings"
)

const tpl = `
  Version:        %s
  Go version:     %s
  Git commit:     %s
  Built:          %s
  OS/Arch:        %s/%s
`

func Command(info *BuildInfo) cli.Command {
	return cli.Command{
		Name:  "version",
		Usage: "Shows app version",
		Action: func(ctx *cli.Context) {
			fmt.Println(strings.Trim(fmt.Sprintf(tpl,
				info.GitTag,
				runtime.Version(),
				info.GitCommit,
				info.BuildTime,
				runtime.GOOS, runtime.GOARCH,
			), "\n"))
		},
	}
}

type BuildInfo struct {
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
	GitTag    string `json:"git_tag"`
}
