package glog

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sqzxcv/glog/log"
)

// captureConsole 把控制台输出重定向到 buf, 返回还原函数
func captureConsole(buf *bytes.Buffer) func() {
	old := log.Std
	log.Std = log.New(buf, "", log.Ldate|log.Ltime|log.Lshortfile)
	oldLevel := LogLevel
	SetConsole(true)
	SetLevel(ALL)
	return func() {
		log.Std = old
		LogLevel = oldLevel
	}
}

// TestOutputAtUsesGivenLocation 包级 OutputAt 应使用参数里的位置
func TestOutputAtUsesGivenLocation(t *testing.T) {
	var buf bytes.Buffer
	defer captureConsole(&buf)()

	OutputAt(INFO, "/a/b/order.go", 42, "下单成功")

	got := buf.String()
	if !strings.Contains(got, "order.go:42") {
		t.Errorf("输出 = %q, 期望包含 order.go:42", got)
	}
	if !strings.Contains(got, "下单成功") {
		t.Errorf("输出 = %q, 期望包含 下单成功", got)
	}
	if !strings.Contains(got, "[INFO]") {
		t.Errorf("输出 = %q, 期望包含 [INFO] 级别标志", got)
	}
}

// TestOutputAtRespectsLogLevel 低于 LogLevel 的级别不应输出
func TestOutputAtRespectsLogLevel(t *testing.T) {
	var buf bytes.Buffer
	defer captureConsole(&buf)()
	SetLevel(ERROR)

	OutputAt(INFO, "/a/b/order.go", 42, "不该出现")

	if buf.Len() != 0 {
		t.Errorf("输出 = %q, 期望为空", buf.String())
	}
}

// TestInfoStillReportsCallerSite 重构后 Info 仍应报告调用方所在文件
func TestInfoStillReportsCallerSite(t *testing.T) {
	var buf bytes.Buffer
	defer captureConsole(&buf)()

	Info("来自测试")

	got := buf.String()
	if !strings.Contains(got, "glog_test.go:") {
		t.Errorf("输出 = %q, 期望包含 glog_test.go:", got)
	}
}

// TestFInfoFormats FInfo 应按格式串输出
func TestFInfoFormats(t *testing.T) {
	var buf bytes.Buffer
	defer captureConsole(&buf)()

	FInfo("订单 %s 金额 %d", "O123", 99)

	got := buf.String()
	if !strings.Contains(got, "订单 O123 金额 99") {
		t.Errorf("输出 = %q, 期望包含格式化后的内容", got)
	}
}

// TestFatalDoesNotExit Fatal 必须不退出进程
func TestFatalDoesNotExit(t *testing.T) {
	var buf bytes.Buffer
	defer captureConsole(&buf)()

	Fatal("致命但不退出")

	if !strings.Contains(buf.String(), "致命但不退出") {
		t.Errorf("输出 = %q, 期望包含消息", buf.String())
	}
}

// TestAllLevelsReportCallerSite 五个级别都应报告调用方位置
func TestAllLevelsReportCallerSite(t *testing.T) {
	cases := []struct {
		name string
		fn   func(...interface{})
	}{
		{"Debug", Debug},
		{"Info", Info},
		{"Warn", Warn},
		{"Error", Error},
		{"Fatal", Fatal},
	}
	for _, c := range cases {
		var buf bytes.Buffer
		restore := captureConsole(&buf)
		c.fn("测试消息")
		restore()

		if !strings.Contains(buf.String(), "glog_test.go:") {
			t.Errorf("%s 输出 = %q, 期望包含 glog_test.go:", c.name, buf.String())
		}
	}
}

// TestAllFormatLevelsReportCallerSite 五个格式化函数都应报告调用方位置
func TestAllFormatLevelsReportCallerSite(t *testing.T) {
	cases := []struct {
		name string
		fn   func(string, ...interface{})
	}{
		{"FDebug", FDebug},
		{"FInfo", FInfo},
		{"FWarn", FWarn},
		{"FError", FError},
		{"FFatal", FFatal},
	}
	for _, c := range cases {
		var buf bytes.Buffer
		restore := captureConsole(&buf)
		c.fn("测试 %d", 1)
		restore()

		if !strings.Contains(buf.String(), "glog_test.go:") {
			t.Errorf("%s 输出 = %q, 期望包含 glog_test.go:", c.name, buf.String())
		}
	}
}
