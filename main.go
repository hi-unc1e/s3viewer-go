package main

import (
	"encoding/xml"
	"fmt"
	"log"
)

// S3BucketListResponse 是 S3 XML 输出的结构映射
type S3BucketListResponse struct {
	XMLName  xml.Name `xml:"ListBucketResult"`
	Name     string   `xml:"Name"`
	Contents []struct {
		Key          string `xml:"Key"`
		LastModified string `xml:"LastModified"`
		Size         int    `xml:"Size"`
	} `xml:"Contents"`
}

func main() {
	// 假设 s3XML 是包含 XML 数据的字符串
	s3XML := `...` // 这里应该是你的 XML 输出

	// 解析 XML
	var listResp S3BucketListResponse
	err := xml.Unmarshal([]byte(s3XML), &listResp)
	if err != nil {
		log.Fatalf("Error unmarshalling XML: %v", err)
	}

	// 创建新的 Excel 文件
	f := excelize.NewFile()

	// 在 Excel 中创建一个新的工作表
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	// 设置表格的表头
	headers := map[string]string{"A1": "文件名", "B1": "上次编辑日期", "C1": "大小"}
	for col, val := range headers {
		f.SetCellValue(sheet, col, val)
	}

	// 填充数据
	for i, content := range listResp.Contents {
		row := i + 2 // 从第二行开始填充数据，因为第一行是表头
		f.SetCellValue(sheet, "A"+fmt.Sprint(row), content.Key)
		f.SetCellValue(sheet, "B"+fmt.Sprint(row), content.LastModified)
		f.SetCellValue(sheet, "C"+fmt.Sprint(row), content.Size)
	}

	// 保存表格
	err = f.SaveAs("S3BucketList.xlsx")
	if err != nil {
		log.Fatalf("无法保存文件, %v", err)
	}

	fmt.Println("Excel 表格已成功创建。")
}
