package frontend

import (
	"embed"
	"io/fs"
)

func init() {
	var err error
	Assets, err = fs.Sub(embedAssets, "assets")
	if err != nil {
		panic(err)
	}
}

//go:embed assets/*
var embedAssets embed.FS

var Assets fs.FS
