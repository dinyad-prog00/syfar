package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func buildPath(args []string) (string, string, error) {
	var path string
	if len(args) >= 1 {
		path = args[0]
		path, err := filepath.Abs(path)
		if err != nil {
			return "", "", err
		}
		if !fileExists(path) {
			return "", "", fmt.Errorf("file %s not found", path)
		}
		if isDir(path) {
			return buildDir(path)
		} else {
			return buildFile(path)
		}

	} else {
		wdir, err := os.Getwd()
		if err != nil {
			return "", "", err
		}
		return buildDir(wdir)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func buildDir(path string) (string, string, error) {
	mainSFPath := filepath.Join(path, "main.sf")
	if _, err := os.Stat(mainSFPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("main.sf does not exist in the folder \"%s\"", path)
	}
	return path, mainSFPath, nil
}

func buildFile(path string) (string, string, error) {
	wdir, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	filedir := filepath.Dir(filepath.Join(wdir, path))
	return filedir, path, nil
}
