package main

import (
    "fmt"
    "github.com/sqzxcv/glog"
    //_ "github.com/sqzxcv/glog/log"
)

func main() {
    logsDir := "./logs/" //GetCurrentDirectory()
    glog.SetConsole(true)
    glog.SetLevel(0)
    glog.Info("日志目录:", logsDir)
    glog.SetRollingDaily(logsDir, "ss_backend.log", false)
    dd := fmt.Sprintf("test22:%s", "aa")

    glog.Error("日志目录:", dd)
    glog.FError("1ddd3etest:%d", 23)
    glog.Debug("test", 232, "aaa")
}
