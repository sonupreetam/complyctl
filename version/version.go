// SPDX-License-Identifier: Apache-2.0

package version

import (
	"fmt"
	"io"
	"runtime"
	"text/template"
)

var (
	// commit is the head commit from git. It should be set during build via -ldflags.
	commit string
	// buildDate in ISO8601 format. IT should with ldflags using $(date -u +'%Y-%m-%dT%H:%M:%SZ').
	buildDate string
	// version describes the version of the client
	// set at build time or detected during runtime.
	version string
	// state of git tree, either "clean" or "dirty"
	gitTreeState string
)

type clientVersion struct {
	Platform  string
	Version   string
	GitCommit string
	GoVersion string
	BuildDate string
}

var versionTemplate = `Version:	{{ .Version }}
Go Version:	{{ .GoVersion }}
Git Commit:	{{ .GitCommit }}
Build Date:	{{ .BuildDate }}
Platform:	{{ .Platform }}
`

// WriteVersion will output the templated version message.
func WriteVersion(writer io.Writer) error {
	if version == "" {
		version = "v0.0.0-unknown"
	}
	versionWithState := version
	if gitTreeState != "" {
		versionWithState = fmt.Sprintf("%s+%s", version, gitTreeState)
	}

	versionInfo := clientVersion{
		Version:   versionWithState,
		GitCommit: commit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}

	tmp, err := template.New("version").Parse(versionTemplate)
	if err != nil {
		return fmt.Errorf("template parsing error: %v", err)
	}

	return tmp.Execute(writer, versionInfo)
}
