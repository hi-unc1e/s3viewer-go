package s3viewer

import (
	"fmt"
	"io/ioutil"
	"log"
)

func LoadFile(path string) {
	// 读取 XML 文件内容
	//fileText, err := ioutil.ReadFile("/Users/dpdu/Desktop/opt/s3view_dev/s3viewer-go/test/h2-html.xml")
	fileText, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read XML file: %v", err)
	}

	// 提取 <ListBucketResult> 标签及其内容
	xmlContent, err := findS3XMLString(string(fileText))
	if err != nil {
		log.Fatalf("Failed to find S3 XML string: %v", err)
	}

	// 预处理 XML 内容
	xmlContent = sanitizeXMLContent(xmlContent)

	// 解析 XML 内容为对象
	result, err := parseXMLToListBucketResult(xmlContent)
	if err != nil {
		log.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// 打印解析后的对象
	for _, file := range result.Files {
		fmt.Printf("%v\n", file)
	}
}
