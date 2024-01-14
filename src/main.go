package main

import (
	"bhumit070/destiny/src/config"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
)

type DirectoryFileList = map[string][]string

func main() {
	args := os.Args[1:]
	flags := parseFlags(&args)
	directory := validateInput(args)

	var directoryFileList DirectoryFileList = make(map[string][]string)

	userInput(directory, &flags)
	findFilesInDirectory(directory, &directoryFileList)
	checkForFolderGroups(&directoryFileList)
	alterDirectory(directory, &directoryFileList)
	showStats(&directoryFileList, &flags)

}

func checkForFolderGroups(directoryFileList *DirectoryFileList) {
	for key, value := range *directoryFileList {
		if _, ok := config.FolderGroups[key]; ok {
			delete(*directoryFileList, key)
			(*directoryFileList)[config.FolderGroups[key]+key] = value

		}
	}
}

func parseFlags(args *[]string) config.InputFlags {
	var inputString string = ""
	for _, arg := range *args {
		inputString += arg + " "
	}

	re := regexp.MustCompile(`-([a-zA-Z]+)`)
	matches := re.FindAllStringSubmatch(inputString, -1)
	var flags = make(map[string]string)
	var indexesToRemoveFromArgs []int = make([]int, 0, len(matches))
	for index, match := range matches {
		_flags := strings.Split(match[1], "")
		for _, _flag := range _flags {
			if _, ok := config.ValidFlags[strings.Trim(_flag, " ")]; ok {
				trimmedFlag := strings.Trim(_flag, " ")
				flags[trimmedFlag] = trimmedFlag
			}
		}
		indexToRemove := -1
		if (*args)[index] == match[0] {
			indexToRemove = index
		} else if (*args)[index+1] == match[0] {
			indexToRemove = index + 1
		}
		indexesToRemoveFromArgs = append(indexesToRemoveFromArgs, indexToRemove)
	}
	if len(indexesToRemoveFromArgs) > 0 {

		for iteration, indexToRemove := range indexesToRemoveFromArgs {
			start := indexToRemove
			if indexToRemove-iteration < 0 {
				start = 0
			} else {
				start = indexToRemove - iteration
			}
			*args = append((*args)[:start], (*args)[indexToRemove+1-iteration:]...)
		}
	}

	return flags
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

func userInput(directory string, flags *config.InputFlags) {

	if isFlagExists("y", flags) {
		return
	}

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

func showStats(directoryFileList *DirectoryFileList, flags *config.InputFlags) {

	if isFlagExists("q", flags) {
		return
	}

	totalFiles := 0
	for fileExtension, fileNames := range *directoryFileList {
		totalFiles += len(fileNames)
		fmt.Println("Moved ", len(fileNames), " ", fileExtension, " file(s) to "+fileExtension)
	}
	fmt.Println("Total files moved: ", totalFiles)
}

func isFlagExists(flag string, flags *config.InputFlags) bool {
	if _, ok := config.ValidFlags[flag]; ok {
		if _, _ok := (*flags)[flag]; _ok {
			return true
		}
	}
	return false
}
