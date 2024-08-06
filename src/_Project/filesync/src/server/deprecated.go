package server

import (

	// "filesync/src/server"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// deprecated
func HandleHttpMaster(w http.ResponseWriter, r *http.Request) {
	// master / slave
	// read: Download

	// master
	// write: Upload / Delete
	// if r.Method == "GET" {
	// 	// HandleDownload(w, r)
	// } else if r.Method == "PUT" {
	// 	// HandleUpload(w, r)
	// } else if r.Method == "DELETE" {
	// 	// HandleDelete(w, r)
	// } else {
	// 	// TODO error
	// }
}

// deprecated
func HandleHttpSlave(w http.ResponseWriter, r *http.Request) {
	// master / slave
	// read: Download
	// if r.Method != "GET" {
	// TODO error
	// }
	// TODO
	// HandleDownload(w, r)

}

// http请求带参数
func LoginServe() {
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9988", nil)
	if err != nil {
		fmt.Println(err)
	}
}

// ***********************************
// http://localhost:9988/login?user=admin&pwd=1234
func login(w http.ResponseWriter, r *http.Request) {
	// 判断参数是否是Get请求，并且参数解析正常
	// if r.Method == "GET" && r.ParseForm() == nil {
	if r.Method == "GET" {
		// 接收参数
		userName := r.FormValue("user") // admin
		fmt.Printf("userName: %s \n", userName)
		passWord := r.FormValue("pwd") // 1234
		fmt.Printf("passWord: %s \n", passWord)
		if userName == "" || passWord == "" {
			w.Write([]byte("用户名或密码不能为空"))
		}
		if userName == "admin" && passWord == "1234" {
			w.Write([]byte("登录成功！"))
		} else {
			w.Write([]byte("用户名或密码错误！"))
		}
	}
}

// http文件操作
const PORT string = "8085"
const IPADDR string = "服务器的外网IP"
const DOWNLOADURL string = "wget http://" + IPADDR + ":" + PORT + "/"
const UPLOADFILEDIR string = "./upload/"
const FILENAMEERR string = "Get file name error!"
const PARAMNUMCHECK int = 3

func GetFileNameFromURL(URL string) (string, error) {
	params := strings.Split(URL, "/")
	fileName := params[len(params)-1]
	fmt.Printf("getFileNameFromURL, len:%d slice=%v\n", len(params), params)
	var err error
	if len(params) > PARAMNUMCHECK {
		err = errors.New(FILENAMEERR)
	}
	return fileName, err
}

func putMethod(w http.ResponseWriter, r *http.Request) {
	fileName, fileNameErr := GetFileNameFromURL(r.URL.Path)
	if fileNameErr != fmt.Errorf("") {
		log.Fatal(fileNameErr.Error())
	}
	fmt.Printf("getHandle URL:%s, filename:%s\n", r.URL.Path, fileName)
	// ***********************************
	// http://localhost:9988/upload/fileuuuu
	// getFileNameFromURL, len:3 slice=[null upload fileuuuu]
	// getHandle URL:/upload/fileuuuu, filename:fileuuuu
	targetFile, err := os.Create(UPLOADFILEDIR + fileName)
	if err != nil {
		panic(err)
	}
	defer targetFile.Close()
	file := r.Body
	n, err := io.Copy(targetFile, file)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\nGet object way: "+DOWNLOADURL+"%s\n", n, fileName)))
}

// upload object
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("method:" + r.Method + "\n")
	if r.Method == "PUT" {
		putMethod(w, r)
	}
}

