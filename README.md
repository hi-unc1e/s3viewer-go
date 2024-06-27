# s3 viewer go
> A CLI tool that helps you view public files in s3 bucket 

> 本工具能帮你更轻松地查看 s3 bucket 中的公开文件
<br>

[![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)](#)
<a href="https://goreportcard.com/report/github.com/hi-unc1e/s3viewer-go"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/hi-unc1e/s3viewer-go"/></a>
<a href="https://github.com/hi-unc1e/s3viewer-go/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/github/license/hi-unc1e/s3viewer-go"/></a>
[![Release](https://github.com/hi-unc1e/s3viewer-go/actions/workflows/releaser.yml/badge.svg)](https://github.com/hi-unc1e/s3viewer-go/actions/workflows/releaser.yml)

<a href="https://github.com/hi-unc1e/s3viewer-go/stargazers"><img alt="GitHub stars" src="https://img.shields.io/github/stars/hi-unc1e/s3viewer-go"/></a>
<a href="https://github.com/hi-unc1e/s3viewer-go/releases"><img alt="GitHub releases" src="https://img.shields.io/github/release/hi-unc1e/s3viewer-go"/></a>
<a href="https://github.com/hi-unc1e/s3viewer-go/releases"><img alt="Downloads" src="https://img.shields.io/github/downloads/hi-unc1e/s3viewer-go/total?color=brightgreen"/></a>
----
## Feature
- can export file_list to `.csv`
- can preview files without messing up your local file system

- 可以将`file_list`导出为`.csv`
- 可以预览文件，而不会搞乱你的本地文件系统


## Usage
### Install（安装）
```bash
$ make
```
### CLI（命令行含义）
```
$ s3viewer -h
Usage of ./s3viewer:    
  -u string
      s3 URL, http://bucket.s3.amazonaws.com/ (default "http://")
  -p int
      max page (default 1)
  -o string
      output file name
  -web
        preview via local_web, such as http://127.0.0.1:30028/static/index.html
```
## ToDo
- [x] 指定`-o`参数时，将会保存 csv 到本地
    ```bash
    $ ./s3v -u https://dl.qianxin.com/ -o qianxin.csv
    2024/06/23 11:41:01 Saved into `qianxin.csv`
    ```
- [x] 不指定`-o`参数时，将会打印出列表，例如下面的命令将会输出
  ```bash
    $  ./s3v -u http://downbs.wan896.com/
  ```
  
  | Key                                                         | Size      | LastModifiedDate         |
  | ----------------------------------------------------------- | --------- | ------------------------ |
  | 20240112144418.png                                          | 347934    | 2024-01-17T09:25:17.000Z |
  | package/android/pzjhcs/35a336f13be958cdc7fefa31a6e953d5.apk | 463869061 | 2024-04-08T06:16:28.000Z |


- [x] 自动翻页，统计文件总数
  ```bash
  $ ./s3v -u https://s3_url/ -o mp.csv -p 2
  ……
  2024/06/23 17:01:45 [+]结果总条数: [1899], 已拉取页数: [2]
  2024/06/23 17:01:45 Saved into mp.csv
  ```
- [x] 下载链接，自动拼接链接～
- [x] 支持预览图片（浏览器支持啥，我就支持啥）
- [ ] 支持预览表格（开发中）
- [ ] 读取文件内容(xlsx，pdf，docx)，判断是否是敏感信息
- [ ] 支持判断文件类型（在考虑要不要开发）

```html
fofa dork: https://fofa.info/result?qbase64=IjxMaXN0QnVja2V0UmVzdWx0IHhtbG5zPVwiaHR0cDovL3MzLmFtYXpvbmF3cy5jb20vZG9jLzIwMDYtMDMtMDEvXCI%2BIiAmJiBjb3VudHJ5PSJDTiIgJiYgaWNvbl9oYXNoPSIyMTAwMDcyMDYyIg%3D%3D


```

## Changelog
| Date       | Desc                                 |
| ---------- |--------------------------------------|
| 2024-06-23 | 发布`v0.0.4`，支付保存为 csv文件，支持翻页。         |
| 2024-06-20 | 项目需要：搞一个 s3 导出工具，支持导出 excel的那种（便于检索） |
|            |                                      |


## ScreenShots

![image](https://github.com/hi-unc1e/s3viewer-go/assets/67778054/e60da15d-6a9d-4582-9fa5-fa9edcbd0331)


![image](https://github.com/hi-unc1e/s3viewer-go/assets/67778054/1013bc8f-c1c2-4b91-b39e-af2589be659c)


    
## sample S3 HTML format
```html

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
</ListBucketResult>
```
