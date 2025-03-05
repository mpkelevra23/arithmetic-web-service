package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

// GetFileSystem возвращает HTTP файловую систему для веб-ресурсов
func GetFileSystem() (http.FileSystem, error) {
	// Получить подфайловую систему из директории static
	fsys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, err
	}
	return http.FS(fsys), nil
}
