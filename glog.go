package glog

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"

	//"llhgo/logger/log"
	"github.com/sqzxcv/glog/log"
	//"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	_VER string = "1.0.2"
)

type LEVEL int32

// LogLevel 日志等级
var LogLevel LEVEL = 1
var maxFileSize int64
var maxFileCount int32
var dailyRolling bool = true
var consoleAppender bool = true
var RollingFile bool = false
var logObj *_FILE
var todb bool = false

//var usock *net.UDPConn

const DATEFORMAT = "2006-01-02"

type UNIT int64

const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

const (
	ALL LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

type _FILE struct {
	dir      string
	filename string
	_suffix  int
	//_todb    bool
	isCover bool
	_date   *time.Time
	mu      *sync.RWMutex
	logfile *os.File
	lg      *log.Logger
}

func SetConsole(isConsole bool) {
	consoleAppender = isConsole
}

func SetLevel(_level LEVEL) {
	LogLevel = _level
	//log.NewSock("127.0.0.1:7777")
}

func SetLevelWithName(_level string) {
	switch _level {
	case "all":
		LogLevel = ALL
	case "debug":
		LogLevel = DEBUG
	case "info":
		LogLevel = INFO
	case "warn":
		LogLevel = WARN
	case "error":
		LogLevel = ERROR
	case "fatal":
		LogLevel = FATAL
	case "off":
		LogLevel = OFF
	}

	//log.NewSock("127.0.0.1:7777")
}

func SetSvrId(svrName string, svrId string) {
	log.Log_svrName = svrName
	log.Log_svrId = svrId
}

//func SetRollingFile(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
//	maxFileCount = maxNumber
//	maxFileSize = maxSize * int64(_unit)
//	RollingFile = true
//	dailyRolling = false
//	mkdirlog(fileDir)
//	logObj = &_FILE{dir: fileDir, filename: fileName, isCover: false, mu: new(sync.RWMutex)}
//	logObj.mu.Lock()
//	defer logObj.mu.Unlock()
//	for i := 1; i <= int(maxNumber); i++ {
//		if isExist(fileDir + "/" + fileName + "." + strconv.Itoa(i)) {
//			logObj._suffix = i
//		} else {
//			break
//		}
//	}
//	if !logObj.isMustRename() {
//		logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
//		logObj.lg = log.New(logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
//	} else {
//		logObj.rename()
//	}
//	//logObj.todb = todb
//	go fileMonitor()
//}

func SetRolling(fileDir, fileName string, todb bool, maxNumber int32, maxSize int64, _unit UNIT) {
	maxFileCount = maxNumber
	maxFileSize = maxSize * int64(_unit)
	RollingFile = true
	dailyRolling = true
	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	mkdirlog(fileDir)
	logObj = &_FILE{dir: fileDir, filename: fileName, _date: &t, isCover: false, mu: new(sync.RWMutex)}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()

	// 从日志文件中过滤所有符合条件的文件 filename.data.num, 然后找出最大的num

	nums, err := getOneDayLogFileNum(fileDir, fileName, t.Format(DATEFORMAT))
	if err != nil {
		logObj.lg.Println("getOneDayLogFileNum err", err.Error())
		os.Exit(1)
	}
	maxNum := 0
	for _, num := range nums {
		if num > maxNum {
			maxNum = num
		}
	}
	logObj._suffix = maxNum

	if !logObj.isMustRename() {
		logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		logObj.lg = log.New(logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logObj.rename()
	}
	logObj.lg.Todb = todb
	go fileMonitor()
}

//func SetLogFile(fileDir string, fileName string, todb bool) {
//	SetRollingDaily(fileDir, fileName, todb)
//}

func Flush() {

}

func mkdirlog(dir string) (e error) {
	_, er := os.Stat(dir)
	b := er == nil || os.IsExist(er)
	if !b {
		if err := os.MkdirAll(dir, 0666); err != nil {
			if os.IsPermission(err) {
				fmt.Println("create dir error:", err.Error())
				e = err
			}
		}
	}
	return
}
func console(level int, calldepth int, s ...interface{}) {
	if consoleAppender {
		log.Std.Output(level, calldepth, fmt.Sprintln(s...))
	}
}
func catchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
	}
}

func Debug(v ...interface{}) {

	//if dailyRolling {
	//	fileCheck()
	//}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if LogLevel <= DEBUG {
		if logObj != nil {
			logObj.lg.Output(4, 2, fmt.Sprintln(v...))
		}
		console(4, 3, v...)
	}
}

