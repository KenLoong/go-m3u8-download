package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var httpClient *http.Client


func init() {
	httpClient = &http.Client{
		Timeout: 3 * time.Minute,
		//Transport: &http.Transport{
		//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//},
	}
}

// https://superuser.com/questions/663687/merge-join-concatenate-hundreds-of-ts-files-into-one-ts-file
func TestDownfile(t *testing.T) {

	urls, err := ParseM3u8("index2.m3u8")
	//urls, err := ParseFile()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = DownloadFiles(urls)
	if err != nil {
		log.Fatal(err)
	}

	tool := NewTool()
	tool.Merge(fmt.Sprintf("%v-merge.ts", time.Now().Unix()), len(urls))

	fmt.Printf("File  downlaod in current working directory")
}

func ParseM3u8(fileName string) ([]string, error) {

	usefulUrls := []string{}
	m3u8Indexs, err := ParseFile(fileName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for _, index := range m3u8Indexs {
		if strings.HasSuffix(index, ".ts") {
			usefulUrls = append(usefulUrls, fmt.Sprintf("https://cdn.bigcloud.click/hls/446301/%v", index))
		}
	}

	return usefulUrls, nil
}

func DownloadFiles(urls []string) error {

	for i, url := range urls {
		//Get the response bytes from the url
		response, err := httpClient.Get(url)
		if err != nil {
			return err
		}

		if response.StatusCode != 200 {
			return errors.New("Received non 200 response code")
		}

		//Create a empty file
		file, err := os.Create(fmt.Sprintf("%v.ts", i))
		if err != nil {
			return err
		}

		//Write the bytes to the file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			return err
		}
		response.Body.Close()
		file.Close()
	}

	return nil

}

func ParseFile(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	defer f.Close()

	return readTxt(f)
}

func readTxt(f *os.File) ([]string, error) {
	reader := bufio.NewReader(f)

	l := make([]string, 0, 64)

	// 按行读取
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		l = append(l, strings.Trim(string(line), " "))
	}

	return l, nil
}
