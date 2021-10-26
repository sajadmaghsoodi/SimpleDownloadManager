package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/sajadmaghsoodi/downloadManager/Utils/downloader"
)

const (
	Command_Download = "Download"
	Command_Quit     = "Quit"
	Command_Help     = "Help"
)

var commands [3]string = [3]string{Command_Download, Command_Quit, Command_Help}

func main() {
	ReadCommands()
	/*currentDownloader := downloader.NewFromURL(GetFileDownloadAddress())
	url := currentDownloader.GetURL()
	currentDownloader.FetchSize()
	currentDownloader.SetThreadCount(10)
	currentDownloader.SetDownloadPath(path.Base(url))
	currentDownloader.Download()

	fmt.Printf("\nDownload Finished")
	*/
}

func ReadCommands() {
	for {
		println("Enter command :")
		var input string
		validCommand := false
		fmt.Scanf("%s", &input)

		for i := range commands {
			if strings.Contains(input, commands[i]) {
				validCommand = true
				ExecuteCommand(i)
			}
		}

		if !validCommand {
			println("command not found. type %s to see the commands", Command_Help)
		}

	}
}

func ExecuteCommand(commandIndex int) {
	if len(commands) < commandIndex || commandIndex < 0 {
		println("command not found. type %s to see the commands", Command_Help)
		return
	}

	ClearConsole()

	switch commands[commandIndex] {

	case Command_Download:
		SetupDownload()

	case Command_Quit:
		os.Exit(3)

	case Command_Help:
		ShowCommandsList()
	}
}

func SetupDownload() {
	currentDownloader := downloader.NewFromURL(GetFileDownloadAddress())
	url := currentDownloader.GetURL()
	currentDownloader.FetchSize()
	currentDownloader.SetThreadCount(10)
	currentDownloader.SetDownloadPath(path.Base(url))
	currentDownloader.Download()

	fmt.Printf("\nDownload Finished\n")
}

func ShowCommandsList() {
	println("Commands :")
	for i := range commands {
		println("   -" + commands[i])
	}
}

func GetFileDownloadAddress() string {
	var input string
	fmt.Printf("Enter the file address to download : \n")
	fmt.Scanf("%s", &input)

	return input
}

func ClearConsole() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
