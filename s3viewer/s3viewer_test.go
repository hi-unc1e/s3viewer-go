package s3viewer

import (
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const sampleXML = `
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <Name>example-bucket</Name>
    <Prefix/>
    <Marker/>
    <MaxKeys>1000</MaxKeys>
    <IsTruncated>false</IsTruncated>
    <Contents>
        <Key>test-file1.txt</Key>
        <LastModified>2023-05-22T08:49:07.000Z</LastModified>
        <ETag>"f19cd76cd7fac68d15f0c40a063519c9"</ETag>
        <Size>1234</Size>
        <Owner>
            <DisplayName>owner</DisplayName>
            <ID>owner-id</ID>
        </Owner>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
</ListBucketResult>`

const specialXML = `
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <Name>example-bucket</Name>
    <Prefix/>
    <Marker/>
    <MaxKeys>1000</MaxKeys>
    <IsTruncated>false</IsTruncated>
    <Contents>
        <Key>test-file1.txt</Key>
        <LastModified>2023-05-22T08:49:07.000Z</LastModified>
        <ETag>"f19cd76cd7fac68d15f0c40a063519c9"</ETag>
        <Size>1234</Size>
        <Owner>
            <DisplayName>owner</DisplayName>
            <ID>owner-id</ID>
        </Owner>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <Contents>
        <Key>004-快速部署指南&安装手册/网康_日志中心 R4.4_x64_安装手册.docx</Key>
        <LastModified>2023-05-22T08:49:07.000Z</LastModified>
        <ETag>"f19cd76cd7fac68d15f0c40a063519c9"</ETag>
        <Size>1234</Size>
        <Owner>
            <DisplayName>owner</DisplayName>
            <ID>owner-id</ID>
        </Owner>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
</ListBucketResult>`

func TestFindS3XMLString(t *testing.T) {
	xmlContent, err := FindS3XMLString(sampleXML)
	if err != nil {
		t.Fatalf("Failed to find S3 XML string: %v", err)
	}

	var result ListBucketResult
	err = xml.Unmarshal(xmlContent, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

}

func TestSanitizeXMLContent(t *testing.T) {
	xmlContent := SanitizeXMLContent([]byte(specialXML))

	var result ListBucketResult
	err := xml.Unmarshal(xmlContent, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal sanitized XML: %v", err)
	}

}

func TestParseXMLToListBucketResult(t *testing.T) {
	xmlContent := SanitizeXMLContent([]byte(sampleXML))

	result, err := ParseXMLToListBucketResult(xmlContent)
	if err != nil {
		t.Fatalf("Failed to parse XML to ListBucketResult: %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(result.Files))
	}

	if result.Files[0].Key != "test-file1.txt" {
		t.Errorf("Expected file key 'test-file1.txt', got '%s'", result.Files[0].Key)
	}
}

func TestLoadFile_GoodXML(t *testing.T) {
	if result, err := LoadFile("../test/h1.xml"); err != nil {
		log.Fatalf("Failed to load file: %v", err)
	} else {
		content, err := SaveResultToTempCSV(t, result)
		assert.NoError(t, err, "Failed to save result to temporary CSV")

		// 断言文件内容不为空
		assert.NotEmpty(t, content, "CSV content is empty")

		// 断言 CSV 文件内容是否符合预期
		expectedCSV := "Key"
		assert.Contains(t, content, expectedCSV, "CSV content does not match expected")
	}
}

func TestLoadFile_XmlInHtml(t *testing.T) {
	if result, err := LoadFile("../test/h2-html.xml"); err != nil {
		log.Fatalf("Failed to load file: %v", err)
	} else {
		content, err := SaveResultToTempCSV(t, result)
		assert.NoError(t, err, "Failed to save result to temporary CSV")

		// 断言文件内容不为空
		assert.NotEmpty(t, content, "CSV content is empty")

		// 断言 CSV 文件内容是否符合预期
		expectedCSV := "Key"
		assert.Contains(t, content, expectedCSV, "CSV content does not match expected")
	}
}

// SaveResultToTempCSV 将给定的 ListBucketResult 保存到临时 CSV 文件中，并返回文件内容
func SaveResultToTempCSV(t *testing.T, result *ListBucketResult) (string, error) {
	t.Helper()

	// 创建临时文件路径
	tmpFile := filepath.Join(t.TempDir(), "test.csv")
	t.Log(tmpFile)

	// 将结果保存到 CSV 文件中
	err := SaveResultToCSVFile(result, tmpFile)
	if err != nil {
		return "", fmt.Errorf("failed to save result to CSV file: %w", err)
	}

	// 读取文件内容
	text, err := os.ReadFile(tmpFile)
	if err != nil {
		return "", fmt.Errorf("failed to read temporary CSV file: %w", err)
	}

	return string(text), nil
}

func TestSaveToCsv(t *testing.T) {
	xmlContent := SanitizeXMLContent([]byte(sampleXML))

	result, err := ParseXMLToListBucketResult(xmlContent)
	if err != nil {
		t.Fatalf("Failed to parse XML to ListBucketResult: %v", err)
	}
	content, err := SaveResultToTempCSV(t, result)
	assert.NoError(t, err, "Failed to save result to temporary CSV")

	// 断言文件内容不为空
	assert.NotEmpty(t, content, "CSV content is empty")

	// 断言 CSV 文件内容是否符合预期
	expectedCSV := "Key"
	assert.Contains(t, content, expectedCSV, "CSV content does not match expected")
}

func TestSaveToCsvSpecial(t *testing.T) {
	xmlContent := SanitizeXMLContent([]byte(specialXML))

	result, err := ParseXMLToListBucketResult(xmlContent)
	if err != nil {
		t.Fatalf("Failed to parse XML to ListBucketResult: %v", err)
	}

	content, err := SaveResultToTempCSV(t, result)
	assert.NoError(t, err, "Failed to save result to temporary CSV")

	// 断言文件内容不为空
	assert.NotEmpty(t, content, "CSV content is empty")

	// 断言 CSV 文件内容是否符合预期
	expectedCSV := "Key"
	assert.Contains(t, content, expectedCSV, "CSV content does not match expected")
}

func TestPrintResult(t *testing.T) {
	if result, err := LoadFile("../test/h2-html.xml"); err != nil {
		log.Fatalf("Failed to load file: %v", err)
	} else {
		if err := PrintResult(result); err != nil {
			t.Fatalf("Failed to print result: %v", err)
		}
	}
}

func TestLoadRemoteHTTP(t *testing.T) {
	result, err := LoadRemoteHTTP("https://dl.qianxin.com/")
	if err != nil {
		t.Fatalf("Failed to load remote HTTP: %v", err)
	}

	content, err := SaveResultToTempCSV(t, result)
	assert.NoError(t, err, "Failed to save result to temporary CSV")

	// 断言文件内容不为空
	assert.NotEmpty(t, content, "CSV content is empty")

	// 断言 CSV 文件内容是否符合预期
	expectedCSV := "Key"
	assert.Contains(t, content, expectedCSV, "CSV content does not match expected")
}

// 翻页
func TestLoadRemoteHTTPRecursive_v1(t *testing.T) {
	result, err := LoadRemoteHTTPRecursive("https://116.178.4.134:9000/", 3)
	if err != nil {
		t.Fatalf("Failed to load remote HTTP: %v", err)
	}

	fmt.Println(result)
	rlen := len(result.Files)
	if rlen < 1000 {
		log.Fatalf("翻页有问题！%+v", rlen)
	}

	content, err := SaveResultToTempCSV(t, result)
	assert.NoError(t, err, "Failed to save result to temporary CSV")

	// 断言文件内容不为空
	assert.NotEmpty(t, content, "CSV content is empty")

	// 断言 CSV 文件内容是否符合预期
	expectedCSV := "Key"
	assert.Contains(t, content, expectedCSV, "CSV content does not match expected")
}

func TestLoadRemoteHTTPRecursive_v1_No_NextMarker(t *testing.T) {
	result, err := LoadRemoteHTTPRecursive("https://gwcache.acer.com.cn/", 3)
	if err != nil {
		t.Fatalf("Failed to load remote HTTP: %v", err)
	}

	fmt.Println(result)
	rlen := len(result.Files)
	if rlen < 1000 {
		log.Fatalf("翻页有问题！%+v", rlen)
	}

	content, err := SaveResultToTempCSV(t, result)
	assert.NoError(t, err, "Failed to save result to temporary CSV")

	// 断言文件内容不为空
	assert.NotEmpty(t, content, "CSV content is empty")

	// 断言 CSV 文件内容是否符合预期
	expectedCSV := "Key"
	assert.Contains(t, content, expectedCSV, "CSV content does not match expected")
}

func TestLoadRemoteHTTPRecursive_v2(t *testing.T) {
	result, err := LoadRemoteHTTPRecursive("https://61.174.60.37/", 2)
	if err != nil {
		t.Fatalf("Failed to load remote HTTP: %v", err)
	}

	fmt.Println(result)
	rlen := len(result.Files)
	if rlen < 1000 {
		log.Fatalf("翻页有问题！%+v", rlen)
	}

	content, err := SaveResultToTempCSV(t, result)
	assert.NoError(t, err, "Failed to save result to temporary CSV")

	// 断言文件内容不为空
	assert.NotEmpty(t, content, "CSV content is empty")

	// 断言 CSV 文件内容是否符合预期
	expectedCSV := "Key"
	assert.Contains(t, content, expectedCSV, "CSV content does not match expected")
}

func TestLoadRemoteHTTPRecursive_CanDeal_CDNIgnoreParameters(t *testing.T) {
	// 测试能处理CDN 忽略参数的场景（也就是无法翻页的情况），这种时候，不理他们

	result, err := LoadRemoteHTTPRecursive("https://dl.qianxin.com/", 3)
	if err != nil {
		t.Fatalf("Failed to load remote HTTP: %v", err)
	}

	rlen := len(result.Files)
	if rlen != 2000 { // 只有 2 页才对！
		log.Fatalf("翻页有问题！没兼容这种 CDN 的场景 %+v", rlen)
	}

	content, err := SaveResultToTempCSV(t, result)
	assert.NoError(t, err, "Failed to save result to temporary CSV")

	// 断言文件内容不为空
	assert.NotEmpty(t, content, "CSV content is empty")

	// 断言 CSV 文件内容是否符合预期
	expectedCSV := "Key"
	assert.Contains(t, content, expectedCSV, "CSV content does not match expected")

}
