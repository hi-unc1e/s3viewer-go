package main

import (
	"flag"
	"fmt"
	"github.com/hi-unc1e/s3viewer-go/s3viewer"
	"log"
	"os"
)

func main() {
	// 定义命令行参数
	url := flag.String("u", "http://", "s3 URL, such as http://bucket.s3.amazonaws.com/")
	output := flag.String("o", "file_list.csv", "output file name")
	flag.Parse()

	// 检查是否提供了所有必需的参数
	if len(os.Args) < 2 {
		fmt.Println("Usage: s3viewer [-u s3_url] [-o output_file]")
		return
	}

	if *url == "" {
		log.Fatalf("s3 URL is required")
	}
	// 从远程 URL 加载内容
	result, err := s3viewer.LoadRemoteHTTP(*url)
	if err != nil {
		log.Fatalf("Failed to load remote URL: %v", err)
	}
	if err := s3viewer.ResultToCSVFile(result, *output); err != nil {
		log.Fatalf("Failed to save result to CSV file: %v", err)
	}
	log.Printf("Saved into %v", *output)
}
