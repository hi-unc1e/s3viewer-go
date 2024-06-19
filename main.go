package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

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
	_            string `xml:"ETag"`
	Size         int    `xml:"Size"`
	_            Owner  `xml:"Owner"`
	_            string `xml:"StorageClass"`
}

type Owner struct {
	DisplayName string `xml:"DisplayName"`
	ID          string `xml:"ID"`
}

func find_s3_xml_string(xml_content string) ([]byte, error) {
	// 定义正则表达式来匹配 <ListBucketResult> 标签及其内容
	re := regexp.MustCompile(`(?s)<ListBucketResult.*?</ListBucketResult>`)

	// 查找匹配的内容
	matches := re.FindAllString(xml_content, -1)

	if len(matches) == 0 {
		match := ""
		return []byte(match), fmt.Errorf("No matches found")
	} else {
		fmt.Println(matches)
		match := matches[0]
		return []byte(match), nil
	}

}

func main() {
	// 读取 XML 文件内容
	file_text, err := ioutil.ReadFile("/Users/dpdu/Desktop/opt/s3view_dev/s3viewer-go/test/h2-html.xml")
	if err != nil {
		log.Fatalf("Failed to read XML file: %v", err)
	}
	xmlContent, err := find_s3_xml_string(string(file_text))
	if err != nil {
		log.Fatalf("Failed to parse XML file: %v", err)
	}

	// 解析 XML 内容为对象
	var result ListBucketResult
	err = xml.Unmarshal(xmlContent, &result)
	if err != nil {
		log.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// 打印解析后的对象
	for _, file := range result.Files {
		fmt.Printf("%v\n", file)
	}

}
