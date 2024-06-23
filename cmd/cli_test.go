package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestMainFunc(t *testing.T) {
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

func TestMissingOutput(t *testing.T) {
	outputFile := "file_list.csv"
	// 设置命令行参数
	os.Args = []string{"cmd", "-u", "https://dl.qianxin.com/"}

	// 运行 main 函数
	main()

	text, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	assert.NotEmpty(t, text, fmt.Sprintf("File %v is empty", outputFile))
}
