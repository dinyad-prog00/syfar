package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	v "syfar/runner"
)

func buildPath(args []string) (string, string, error) {
	var path string
	if len(args) >= 1 {
		path = args[0]
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

// CopyDir copies a directory recursively from src to dst.
func CopyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}

	err = os.MkdirAll(dst, info.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func initSyfarJson(projectDir string) error {
	syfarJsonPath := filepath.Join(projectDir, "syfar.json")
	if fileExists(syfarJsonPath) {
		return fmt.Errorf("project is already created")
	}
	syfarJson := make(map[string]interface{})
	syfarJson["version"] = "1.0.0"
	syfarJson["syfar_version"] = v.Version

	jsonData, err := json.MarshalIndent(syfarJson, "", "\t")
	if err != nil {
		return err
	}
	err = os.MkdirAll(projectDir, 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(syfarJsonPath, jsonData, 0755)
	if err != nil {
		return err
	}

	return nil
}
