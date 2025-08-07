package assets

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var assets embed.FS

func GetAssets() fs.FS {
	files, err := fs.Sub(assets, "dist")
	if err != nil {
		panic("failed to get assets: " + err.Error())
	}
	return files
}
