# s3 viewer go
> A CLI tool that helps you view public files in s3 bucket 

## Feature
- can export file_list to `.csv`
- can preview files without messing up your local file system


## Usage
```
$ s3viewer -h
Usage of ./s3viewer:    
  -u string
      s3 URL, http://bucket.s3.amazonaws.com/ (default "http://")
  -o string
      output file name (default "file_list.csv")
 
```

    
## sample HTML format
```html
<html xmlns="http://www.w3.org/1999/xhtml"><head><style id="xml-viewer-style">/* Copyright 2014 The Chromium Authors
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

:root {
  color-scheme: light dark;
}

div.header {
    border-bottom: 2px solid black;
    padding-bottom: 5px;
    margin: 10px;
}

@media (prefers-color-scheme: dark) {
  div.header {
    border-bottom: 2px solid white;
  }
}

div.folder &gt; div.hidden {
    display:none;
}

div.folder &gt; span.hidden {
    display:none;
}

.pretty-print {
    margin-top: 1em;
    margin-left: 20px;
    font-family: monospace;
    font-size: 13px;
}

#webkit-xml-viewer-source-xml {
    display: none;
}

.opened {
    margin-left: 1em;
}

.comment {
    white-space: pre;
}

.folder-button {
    user-select: none;
    cursor: pointer;
    display: inline-block;
    margin-left: -10px;
    width: 10px;
    background-repeat: no-repeat;
    background-position: left top;
    vertical-align: bottom;
}

.fold {
    background: url("data:image/svg+xml,&lt;svg xmlns='http://www.w3.org/2000/svg' fill='%23909090' width='10' height='10'&gt;&lt;path d='M0 0 L8 0 L4 7 Z'/&gt;&lt;/svg&gt;");
    height: 10px;
}

.open {
    background: url("data:image/svg+xml,&lt;svg xmlns='http://www.w3.org/2000/svg' fill='%23909090' width='10' height='10'&gt;&lt;path d='M0 0 L0 8 L7 4 Z'/&gt;&lt;/svg&gt;");
    height: 10px;
}
</style></head><body><div id="webkit-xml-viewer-source-xml"><ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Name>bofiles</Name><Prefix/><Marker/><MaxKeys>1000</MaxKeys><IsTruncated>true</IsTruncated><Contents><Key>book/</Key><LastModified>2021-06-24T03:30:23.000Z</LastModified><ETag>"d41d8cd98f00b204e9800998ecf8427e"</ETag><Size>0</Size><Owner><DisplayName/><ID>szjinxingwei</ID></Owner><StorageClass>STANDARD</StorageClass></Contents><Contents><Key>book/list.txt</Key><LastModified>2023-05-22T08:49:07.000Z</LastModified><ETag>"f19cd76cd7fac68d15f0c40a063519c9"</ETag><Size>1499</Size><Owner><DisplayName/><ID>szjinxingwei</ID></Owner><StorageClass>STANDARD</StorageClass></Contents><Contents><Key>book/中华书局版历史七年级上册_5344_20201016.jlk</Key><LastModified>2021-06-24T03:33:36.000Z</LastModified><ETag>"abbd6a5d7dc71bc625fde1b6ed702b7e"</ETag><Size>5803443</Size><Owner><DisplayName/><ID>szjinxingwei</ID></Owner><StorageClass>STANDARD</StorageClass></Contents><
```