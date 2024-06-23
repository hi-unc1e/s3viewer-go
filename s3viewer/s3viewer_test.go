package s3viewer

import (
	"encoding/xml"
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

const invalidXML = `
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

	if result.Name != "example-bucket" {
		t.Errorf("Expected bucket name 'example-bucket', got '%s'", result.Name)
	}
}

func TestSanitizeXMLContent(t *testing.T) {
	xmlContent := SanitizeXMLContent([]byte(invalidXML))

	var result ListBucketResult
	err := xml.Unmarshal(xmlContent, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal sanitized XML: %v", err)
	}

	if result.Name != "example-bucket" {
		t.Errorf("Expected bucket name 'example-bucket', got '%s'", result.Name)
	}
}

func TestParseXMLToListBucketResult(t *testing.T) {
	xmlContent := SanitizeXMLContent([]byte(sampleXML))

	result, err := ParseXMLToListBucketResult(xmlContent)
	if err != nil {
		t.Fatalf("Failed to parse XML to ListBucketResult: %v", err)
	}

	if result.Name != "example-bucket" {
		t.Errorf("Expected bucket name 'example-bucket', got '%s'", result.Name)
	}

	if len(result.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(result.Files))
	}

	if result.Files[0].Key != "test-file1.txt" {
		t.Errorf("Expected file key 'test-file1.txt', got '%s'", result.Files[0].Key)
	}
}