func Info(v ...interface{}) {
	//if dailyRolling {
	//	fileCheck()
	//}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	if LogLevel <= INFO {
		if logObj != nil {
			logObj.lg.Output(3, 2, fmt.Sprintln(v...))
		}
		console(3, 3, v...)
	}
}
func Warn(v ...interface{}) {
	//if dailyRolling {
	//	fileCheck()
	//}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if LogLevel <= WARN {
		if logObj != nil {
			logObj.lg.Output(2, 2, fmt.Sprintln(v...))
		}
		console(2, 3, v...)
	}
}

func Error(v ...interface{}) {
	//if dailyRolling {
	//	fileCheck()
	//}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	if LogLevel <= ERROR {
		if logObj != nil {
			logObj.lg.Output(1, 2, fmt.Sprintln(v...))
		}
		console(1, 3, v...)
	}
}

func Fatal(v ...interface{}) {
	//if dailyRolling {
	//	fileCheck()
	//}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	if LogLevel <= FATAL {
		if logObj != nil {
			logObj.lg.Output(0, 2, fmt.Sprintln(v...))
		}
		console(0, 3, v...)
	}
}

func FDebug(format string, v ...interface{}) {

	//if dailyRolling {
	//	fileCheck()
	//}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if LogLevel <= DEBUG {
		if logObj != nil {
			logObj.lg.Output(4, 2, fmt.Sprintf(format, v...))
		}
		console(4, 3, fmt.Sprintf(format, v...))
	}
}

func FInfo(format string, v ...interface{}) {
	//if dailyRolling {
	//	fileCheck()
	//}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	if LogLevel <= INFO {
		if logObj != nil {
			logObj.lg.Output(3, 2, fmt.Sprintf(format, v...))
		}
		console(3, 3, fmt.Sprintf(format, v...))
	}
}
func FWarn(format string, v ...interface{}) {
	//if dailyRolling {
	//	fileCheck()
	//}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if LogLevel <= WARN {
		if logObj != nil {
			logObj.lg.Output(2, 2, fmt.Sprintf(format, v...))
		}
		console(2, 3, fmt.Sprintf(format, v...))
	}
}

func FError(format string, v ...interface{}) {
	//if dailyRolling {
	//	fileCheck()
	//}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	if LogLevel <= ERROR {
		if logObj != nil {
			logObj.lg.Output(1, 2, fmt.Sprintf(format, v...))
		}
		console(1, 3, fmt.Sprintf(format, v...))
	}
}

func FFatal(format string, v ...interface{}) {
	//if dailyRolling {
	//	fileCheck()
	//}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	if LogLevel <= FATAL {
		if logObj != nil {
			logObj.lg.Output(0, 2, fmt.Sprintf(format, v...))
		}
		console(0, 3, fmt.Sprintf(format, v...))
	}
}

func (f *_FILE) isMustRename() bool {

	if f.isNewDate() {
		return true
	}
	if f.isCurrentLogFileTooBig() {
		return true
	}
	return false
}

func (f *_FILE) isNewDate() bool {
	if dailyRolling {
		t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
		if t.After(*f._date) {
			return true
		}
	}
	return false
}

func (f *_FILE) isCurrentLogFileTooBig() bool {
	if maxFileCount > 1 {
		if fileSize(f.dir+"/"+f.filename) >= maxFileSize {
			return true
		}
	}
	return false
}

func (f *_FILE) rename() {
	if f.isNewDate() {

		// 如果存在符合 filename.date.num 格式的文件，则不进行重命名
		if hasPrefixFile(f.dir, f.filename, f._date.Format(DATEFORMAT)) == false {
			if f.logfile != nil {
				f.logfile.Close()
			}
			fn := fmt.Sprintf("%s/%s.%s.%d", f.dir, f.filename, f._date.Format(DATEFORMAT), 1)
			f._suffix = 1
			err := os.Rename(f.dir+"/"+f.filename, fn)
			if err != nil {
				f.lg.Println("rename err", err.Error())
			}
			t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
			f._date = &t
			f.logfile, _ = os.Create(f.dir + "/" + f.filename)
			//f.lg = log.New(logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
			f.lg = log.New(logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
			return
		}
	}
	if f.isCurrentLogFileTooBig() {
		num, err := getOneDayLogFileNum(f.dir, f.filename, f._date.Format(DATEFORMAT))
		if err != nil {
			f.lg.Println("getOneDayLogFileNum err", err.Error())
			return
		}
		maxNum := f._suffix
		for _, n := range num {
			if n > maxNum {
				maxNum = n
			}
		}
		f._suffix = maxNum + 1

		if f.logfile != nil {
			f.logfile.Close()
		}
		fn := fmt.Sprintf("%s/%s.%s.%d", f.dir, f.filename, f._date.Format(DATEFORMAT), f._suffix)

		os.Rename(f.dir+"/"+f.filename, fn)
		f.logfile, _ = os.Create(f.dir + "/" + f.filename)
		t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
		f._date = &t
		f.lg = log.New(logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)

	}
}

func (f *_FILE) nextSuffix() int {
	return int(f._suffix%int(maxFileCount) + 1)
}

func (f *_FILE) coverNextOne() {
	f._suffix = f.nextSuffix()
	if f.logfile != nil {
		f.logfile.Close()
	}
	if isExist(f.dir + "/" + f.filename + "." + strconv.Itoa(int(f._suffix))) {
		os.Remove(f.dir + "/" + f.filename + "." + strconv.Itoa(int(f._suffix)))
	}
	os.Rename(f.dir+"/"+f.filename, f.dir+"/"+f.filename+"."+strconv.Itoa(int(f._suffix)))
	f.logfile, _ = os.Create(f.dir + "/" + f.filename)
	f.lg = log.New(logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func fileSize(file string) int64 {
	//fmt.Println("fileSize", file)
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
		return 0
	}
	return f.Size()
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func fileMonitor() {
	timer := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-timer.C:
			fileCheck()
		}
	}
}

func fileCheck() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	if logObj != nil && logObj.isMustRename() {
		logObj.mu.Lock()
		defer logObj.mu.Unlock()
		logObj.rename()
		removeMoreOldLogFile(logObj.dir, logObj.filename, int(maxFileCount))
	}
}

