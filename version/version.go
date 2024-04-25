package version

import (
	_ "embed"
	"strings"

	v "github.com/hashicorp/go-version"
)

//go:embed VERSION
var version string

var Version string

var SemVer *v.Version

func init() {
	verFull := v.Must(v.NewVersion(strings.TrimSpace(version)))
	Version = verFull.String()
}
