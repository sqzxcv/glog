package log

import (
	"bytes"
	"strings"
	"testing"
)

// TestOutputAtUsesGivenLocation OutputAt 应当使用参数里的位置, 而不是现场推导调用栈
func TestOutputAtUsesGivenLocation(t *testing.T) {
	var buf bytes.Buffer
	lg := New(&buf, "", Ldate|Ltime|Lshortfile)

	if err := lg.OutputAt(3, "/a/b/order.go", 42, "下单成功"); err != nil {
		t.Fatalf("OutputAt 返回错误: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "order.go:42") {
		t.Errorf("输出 = %q, 期望包含 order.go:42", got)
	}
	if !strings.Contains(got, "下单成功") {
		t.Errorf("输出 = %q, 期望包含 下单成功", got)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("输出 = %q, 期望以换行结尾", got)
	}
}

// TestOutputStillResolvesCaller 重构后 Output 仍应自己推导调用位置
func TestOutputStillResolvesCaller(t *testing.T) {
	var buf bytes.Buffer
	lg := New(&buf, "", Ldate|Ltime|Lshortfile)

	if err := lg.Output(3, 1, "来自测试"); err != nil {
		t.Fatalf("Output 返回错误: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "log_test.go:") {
		t.Errorf("输出 = %q, 期望包含 log_test.go:", got)
	}
}
