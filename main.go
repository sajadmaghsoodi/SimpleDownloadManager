package main

import (
	"fmt"
	"path"

	"github.com/sajadmaghsoodi/downloadManager/Utils/downloader"
)

func main() {
	currentDownloader := downloader.NewFromURL(GetFileDownloadAddress())
	url := currentDownloader.GetURL()
	currentDownloader.FetchSize()
	currentDownloader.SetThreadCount(10)
	currentDownloader.SetDownloadPath(path.Base(url))
	currentDownloader.Download()

	fmt.Printf("\nDownload Finished")
}

func GetFileDownloadAddress() string {
	var input string
	fmt.Printf("Enter the file address to download : \n")
	fmt.Scanf("%s", &input)

	return input
}
