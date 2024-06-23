package main

import (
	"flag"
	"fmt"
	"github.com/rendon/testcli"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

func TestMissingOutput(t *testing.T) {
	// 重置 flag
	resetFlags()

	// 捕获标准输出
	var ttyName string
	stdOut := make([]byte, 1024)

	if runtime.GOOS == "windows" {
		//log.Println("*** Using `con`")
		ttyName = "con"
	} else {
		//fmt.Println("*** Using `/dev/tty`")
		ttyName = "/dev/tty"
	}

	f, err := os.OpenFile(ttyName, os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 设置命令行参数并运行 main 函数
	os.Args = []string{"cmd", "-u", "https://dl.qianxin.com/"}
	main()

	// 读取标准输出
	_, err = f.Read(stdOut)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading from file:", err)
	}
	log.Printf("Stdout: %v", stdOut)

	assert.NotEmpty(t, string(stdOut), "StdOutput is empty")
	assert.Contains(t, string(stdOut), "Key", "StdOutput does not contains kw")
}

func TestCLI_V2_MissingOutput(t *testing.T) {
	// Using the struct version, if you want to test multiple commands
	c := testcli.Command("go", "run", "main.go", "-u", "https://dl.qianxin.com/")
	c.Run()
	if !c.Success() {
		t.Fatalf("Expected to succeed, but failed with error: %s", c.Error())
	}

	if !c.StdoutContains("Hello John!") {
		t.Fatalf("Expected %q to contain %q", c.Stdout(), "Key")
	}
}
