package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

// 定义结构体以匹配 XML 内容
type ListBucketResult struct {
	XMLName     xml.Name `xml:"ListBucketResult"`
	Name        string   `xml:"Name"`
	Prefix      string   `xml:"Prefix"`
	Marker      string   `xml:"Marker"`
	MaxKeys     int      `xml:"MaxKeys"`
	IsTruncated bool     `xml:"IsTruncated"`
	Files       []File   `xml:"Contents"`
}

type File struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	Size         int    `xml:"Size"`
}

// 查找并提取 <ListBucketResult> 标签及其内容
func findS3XMLString(xmlContent string) ([]byte, error) {
	re := regexp.MustCompile(`(?s)<ListBucketResult.*?</ListBucketResult>`)
	matches := re.FindAllString(xmlContent, -1)

	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches found")
	}

	return []byte(matches[0]), nil
}

// 预处理 XML 内容，替换无效字符实体
func sanitizeXMLContent(xmlContent []byte) []byte {
	re := regexp.MustCompile(`&[^;]+`)
	return re.ReplaceAllFunc(xmlContent, func(b []byte) []byte {
		if bytes.HasPrefix(b, []byte("&")) && !bytes.Contains(b, []byte(";")) {
			return bytes.Replace(b, []byte("&"), []byte("&amp;"), 1)
		}
		return b
	})
}

// 解析 XML 内容为 ListBucketResult 结构体
func parseXMLToListBucketResult(xmlContent []byte) (*ListBucketResult, error) {
	var result ListBucketResult
	err := xml.Unmarshal(xmlContent, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func main() {
	// 读取 XML 文件内容
	fileText, err := ioutil.ReadFile("/Users/dpdu/Desktop/opt/s3view_dev/s3viewer-go/test/h2-html.xml")
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
