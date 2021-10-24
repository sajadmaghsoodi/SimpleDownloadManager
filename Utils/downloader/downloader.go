package downloader

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	url             string
	size            int64
	threadCount     int
	downloadPath    string
	downloadThreads []Thread
}

var wg sync.WaitGroup

func NewFromURL(url string) *Downloader {
	return &Downloader{
		url: url,
	}
}

func New() *Downloader {
	return &Downloader{}
}

func (d *Downloader) SetDownloadPath(path string) {
	d.downloadPath = path
}

func (d *Downloader) SetThreadCount(count int) {
	d.threadCount = count
}

func (d *Downloader) SetURL(url string) {
	d.url = url
}

func (d *Downloader) GetURL() string {
	return d.url
}

func (d *Downloader) FetchSize() error {
	if d.url == "" {
		return fmt.Errorf("URL is not set")
	}

	response, err := http.Head(d.url)
	if err != nil {
		return fmt.Errorf("error fetching file size: %+v", err)
	}

	d.size = response.ContentLength

	return nil
}

func (d *Downloader) ShowProgress() {

	bar := progressbar.Default(100)
	for {
		bar.Clear()
		for i := 0; i < d.threadCount; i++ {
			bar.Add(int(d.downloadThreads[i].GetProgress() * 10))
		}

		time.Sleep(time.Second * 1)
	}
}

func (d *Downloader) Download() {
	d.startDownloadThreads()
}

func (d *Downloader) startDownloadThreads() {
	wg = sync.WaitGroup{}

	defer d.CreateDownloadedFile()
	defer wg.Wait()

	threadSize := int(d.size) / d.threadCount
	additionBytesForLastThread := int(d.size) % d.threadCount

	wg.Add(d.threadCount)
	for i := 0; i < d.threadCount; i++ {
		min := i * threadSize
		max := (i + 1) * threadSize

		if i == d.threadCount-1 { //if its the last thread it needs to download to the end
			max += additionBytesForLastThread + 1
		}

		thread := *NewThread(d.url, min, max, i, d)
		d.downloadThreads = append(d.downloadThreads, thread)

		go thread.Start()
	}

	go d.ShowProgress()
}

func (d *Downloader) CreateDownloadedFile() {

	result := make([]byte, 0)
	for i := 0; i < d.threadCount; i++ {

		tmpThreadData, err := ioutil.ReadFile(strconv.Itoa(i))
		if err != nil {
			println(err)
		}
		result = append(result, tmpThreadData...)
		os.Remove(strconv.Itoa(i))
	}

	fileName := d.downloadPath
	suf := ""

	for i := 0; true; i++ {
		fileExists, _ := fileExists(suf + fileName)

		if !fileExists {
			break
		} else {
			suf = fmt.Sprintf("(%v) ", strconv.Itoa(i))
		}
	}

	err := ioutil.WriteFile(suf+fileName, result, 0x777)

	if err != nil {
		fmt.Printf("%v \n", err)
	}

}

func (d *Downloader) TerminateThread(index int) {
	wg.Done()
}

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
