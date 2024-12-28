package index

import (
	"embed"
	"io/fs"
)

//go:embed web
var assets embed.FS

func GetWebFS() fs.FS {
	serverRoot, _ := fs.Sub(assets, "web")
	return serverRoot
}