// getOneDayLogFileNum 获取指定目录下，符合特定文件名前缀和日期的文件，并返回文件名中的 num 数组
func getOneDayLogFileNum(dir, filename, date string) ([]int, error) {
	var nums []int
	// 使用正则表达式匹配文件名，假设 num 为非负整数
	regexPattern := fmt.Sprintf(`^%s\.%s\.(\d+)$`, regexp.QuoteMeta(filename), date)
	regex := regexp.MustCompile(regexPattern)

	// 读取目录下的所有文件
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			// 获取文件名
			baseName := filepath.Base(file.Name())
			// 使用正则表达式匹配文件名
			matches := regex.FindStringSubmatch(baseName)
			if len(matches) == 2 {
				// 提取匹配的 num，并转换为整数
				num, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, err
				}
				nums = append(nums, num)
			}
		}
	}
	return nums, nil
}

// hasPrefixFile 检查指定目录下是否存在符合 filename.date.num 格式的文件
func hasPrefixFile(dir, filename, date string) bool {
	// 构造正则表达式以匹配 filename.date.num 格式的文件
	regexPattern := fmt.Sprintf(`^%s\.%s\.\d+$`, regexp.QuoteMeta(filename), date)
	regex := regexp.MustCompile(regexPattern)

	// 读取目录下的所有文件
	files, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	// 遍历目录下的文件，检查是否有匹配的文件名
	for _, file := range files {
		if !file.IsDir() {
			// 获取文件名
			baseName := filepath.Base(file.Name())
			// 检查文件名是否符合正则表达式
			if regex.MatchString(baseName) {
				return true
			}
		}
	}
	return false
}

// removeMoreOldLogFile 删除 dir 下以 filename 开头的文件，保留最新的 fileCount 个文件
func removeMoreOldLogFile(dir, filename string, fileCount int) {
	// 构造正则表达式以匹配 filename.date.num 格式的文件
	regexPattern := fmt.Sprintf(`^%s\.(\d{4}-\d{2}-\d{2})(?:\.(\d+))?$`, regexp.QuoteMeta(filename))
	regex := regexp.MustCompile(regexPattern)

	// FileInfo 包含文件名和提取的日期、num
	type FileInfo struct {
		name string
		date string
		num  int
	}

	var files []FileInfo

	// 读取目录下的所有文件
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// 遍历目录下的文件，匹配符合条件的文件
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			baseName := filepath.Base(entry.Name())
			matches := regex.FindStringSubmatch(baseName)
			if len(matches) > 1 {
				num := 0
				if len(matches) > 2 && matches[2] != "" {
					num, err = strconv.Atoi(matches[2])
					if err != nil {
						fmt.Println("Error converting num to integer:", err)
						continue
					}
				}
				files = append(files, FileInfo{name: baseName, date: matches[1], num: num})
			}
		}
	}

	// 如果文件数量小于或等于需要保留的数量，直接返回
	if len(files) <= fileCount {
		return
	}

	// 排序文件：先按日期升序，再按 num 升序
	sort.Slice(files, func(i, j int) bool {
		if files[i].date == files[j].date {
			return files[i].num < files[j].num
		}
		return files[i].date < files[j].date
	})

	// 计算需要删除的文件数量
	numToDelete := len(files) - fileCount

	// 删除旧文件
	for i := 0; i < numToDelete; i++ {
		filePath := filepath.Join(dir, files[i].name)
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Error deleting file:", filePath, err)
		} else {
			fmt.Println("Deleted file:", filePath)
		}
	}
}
