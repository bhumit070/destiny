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

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DirectoryFileList = map[string][]string

func main() {
	args := os.Args[1:]
	flags := parseFlags(&args)
	directory := validateInput(args)

	var directoryFileList DirectoryFileList = make(map[string][]string)

	userInput(directory, &flags)
	findFilesInDirectory(directory, &directoryFileList)
	checkForFolderGroups(&directoryFileList, &flags)
	alterDirectory(directory, &directoryFileList)
	showStats(&directoryFileList, &flags)

}

func checkForFolderGroups(directoryFileList *DirectoryFileList, flags *config.InputFlags) {

	if config.IsFlagExists("nfg", flags) {
		return
	}

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
		if config.IsValidFlag(strings.Trim(match[1], "")) {
			trimmedFlag := strings.Trim(match[1], " ")
			flags[trimmedFlag] = trimmedFlag
		} else {
			_flags := strings.Split(match[1], "")
			for _, _flag := range _flags {
				trimmedFlag := strings.Trim(_flag, " ")
				if config.IsValidFlag(trimmedFlag) {
					flags[trimmedFlag] = trimmedFlag
				}
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

	userHomeDirectory, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	if config.IsFlagExists("y", flags) && directory != userHomeDirectory {
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

		fileExtension := path.Ext(fileName)

		if fileExtension == "" {
			fileExtension = "others"
		} else {
			fileExtension = fileExtension[1:]
			fileExtension = strings.ToLower(fileExtension)
		}
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

				// check if destination folder exists or not
				destinationFolderPath := path.Join(directory, fileExtension)
				_, destinationDirError := os.Stat(destinationFolderPath)

				if destinationDirError != nil {
					if os.IsNotExist(destinationDirError) {
						os.MkdirAll(destinationFolderPath, os.ModePerm)
					}
				}

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

				fmt.Sprintln("moved %s to %s", source, destination)
				os.Rename(source, destination)
			}
		}(fileExtension, fileNames)
	}
	wg.Wait()
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "down":
			m.table.MoveDown(1)
		case "up":
			m.table.MoveUp(1)
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func showStats(directoryFileList *DirectoryFileList, flags *config.InputFlags) {

	if config.IsFlagExists("q", flags) {
		return
	}

	totalFiles := 0

	var rows []table.Row
	for fileExtension, fileNames := range *directoryFileList {
		totalFiles += len(fileNames)
		str := fmt.Sprintf("%d", len(fileNames))
		rows = append(rows, table.Row{
			fileExtension,
			str,
		})
	}

	if totalFiles == 0 {
		fmt.Println("No files moved")
		return
	}

	columns := []table.Column{
		{
			Title: "File Extension",
			Width: 20,
		},
		{
			Title: "Moved Count",
			Width: 20,
		},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(len(rows)+1),
	)
	m := model{table: t}
	time.AfterFunc(1*time.Second, func() {
		os.Exit(0)
	})
	tea.NewProgram(m).Run()
}
