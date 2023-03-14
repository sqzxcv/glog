package main

import (
    "fmt"
    "github.com/sqzxcv/glog"
    //_ "github.com/sqzxcv/glog/log"
)

func main() {
    glog.SetConsole(true)
    dd := fmt.Sprintf("test22:%s", "aa")

    glog.Error("日志目录:", dd)
    glog.ErrorF("test:%d", 23)
}
