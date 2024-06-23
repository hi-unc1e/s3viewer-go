# s3 viewer go
> A CLI tool that helps you view public files in s3 bucket 

> 本工具能帮你更轻松地查看 s3 bucket 中的公开文件
<br>

[![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)](#)
<a href="https://goreportcard.com/report/github.com/hi-unc1e/s3viewer-go"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/hi-unc1e/s3viewer-go"/></a>
<a href="https://github.com/hi-unc1e/s3viewer-go/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/github/license/hi-unc1e/s3viewer-go"/></a>

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
  -o string
      output file name (default "file_list.csv")
 
```
## ToDo
- [x] 指定`-o`参数时，将会保存 csv 到本地
    ```bash
    $ ./s3v -u https://dl.qianxin.com/ -o qianxin.csv
    2024/06/23 11:41:01 Saved into `qianxin.csv`
    ```
- [ ] 不指定`-o`参数时，将会打印出列表，例如
    ```bash
    $ ./s3v -u https://dl.qianxin.com/
    TODO
    ```
- [ ] 自动翻页


## Changelog
| Date       | Desc                                 |
| ---------- |--------------------------------------|
| 2024-06-23 | 发布`v0.0.1`，至少能用了。                    |
| 2024-06-20 | 项目需要：搞一个 s3 导出工具，支持导出 excel的那种（便于检索） |
|            |                                      |



## ToDo

    
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