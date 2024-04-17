package main

import (
	"fmt"
	"os"
	"path/filepath"
	file "syfar/providers/file"
	http "syfar/providers/http"
	runner "syfar/runner"
)

func main() {

	args := os.Args
	var filename string
	if len(args) >= 2 {
		filename = args[1]
	} else {
		panic("syfar: Error - no file is provided")
	}

	wdir, err := os.Getwd()
	if err != nil {
		fmt.Println("syfar: Erreur lors de l'obtention du répertoire de travail:", err)
		return
	}

	filedir := filepath.Dir(filepath.Join(wdir, filename))

	syfar := runner.NewSyfar()

	//File provider
	fileProvider := file.ActionProvider{}
	syfar.RegisterActionProvider("file", &fileProvider)

	//HTTP Provider
	httpProvider := http.ActionProvider{}
	syfar.RegisterActionProvider("http", &httpProvider)

	syfar.Init()
	syfar.Run(filedir, filename)

}
