package s3viewer

import (
	"bytes"
	"encoding/xml"
	"fmt"
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
