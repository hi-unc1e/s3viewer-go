package main

import (
	"flag"
	"fmt"
	"github.com/hi-unc1e/s3viewer-go/s3viewer"
	"github.com/hi-unc1e/s3viewer-go/web"
	"log"
	"os"
)

func main() {
	// 定义命令行参数
	url := flag.String("u", "http://", "s3 URL, such as http://bucket.s3.amazonaws.com/")
	output := flag.String("o", "", "output file name")
	maxPage := flag.Int("p", 1, "max page")
	webFlag := flag.Bool("web", false, "preview via local_web, such as http://127.0.0.1:30028/static/index.html")
	flag.Parse()
	// 2nd param
	isUseFileOutput := *output != ""
	isRecursively := *maxPage > 0

	// 检查是否提供了所有必需的参数
	if len(os.Args) < 2 {
		fmt.Println("Usage: s3viewer -u s3_url [-o output_file] [-p max_page]")
		return
	}

	if *url == "" {
		log.Fatalf("s3 URL is required")
	}

	// 从远程 URL 加载内容
	var result = new(s3viewer.ListBucketResult)

	// 初始化 URL
	result.Url = *url
	var err error

	if isRecursively {
		result, err = s3viewer.LoadRemoteHTTPRecursive(*url, *maxPage)

		if err != nil {
			log.Fatalf("Failed to load remote URL: %v", err)
		}
	} else {
		result, err = s3viewer.LoadRemoteHTTP(*url)
		if err != nil {
			log.Fatalf("Failed to load remote URL: %v", err)
		}
	}

	if isUseFileOutput {
		// 保存结果到文件
		if err := s3viewer.SaveResultToCSVFile(result, *output); err != nil {
			log.Fatalf("Failed to save result to CSV file: %v", err)
		}
		log.Printf("Saved into %v", *output)
	} else {
		// 打印结果到终端
		if err := s3viewer.PrintResult(result); err != nil {
			log.Fatalf("Failed to print result: %v", err)
		}
	}

	if *webFlag {
		imageUrls := make([]string, 0)
		for _, file := range result.Files {
			imageUrls = append(imageUrls, file.Link)
		}
		web.ServeHttp(imageUrls)
	}
}
