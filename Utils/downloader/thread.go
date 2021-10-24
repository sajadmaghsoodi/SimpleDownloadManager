package downloader

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Thread struct {
	url        string
	start      int
	end        int
	index      int
	downloader *Downloader
}

func NewThread(url string, startByte int, endByte int, index int, downloader *Downloader) *Thread {

	t := Thread{}
	t.url = url
	t.start = startByte
	t.end = endByte
	t.index = index
	t.downloader = downloader

	return &t
}

func (t *Thread) Start() {

	client := &http.Client{}

	req, err := http.NewRequest("GET", t.url, nil)
	if err != nil {
		fmt.Println("fuck")
		fmt.Println(err)
		return
	}

	range_header := "bytes=" + strconv.Itoa(t.start) + "-" + strconv.Itoa(t.end-1) // Add the data for the Range header of the form "bytes=0-100"
	req.Header.Add("Range", range_header)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	reader, _ := ioutil.ReadAll(resp.Body)

	ioutil.WriteFile(strconv.Itoa(t.index), []byte(string(reader)), 0x777)

	t.downloader.TerminateThread(t.index)
}

func (t *Thread) GetProgress() float32 {

	file, _ := ioutil.ReadFile(strconv.Itoa(t.index))
	currentSize := len(file)
	fullSize := t.end - t.start
	res := float32(currentSize) / float32(fullSize)
	return res
}
