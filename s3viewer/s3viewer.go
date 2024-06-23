package s3viewer

import (
	"bytes"
	"crypto/tls"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"text/tabwriter"
	"time"
)

// 定义结构体以匹配 XML 内容
type ListBucketResult struct {
	Url         string
	Prefix      string `xml:"Prefix"`
	NextMarker  string `xml:"NextMarker"` // v1_翻页用
	Marker      string `xml:"Marker"`     // v1_翻页用（备选）
	KeyCount    int    `xml:"KeyCount"`   // 当前数量
	MaxKeys     int    `xml:"MaxKeys"`
	IsTruncated bool   `xml:"IsTruncated"`
	/* IsTruncated
	请求中返回的结果是否被截断。
	- true表示本次没有返回全部结果。
	- false表示本次已经返回了全部结果。
	*/
	NextContinuationToken string `xml:"NextContinuationToken"` //翻页用
	Files                 []File `xml:"Contents"`
}

type File struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	Size         int    `xml:"Size"`
}

func HttpGet(url string) (resp *http.Response, err error) {
	// 创建一个自定义的HTTP客户端(30秒超时，忽略 TLS 证书问题)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略SSL证书验证
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 设置建立连接的超时时间
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   45 * time.Second, // 设置请求的总超时时间
	}
	// 使用自定义的客户端发起GET请求
	log.Printf("Http Get [%v]\n", url)
	response, err := client.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return nil, err
	} else {
		return response, nil
	}
}

func tryGetNextPageURL(currentUrl string, result ListBucketResult) (string, error) {
	/*
		* 方法1
			GET /?marker=?
			拿到第一页的XML中的NextMarker*，构造 URL，继续请求下一页，这样周而复始。
			注：*， 这个 NextMarker* 可能来自NextMarker字段，也可能是最后一个元素的 key
			参考资料：
			- https://help.aliyun.com/zh/oss/developer-reference/listobjects#:~:text=%E5%BC%80%E5%A7%8B%E8%BF%94%E5%9B%9EObject%E3%80%82-,marker%E7%94%A8%E6%9D%A5%E5%AE%9E%E7%8E%B0%E5%88%86%E9%A1%B5,-%E6%98%BE%E7%A4%BA%E6%95%88%E6%9E%9C%EF%BC%8C%E5%8F%82%E6%95%B0
			- 设定从marker之后按字母排序开始返回Object。marker用来实现分页显示效果，参数的长度必须小于1024字节。做条件查询时，即使marker在列表中不存在，也会从符合marker字母排序的下一个开始打印。

		* 方法2
			GET /?list-type=2&continuation-token=?
			拿到XML 中NextContinuationToken的值，构造 URL，继续请求下一页，这样周而复始。
			参考资料：
			- https://help.aliyun.com/zh/oss/developer-reference/listobjectsv2#:~:text=%E9%BB%98%E8%AE%A4%E5%80%BC%EF%BC%9A%E6%97%A0-,continuation%2Dtoken,-%E5%AD%97%E7%AC%A6%E4%B8%B2
			- 指定List操作需要从此token开始。您可从ListObjectsV2（GetBucketV2）结果中的NextContinuationToken获取此token。
	*/
	var nextUrl = currentUrl
	var err error = nil
	// 创建一个新的 URL 查询参数结构体
	query := url.Values{}

	u, error := url.Parse(currentUrl)
	if error != nil {
		return currentUrl, fmt.Errorf("invalid URL: %w", err)
	}

	//fmt.Printf("result-> %+v \n", result)
	// try v2
	if result.NextContinuationToken != "" {
		// 	("/?list-type=2&continuation-token=%s", result.NextContinuationToken)

		query.Set("list-type", "2")
		query.Set("continuation-token", result.NextContinuationToken)
		// 更新 URL 的查询参数
		u.RawQuery = query.Encode()
		nextUrl = u.String()

	}

	// try v1
	// 	("/?marker=%s", NextMarker)
	if result.NextMarker != "" {
		NextMarker := result.Files[len(result.Files)-1].Key
		query.Set("marker", NextMarker)
		nextUrl = u.String()
	} else {
		// try last one, or die
		lastItemMarker := result.Files[len(result.Files)-1].Key
		if lastItemMarker != "" {
			query.Set("marker", lastItemMarker)
			nextUrl = u.String()
		} else {
			err = fmt.Errorf("[!]无法翻页，应该是不支持翻页")
		}
	}

	// 更新 URL 的查询参数
	u.RawQuery = query.Encode()
	nextUrl = u.String()

	log.Printf("尝试请求下一页: %v", nextUrl)
	return nextUrl, err
}

func LoadRemoteHTTPRecursive(url string, maxPage int) (*ListBucketResult, error) {
	// e.g.: http://s3.example.com/
	// 如果不支持翻页，就打印warning，退化到LoadRemoteHTTP
	if maxPage <= 0 {
		maxPage = 1
	}
	var acutalPage int

	var allResults ListBucketResult

	for page := 0; page < maxPage; page++ {
		acutalPage = page + 1
		response, err := HttpGet(url)
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch remote URL: %w", err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Failed to fetch remote URL: %s", response.Status)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed to read response body: %w", err)
		}

		xmlContent, err := FindS3XMLString(string(body))
		if err != nil {
			return nil, fmt.Errorf("Failed to find S3 XML string: %w", err)
		}

		xmlContent = SanitizeXMLContent(xmlContent)

		result, err := ParseXMLToListBucketResult(xmlContent)
		log.Printf("第 %v 页结果条数: %v", acutalPage, len(result.Files))

		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal XML: %w", err)
		}

		allResults.Files = append(allResults.Files, result.Files...)

		// 判断是否有必要翻页
		if !result.IsTruncated {
			log.Printf("不必翻页，本页已经返回了全部结果（%v）", len(result.Files))
			break
		}

		url, err = tryGetNextPageURL(url, *result)
		if err != nil {
			log.Printf("翻页失败，错误: %v", err)
			break
		}
	}
	log.Printf("[+]结果总条数: [%v], 已拉取页数: [%v]", len(allResults.Files), acutalPage)
	return &allResults, nil
}

func LoadRemoteHTTP(url string) (*ListBucketResult, error) {
	// 获取远程 URL 的内容
	// e.g.: http://s3.example.com
	response, err := HttpGet(url)
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
