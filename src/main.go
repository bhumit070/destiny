package main

import (
	"bhumit070/destiny/src/config"
	"fmt"
	"os"
	"path"
	"strings"
)

func main() {

	args := os.Args[1:]

	var directory string
	if len(args) > 0 {
		directory = args[0]
		if directory == "" {
			fmt.Println("No directory given")
			os.Exit(1)
		}
	} else {
		_directory, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		directory = _directory
	}

	fileInfo, err := os.Stat(directory)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if !fileInfo.IsDir() {
		fmt.Println(directory, "is not a directory")
		os.Exit(0)
	}

	fileList, err := os.ReadDir(directory)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var directoryFileList map[string][]string = make(map[string][]string)

	for _, file := range fileList {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()

		if config.ExcludedFileList[fileName] == fileName {
			continue
		}

		splitFileName := strings.Split(fileName, ".")
		fileExtension := splitFileName[len(splitFileName)-1]
		directoryFileList[fileExtension] = append(directoryFileList[fileExtension], fileName)
	}

	for fileExtension, fileNames := range directoryFileList {
		createDirectoryError := os.MkdirAll(fileExtension, 0755)
		if createDirectoryError != nil {
			fmt.Println(createDirectoryError.Error())
			os.Exit(1)
		}
		for _, fileName := range fileNames {
			source := path.Join(directory, fileName)
			destination := path.Join(directory, fileExtension, fileName)
			os.Rename(source, destination)
		}
	}

}