func getHandle(w http.ResponseWriter, r *http.Request) {
	fileName, err := GetFileNameFromURL(r.URL.Path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("getHandle URL:%s, filename:%s\n", r.URL.Path, fileName)
	objectPath := UPLOADFILEDIR + fileName
	http.ServeFile(w, r, objectPath)
}

// func GetFilelistJSON(flist []string) []byte {
// 	files := GetFilelist(flist)
// 	b, err := json.Marshal(files)
// 	if err != nil {
// 		fmt.Println("GetFilelistJSON error:", err)
// 		return []byte{}
// 	}
// 	fmt.Println(string(b))
// 	return b
// }

func GetFilelist(files []string, path string) *Files {
	flist := Files{
		Files: make([]*FileInfo, 0, 10),
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	// 读取全部文件列表
	for _, fname := range files {
		// 获取文件信息
		info, err := getFileInfo(path + fname)
		if err != nil {
			continue
		}
		flist.Files = append(flist.Files, info)

	}
	return &flist
}

func getFileInfo(fname string) (res *FileInfo, err error) {
	info, err := os.Stat("./src/" + fname)
	if err != nil {
		fmt.Printf("getFileInfo error, filename: %v, error: %v\n", fname, err)
		return nil, err
	}
	// linuxFileAttr := info.Sys().(*syscall.Stat_t)
	res = &FileInfo{
		FileName: fname,
		Size:     info.Size(),
		CTime:    info.ModTime().UnixMicro(),
		MTime:    info.ModTime().UnixMicro(),
	}
	fmt.Println("getFileInfo: ", *res)
	return res, nil
}

// func getFileInfo(fname string) (res *server.FileInfo, err error) {
// 	info, err := os.Stat("./src/" + fname)
// 	if err != nil {
// 		fmt.Printf("getFileInfo error, filename: %v, error: %v\n", fname, err)
// 		return nil, err
// 	}
// 	linuxFileAttr := info.Sys().(*syscall.Stat_t)
// 	res = &server.FileInfo{
// 		FileName: fname,
// 		Size:     linuxFileAttr.Size,
// 		CTime:    linuxFileAttr.Ctim.Sec,
// 		MTime:    linuxFileAttr.Mtim.Sec,
// 	}
// 	fmt.Println("getFileInfo: ", *res)
// 	return res, nil
// }

// func main() {
// http.HandleFunc("GET /files/{filename}", HandleDownload)
// http.HandleFunc("PUT /files/{filename}", HandleUpload)

// http.ListenAndServe(":9988", nil)
// GetFilelistJSON([]string{"test1", "test2", "test3"})
// server.ServerListTest()
// server.ServerSlaveTest()

// server.IsAliveTestServer()
// server.IsAliveTestClient()

// server.SyncOneTestServer()
// server.SyncOneTestClient()

// 同步功能测试
// server.ServerSyncOneSlaveTest()
// server.ServerSyncOneMasterTest()

// 全部功能测试
// server.ServerSyncOneSlaveFullTest()
// server.ServerSyncOneMasterFullTest()

// 	cmd.Execute()

// }

// func HandleHttpMaster(w http.ResponseWriter, r *http.Request) {
// 	// master / slave
// 	// read: Download

// 	// master
// 	// write: Upload / Delete
// 	if r.Method == "GET" {
// 		// HandleDownload(w, r)
// 	} else if r.Method == "POST" {
// 		// res := HandleUpload(r)
// 		// w.Write(res)
// 	} else if r.Method == "DELETE" {
// 		// HandleDelete(w, r)
// 	}
// }

// const PORT string = "8085"
// const IPADDR string = "服务器的外网IP"
// const DOWNLOADURL string = "wget http://" + IPADDR + ":" + PORT + "/"
// const UPLOADFILEDIR string = "./upload/"
// const FILENAMEERR string = "Get file name error !"
// const PARAMNUMCHECK int = 3

// func getFileNameFromURL(URL string) (string, string) {
// 	params := strings.Split(URL, "/")
// 	fileName := params[len(params)-1]
// 	fmt.Printf("getFileNameFromURL, len:%d slice=%v\n", len(params), params)
// 	var err string
// 	if len(params) <= PARAMNUMCHECK+1 {
// 		err = FILENAMEERR
// 	}
// 	return fileName, err
// }

// func HandleUpload(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r.URL.RawQuery)
// 	// fileName, err := getFileNameFromURL(r.URL.Path)
// 	fileName := r.PathValue("filename")
// 	fmt.Println("fileName: ", fileName)
// 	w.Write([]byte(fmt.Sprintf("Upload fileName: %s", fileName)))
// }

// func HandleUploadRAW(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r.URL.RawQuery)
// 	// fileName, err := getFileNameFromURL(r.URL.Path)
// 	fileName := r.PathValue("filename")
// 	fmt.Println("fileName:", fileName)
// 	w.Write([]byte(fmt.Sprintf("fileName: %s", fileName)))
// }

// func HandleDownload(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r.URL.RawQuery)
// 	// fileName, err := getFileNameFromURL(r.URL.Path)
// 	fileName := r.PathValue("filename")
// 	fmt.Println("fileName: ", fileName)
// 	w.Write(TestGetFileList())
// }

// func TestGetFileList() []byte {
// 	type File struct {
// 		Name string `json:"name"`
// 		Size int64  `json:"size"`
// 	}
// 	type Files struct {
// 		Files []*File `json:"files"`
// 	}

// 	files := Files{
// 		Files: make([]*File, 0, 2),
// 	}
// 	files.Files = append(files.Files, &File{
// 		Name: "test1",
// 		Size: 1024,
// 	})
// 	files.Files = append(files.Files, &File{
// 		Name: "test2",
// 		Size: 2048,
// 	})
// 	fmt.Println("files RAW:", files)

// 	b, _ := json.Marshal(files)
// 	fmt.Println("files JSON:", string(b))
// 	return b
// }

// func (s *Server) DownloadFile(filename string, fileContent io.ReadCloser) error {
// 	filename = s.FilesPath + filename
// 	s.Files.Mut.Lock()
// 	// 修改Files,FilesMap
// 	f, ok := s.Files.FilesMap[filename]
// 	if !ok { // 文件不存在于Files中
// 		s.Files.FilesMap[filename] = NewFile(filename, s.Files)
// 		f = s.Files.FilesMap[filename]
// 	}
// 	f.Status = def.FILE_USING_WRITING
// 	f.Mut.Lock()
// 	defer f.Mut.Unlock()
// 	s.Files.Mut.Unlock()

// 	// 判断文件是否存在和是否为目录
// 	_, err := os.Stat(filename)
// 	if err != nil {
// 		s.Logger.Printf("os.Stat-1 error, filename: %v, error: %v", filename, err)
// 	}
// 	targetFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
// 	if err != nil {
// 		s.Logger.Printf("open file error, filename: %v, error: %v\n", filename, err)
// 		return errors.New(def.OPEN_FILE_ERROR)
// 	}
// 	defer targetFile.Close()
// 	n, err := io.Copy(targetFile, fileContent)
// 	if err != nil {
// 		s.Logger.Printf("copy req body to file error, filename: %v, error: %v\n", filename, err)
// 		return errors.New(def.WRITE_FILE_ERROR)
// 	}
// 	s.Logger.Printf("write %d bytes to file %s\n", n, filename)
// 	f.Status = def.FILE_USING_FREE
// 	// 更新文件info
// 	info, err := os.Stat(filename)
// 	if err != nil {
// 		s.Logger.Printf("os.Stat-2 error, filename: %v, error: %v", filename, err)
// 	}
// 	f.Info = &FileInfo{
// 		FileName: filename,
// 		Size:     info.Size(),
// 		CTime:    info.ModTime().UnixMilli(),
// 		MTime:    info.ModTime().UnixMilli(),
// 	}
// 	return nil
// }

// JSON FileList test
func TestGetFileList() []byte {
	type File struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
	}
	type Files struct {
		Files []*File `json:"files"`
	}

	files := Files{
		Files: make([]*File, 0, 2),
	}
	files.Files = append(files.Files, &File{
		Name: "test1",
		Size: 1024,
	})
	files.Files = append(files.Files, &File{
		Name: "test2",
		Size: 2048,
	})
	fmt.Println("files RAW:", files)

	b, _ := json.Marshal(files)
	fmt.Println("files JSON:", string(b))
	return b
}

// func ServerListTest() {
// 	server, err := NewServer(RoleMaster, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	server.Files.FilesMap["test1"] = NewFile("test1", server.Files)
// 	server.Files.FilesMap["test2"] = NewFile("test2", server.Files)
// 	server.Files.FilesMap["test3"] = NewFile("test3", server.Files)

// 	server.Run()
// }

// func ServerUploadTest() {
// 	server, err := NewServer(RoleMaster, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	server.Run()
// }

// func ServerDownloadTest() {
// 	server, err := NewServer(RoleMaster, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	server.Run()
// }

// func ServerDeleteTest() {
// 	server, err := NewServer(RoleMaster, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	server.Run()
// }

// func ServerSlaveTest() {
// 	server, err := NewServer(RoleSlave, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	server.Run()
// }

// func ServerSyncOneSlaveTest() {
// 	server, err := NewServer(RoleSlave, httpIP, httpPort, RPCIP, RPCPort, SlaveFilespath, encryption)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	server.Run()
// }

// func ServerSyncOneMasterTest() {
// 	server, err := NewServer(RoleMaster, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	go server.Run()

// 	time.Sleep(time.Second)
// 	fmt.Printf("start to sync one upload\n")
// 	filenameWithoutPrefix := "FILENAMEU"
// 	err = server.SyncCli.SyncOneUpload(filenameWithoutPrefix)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	time.Sleep(time.Second * 5)
// 	fmt.Printf("start to sync one delete\n")
// 	err = server.SyncCli.SyncOneDelete(filenameWithoutPrefix)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// func ServerSyncOneSlaveFullTest() {
// 	server, err := NewServer(RoleSlave, httpIP, httpPort, RPCIP, RPCPort, SlaveFilespath, encryption)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	server.RunWithoutHttp()
// }

// func ServerSyncOneMasterFullTest() {
// 	server, err := NewServer(RoleMaster, httpIP, httpPort, RPCIP, RPCPort, MasterFilespath, encryption)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	server.Run()
// }

// func (s *Server) HandleDownload2(w http.ResponseWriter, r *http.Request) {
// 	s.Logger.Printf("[Download] New Query: %v\n", r.URL.RequestURI())
// 	filenameWithoutPrefix := r.PathValue("filename")
// 	// w.Write([]byte(fmt.Sprintf("fileName: %s", filename)))

// 	// 下载文件
// 	// err := s.DownloadFile(filename)

// 	filename := s.FilesPath + filenameWithoutPrefix
// 	s.Files.Mut.RLock()
// 	// 修改Files,FilesMap
// 	f, exist := s.Files.FilesMap[filenameWithoutPrefix]
// 	if !exist { // 文件不存在于Files中
// 		s.Logger.Printf("file not exist in Files, filename: %v\n", filenameWithoutPrefix)
// 		w.WriteHeader(http.StatusNotFound)
// 		s.Files.Mut.RUnlock()
// 		return
// 	}
// 	s.Files.Mut.RUnlock()
// 	f.Mut.RLock()
// 	defer f.Mut.RUnlock()

// 	// 判断文件是否存在和是否为目录
// 	info, err := os.Stat(filename)
// 	if err != nil && !os.IsNotExist(err) {
// 		s.Logger.Printf("download os.Stat error, filename: %v, error: %v", filename, err)
// 	} else if os.IsNotExist(err) {
// 		s.Logger.Printf("download os.Stat error, filename: %v, error: %v", filename, err)
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}
// 	if info.IsDir() {
// 		s.Logger.Printf("file is a dir, filename: %v, error: %v", filename, err)
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	f.Info = &FileInfo{
// 		FileName: filename,
// 		Size:     info.Size(),
// 		CTime:    info.ModTime().UnixMilli(),
// 		MTime:    info.ModTime().UnixMilli(),
// 	}

// 	w.Header().Set("Last-Modified", info.ModTime().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
// 	http.ServeFile(w, r, filename)
// 	s.Logger.Println("download success, filename: ", filename)
// }

func SyncAllTestServer() {
	s := NewSyncServer(nil)
	s.Run(TestSyncServerAddr)
}

func SyncAllTestClient() {
	cli := SyncClient{
		ClientName: "SyncOneTestClient",
	}
	cli.Connect(TestSyncServerAddr)
	// cli.SyncAll([]string{"test.txt", "test2.txt"}, [][]byte{[]byte("hello world1"), []byte("hello world2")})
}

func SyncOneTestServer() {
	s := NewSyncServer(nil)
	s.Run(TestSyncServerAddr)
}

func SyncOneTestClient() {
	cli := SyncClient{
		ClientName: "SyncOneTestClient",
	}
	TestFilename := "FILENAMEU"
	cli.Connect(TestSyncServerAddr)
	cli.SyncOneUpload(TestFilename)
}

func IsAliveTestServer() {
	s := NewSyncServer(nil)
	s.Run(TestSyncServerAddr)
}

func IsAliveTestClient() {
	cli := SyncClient{
		ClientName: "IsAliveTestClient",
	}
	cli.Connect(TestSyncServerAddr)
	cli.SyncIsAlive()
}
