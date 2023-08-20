package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	midOutPut("mr-2-4")
}

func midOutPut(
	fileName string,
) {
	tmpFile, err := ioutil.TempFile(".", fileName)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	//TODO：获取文件大小
	stat, _ := tmpFile.Stat()
	size := stat.Size()
	fmt.Printf("size:%d", size)

}
