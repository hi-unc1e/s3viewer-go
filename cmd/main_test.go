package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestMainFunc(t *testing.T) {
	// 重置 flag
	resetFlags()

	// 创建临时目录
	tmpDir, err := ioutil.TempDir("", "s3viewer_cli_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 设置命令行参数
	os.Args = []string{"cmd", "-u", "https://dl.qianxin.com/", "-o", filepath.Join(tmpDir, "dl.qianxin.com.csv")}

	// 运行 main 函数
	main()

	// 检查文件是否生成
	outputFile := filepath.Join(tmpDir, "dl.qianxin.com.csv")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Expected output file %s, but it does not exist", outputFile)
	}
	text, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	log.Println(string(text[:100]))
}

func CaptureOutput(f func()) string {
	//捕获标准输出
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	// Copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// Run the function
	f()

	// Back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	return <-outC
}

func TestMissingOutput(t *testing.T) {
	// 重置 flag
	resetFlags()

	output := CaptureOutput(func() {
		os.Args = []string{"cmd", "-u", "https://dl.qianxin.com/"}
		main()
	})
	expected := "Key"
	if !strings.Contains(output, expected) {
		t.Fatalf("Expected %q to contain %q", output, expected)
	}
}
