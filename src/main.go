package main

import (
	"bhumit070/destiny/src/config"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type DirectoryFileList = map[string][]string

func main() {
	args := os.Args[1:]
	directory := validateInput(args)

	var directoryFileList DirectoryFileList = make(map[string][]string)

	userInput(directory)
	findFilesInDirectory(directory, &directoryFileList)
	alterDirectory(directory, &directoryFileList)
	showStats(&directoryFileList)

}

func validateInput(args []string) string {
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
	return directory
}

func userInput(directory string) {
	fmt.Print("Are you sure it will alter the files in the " + directory + "? (y/N)")
	var answer string
	_, inputReadError := fmt.Scanln(&answer)

	if inputReadError != nil {
		log.Fatal(inputReadError.Error())
		os.Exit(1)
	}

	answer = strings.TrimSpace(answer)

	if answer != "y" {
		fmt.Println("Aborting")
		os.Exit(0)
	}
}

func findFilesInDirectory(directory string, directoryFileList *DirectoryFileList) {
	fileList, err := os.ReadDir(directory)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

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
		(*directoryFileList)[fileExtension] = append((*directoryFileList)[fileExtension], fileName)
	}
}

func alterDirectory(directory string, directoryFileList *DirectoryFileList) {
	var wg sync.WaitGroup
	for fileExtension, fileNames := range *directoryFileList {
		wg.Add(1)
		go func(fileExtension string, fileNames []string) {
			defer wg.Done()
			createDirectoryError := os.MkdirAll(fileExtension, 0755)
			if createDirectoryError != nil {
				fmt.Println(createDirectoryError.Error())
				os.Exit(1)
			}
			for _, fileName := range fileNames {
				source := path.Join(directory, fileName)
				destination := path.Join(directory, fileExtension, fileName)

				// check if destination exists
				_, err := os.Stat(destination)
				if err == nil {
					splitFileWithExtension := strings.Split(fileName, ".")

					destination = ""

					i := 0

					for i < len(splitFileWithExtension)-1 {
						if i == len(splitFileWithExtension)-2 {
							destination += splitFileWithExtension[i] + "_" + fmt.Sprintf("%d", time.Now().Unix())
						} else {
							destination += (splitFileWithExtension[i] + ".")
						}
						i += 1
					}

					destination += "." + splitFileWithExtension[len(splitFileWithExtension)-1]

					destination = path.Join(directory, fileExtension, destination)
				}

				os.Rename(source, destination)
			}
		}(fileExtension, fileNames)
	}
	wg.Wait()
}

func showStats(directoryFileList *DirectoryFileList) {
	totalFiles := 0
	for fileExtension, fileNames := range *directoryFileList {
		totalFiles += len(fileNames)
		fmt.Println("Moved ", len(fileNames), " ", fileExtension, " file(s) to "+fileExtension)
	}
	fmt.Println("Total files moved: ", totalFiles)
}
