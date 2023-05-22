package glog

import (
    "fmt"
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

func SetRollingFile(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
    maxFileCount = maxNumber
    maxFileSize = maxSize * int64(_unit)
    RollingFile = true
    dailyRolling = false
    mkdirlog(fileDir)
    logObj = &_FILE{dir: fileDir, filename: fileName, isCover: false, mu: new(sync.RWMutex)}
    logObj.mu.Lock()
    defer logObj.mu.Unlock()
    for i := 1; i <= int(maxNumber); i++ {
        if isExist(fileDir + "/" + fileName + "." + strconv.Itoa(i)) {
            logObj._suffix = i
        } else {
            break
        }
    }
    if !logObj.isMustRename() {
        logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
        logObj.lg = log.New(logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
    } else {
        logObj.rename()
    }
    //logObj.todb = todb
    go fileMonitor()
}

func SetRollingDaily(fileDir, fileName string, todb bool) {
    RollingFile = false
    dailyRolling = true
    t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
    mkdirlog(fileDir)
    logObj = &_FILE{dir: fileDir, filename: fileName, _date: &t, isCover: false, mu: new(sync.RWMutex)}
    logObj.mu.Lock()
    defer logObj.mu.Unlock()

    if !logObj.isMustRename() {
        logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
        logObj.lg = log.New(logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
    } else {
        logObj.rename()
    }
    logObj.lg.Todb = todb
}

func SetLogFile(fileDir string, fileName string, todb bool) {
    SetRollingDaily(fileDir, fileName, todb)
}

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

    if dailyRolling {
        fileCheck()
    }
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
    if dailyRolling {
        fileCheck()
    }
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
    if dailyRolling {
        fileCheck()
    }
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
    if dailyRolling {
        fileCheck()
    }
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
    if dailyRolling {
        fileCheck()
    }
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

    if dailyRolling {
        fileCheck()
    }
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
    if dailyRolling {
        fileCheck()
    }
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
    if dailyRolling {
        fileCheck()
    }
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
    if dailyRolling {
        fileCheck()
    }
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
    if dailyRolling {
        fileCheck()
    }
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
    if dailyRolling {
        t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
        if t.After(*f._date) {
            return true
        }
    } else {
        if maxFileCount > 1 {
            if fileSize(f.dir+"/"+f.filename) >= maxFileSize {
                return true
            }
        }
    }
    return false
}

func (f *_FILE) rename() {
    if dailyRolling {
        fn := f.dir + "/" + f.filename + "." + f._date.Format(DATEFORMAT)
        if !isExist(fn) && f.isMustRename() {
            if f.logfile != nil {
                f.logfile.Close()
            }
            err := os.Rename(f.dir+"/"+f.filename, fn)
            if err != nil {
                f.lg.Println("rename err", err.Error())
            }
            t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
            f._date = &t
            f.logfile, _ = os.Create(f.dir + "/" + f.filename)
            //f.lg = log.New(logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
            f.lg = log.New(logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
        }
    } else {
        f.coverNextOne()
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
    fmt.Println("fileSize", file)
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
    timer := time.NewTicker(1 * time.Second)
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
    }
}
