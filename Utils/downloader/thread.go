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
	progress   float32
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

	count, size := CalculateBytesPerStep(t.end - t.start - 1)

	buffer := make([]byte, 0)

	for i := 0; i < count; i++ {
		req, err := http.NewRequest("GET", t.url, nil)
		if err != nil {
			fmt.Println("fuck")
			fmt.Println(err)
			return
		}
		start := t.start + (i * size)
		end := start + size - 1

		if i == count-1 {
			end = t.end - 1
		}

		range_header := "bytes=" + strconv.Itoa(start) + "-" + strconv.Itoa(end) // Add the data for the Range header of the form "bytes=0-100"
		req.Header.Add("Range", range_header)

		resp, err := client.Do(req)
		if err != nil {
			return
		}

		reader, _ := ioutil.ReadAll(resp.Body)
		buffer = append(buffer, reader...)

		t.progress = float32(i) / float32(count)

		resp.Body.Close()
	}

	ioutil.WriteFile(strconv.Itoa(t.index), buffer, 0x777)

	t.downloader.TerminateThread(t.index)
}

func CalculateBytesPerStep(byteCount int) (numberOfTry int, trySize int) {

	count := byteCount / (1024 * 100)
	return count, 1024 * 100
}
