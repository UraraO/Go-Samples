package heartcheck

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type tickData struct {
	Reqtime  time.Time
	Recvtime time.Time
	Received bool
}

type tickDataBlock struct {
	dataMut sync.Mutex
	Datas   []tickData
	Url     string
}

func InittickDataBlock(url string) *tickDataBlock {
	return &tickDataBlock{
		Datas: make([]tickData, 0, 10),
		Url:   url,
	}
}

func (bs *tickDataBlock) Add(reqtime, recvtime time.Time, received bool) {
	bs.dataMut.Lock()
	defer bs.dataMut.Unlock()
	bs.Datas = append(bs.Datas, tickData{
		Reqtime:  reqtime,
		Recvtime: recvtime,
		Received: received,
	})
}

func (bs *tickDataBlock) PersistToFile(filename string) error {
	bs.dataMut.Lock()
	defer bs.dataMut.Unlock()
	if len(bs.Datas) == 0 {
		return nil
	}
	if filename == "" {
		filename = bs.Url
	}
	// 判断文件是否存在
	var f *os.File
	// _, err := os.Stat(filename)
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		// panic("文件不存在")
	// 		fmt.Println("PersistToFile filename not exist")
	// 		f, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	// 	} else {
	// 		fmt.Println("PersistToFile os.Stat error:", err.Error())
	// 		return err
	// 	}
	// } else {
	// 	f, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	// }
	filename = strings.TrimPrefix(filename, "http")
	filename = strings.TrimPrefix(filename, "s")
	filename = strings.TrimPrefix(filename, "://")
	filename += ".log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println("PersistToFile os.OpenFile error:", err.Error())
		return err
	}

	for i := range bs.Datas {
		// json序列化
		bytes, err := json.Marshal(bs.Datas[i])
		if err != nil {
			fmt.Println("PersistToFile json.Marshal error:", err.Error())
			return err
		}
		// 写入文件
		_, err = f.Write(bytes)
		if err != nil {
			fmt.Println("PersistToFile file.Write error:", err.Error())
			return err
		}
		_, err = f.Write([]byte("\n"))
		if err != nil {
			fmt.Println("PersistToFile file.Write error:", err.Error())
			return err
		}
	}
	bs.Datas = make([]tickData, 0, 10)
	return nil
}
