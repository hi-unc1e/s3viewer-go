package s3viewer

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/tabwriter"
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

func LoadRemoteHTTP(url string) (*ListBucketResult, error) {
	// 获取远程 URL 的内容
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch remote URL: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch remote URL: %s", response.Status)
	}

	// 读取响应体
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %w", err)
	}

	// 提取 <ListBucketResult> 标签及其内容
	// 假设 findS3XMLString 和 sanitizeXMLContent 函数已经更新为处理字符串输入
	xmlContent, err := FindS3XMLString(string(body))
	if err != nil {
		return nil, fmt.Errorf("Failed to find S3 XML string: %w", err)
	}

	xmlContent = SanitizeXMLContent(xmlContent)

	// 解析 XML 内容为对象
	result, err := ParseXMLToListBucketResult(xmlContent)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal XML: %w", err)
	}
	return result, err
}

func LoadFile(path string) (*ListBucketResult, error) {
	// 读取 XML 文件内容
	//fileText, err := ioutil.ReadFile("/Users/dpdu/Desktop/opt/s3view_dev/s3viewer-go/test/h2-html.xml")
	fileText, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read XML file: %v", err)
	}

	// 提取 <ListBucketResult> 标签及其内容
	xmlContent, err := FindS3XMLString(string(fileText))
	if err != nil {
		log.Fatalf("Failed to find S3 XML string: %v", err)
	}

	// 预处理 XML 内容
	xmlContent = SanitizeXMLContent(xmlContent)

	// 解析 XML 内容为对象
	result, err := ParseXMLToListBucketResult(xmlContent)
	if err != nil {
		log.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// 打印解析后的对象
	for _, file := range result.Files {
		fmt.Printf("%v\n", file)
	}
	return result, nil
}

func PrintResult(result *ListBucketResult) error {
	// 如果 filePath 不为空，则将输出写入文件
	var output *os.File

	// 输出到标准输出
	output = os.Stdout

	// 创建一个新的 tabwriter
	writer := tabwriter.NewWriter(output, 0, 0, 2, ' ', 0)

	// 打印表头
	fmt.Fprintln(writer, "Key\tSize\tLastModifiedDate")

	// 遍历文件并打印每一行的内容
	for _, file := range result.Files {
		// Todo: fileSize 可以统一换算成合适的单位
		// Todo: Date 可以排个倒序，最新的在最前面
		fmt.Fprintf(writer, "%s\t%s\t%s\n", file.Key, fmt.Sprint(file.Size), file.LastModified)
	}

	// 刷新和清理 tabwriter
	writer.Flush()

	return nil
}

// 将 ListBucketResult 对象转换为 CSV 格式，并保存到指定的文件中
func SaveResultToCSVFile(result *ListBucketResult, filePath string) error {
	// 创建输出文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Failed to create output file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 创建 CSV 写入器
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入 CSV 头部
	headers := []string{"Key", "Size", "LastModified"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("Failed to write CSV headers: %w", err)
	}

	// 写入文件条目
	for _, entry := range result.Files {
		record := []string{entry.Key, fmt.Sprintf("%d", entry.Size), entry.LastModified}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("Failed to write CSV record: %w", err)
		}
	}

	return nil
}

// 查找并提取 <ListBucketResult> 标签及其内容
func FindS3XMLString(xmlContent string) ([]byte, error) {
	re := regexp.MustCompile(`(?s)<ListBucketResult.*?</ListBucketResult>`)
	matches := re.FindAllString(xmlContent, -1)

	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches found")
	}

	return []byte(matches[0]), nil
}

// 预处理 XML 内容，替换无效字符实体
func SanitizeXMLContent(xmlContent []byte) []byte {
	re := regexp.MustCompile(`&[^;]+`)
	return re.ReplaceAllFunc(xmlContent, func(b []byte) []byte {
		if bytes.HasPrefix(b, []byte("&")) && !bytes.Contains(b, []byte(";")) {
			return bytes.Replace(b, []byte("&"), []byte("&amp;"), 1)
		}
		return b
	})
}

// 解析 XML 内容为 ListBucketResult 结构体
func ParseXMLToListBucketResult(xmlContent []byte) (*ListBucketResult, error) {
	var result ListBucketResult
	err := xml.Unmarshal(xmlContent, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
