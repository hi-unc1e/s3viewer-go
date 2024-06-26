package web

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Image 结构体用于存储图片信息
type Image struct {
	Name string
	URL  string
}

// ImagePage 结构体用于存储页面上的所有图片信息
type ImagePage struct {
	Title  string
	Images []Image
}

// generateImagePage 函数生成HTML页面并保存到本地
func generateImagePage(imageURLs []string, outputPath string) error {
	// 创建一个ImagePage实例，包含所有图片信息
	page := ImagePage{
		Title:  "图片展示",
		Images: make([]Image, len(imageURLs)),
	}

	// 填充图片信息
	for i, url := range imageURLs {
		// 从URL中提取文件名
		name := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
		page.Images[i] = Image{Name: name, URL: url}
	}

	// 定义HTML模板
	tmpl, err := template.New("imagePage").Parse(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<link rel="stylesheet" href="https://cdn.staticfile.org/twitter-bootstrap/4.5.0/css/bootstrap.min.css">
</head>
<body>
<div class="container">
<div class="row">
{{range .Images}}
<div class="col-md-4">
	<div class="card mb-4 shadow-sm">
		<img src="{{.URL}}" class="card-img-top" alt="{{.Name}}">
		<div class="card-body">
			<h5 class="card-title">{{.Name}}</h5>
			<p class="card-text"><a href="{{.URL}}" download="{{.Name}}" class="btn btn-primary">下载图片</a></p>
		</div>
	</div>
</div>
{{end}}
</div>
</div>
<script src="https://cdn.staticfile.org/jquery/3.5.1/jquery.min.js"></script>
<script src="https://cdn.staticfile.org/popper.js/1.16.0/umd/popper.min.js"></script>
<script src="https://cdn.staticfile.org/twitter-bootstrap/4.5.0/js/bootstrap.min.js"></script>
</body>
</html>
`)
	if err != nil {
		return err
	}

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	// 创建HTML文件
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用模板生成HTML页面并写入文件
	return tmpl.Execute(file, page)
}

// staticHandler 用于提供静态文件服务
func staticHandler(w http.ResponseWriter, r *http.Request) {
	// 设置静态文件服务目录为./static
	staticDir := "./static"
	// 获取请求的文件路径
	filePath := r.URL.Path[len("/static"):]

	// 读取本地文件并返回
	http.ServeFile(w, r, filepath.Join(staticDir, filePath))
}

// 获取可用端口
func GetAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil

}

func ServeHttp(imageURLs []string) {
	// 生成HTML页面并保存在./static/index.html
	if err := generateImagePage(imageURLs, "./static/index.html"); err != nil {
		fmt.Println("生成页面出错:", err)
		return
	}

	// 设置HTTP服务器的路由
	http.HandleFunc("/static/", staticHandler)

	// 启动HTTP服务器
	// 创建一个TCP Listener，自动选择端口号
	port, err := GetAvailablePort() // :0 表示随机选择一个端口号
	if err != nil {
		log.Fatal("监听端口失败:", err)
	}
	fmt.Printf("服务器正在启动，访问 http://localhost:%v/static/index.html 查看图片展示页面\n", port)
	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", port), nil); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

func main() {
	// 图片URL列表
	imageURLs := []string{
		"https://182.151.13.56/FktrDTIrIu772991c5e27bf55e1da2f5cb39a34a8839_20231122212957A493.png",
		"https://182.151.13.56/FnlBd3T9euFh8f8b0bcc5fe8a7c24aba59f15f347ac1_20230818230913A068.ico",
		"https://182.151.13.56/FxGmAai2qeZ3b2333e3255fc93f3a7bd9cf1ee7c4de7_20230818201741A012.png",
		"https://182.151.13.56/FypVHPawYd4b2991c5e27bf55e1da2f5cb39a34a8839_20231122224230A499.png",
		"https://182.151.13.56/G6At7XwmCb4p76212d298a6c0013ee9a52e6e09ef6fd_20230416170633A270.jpg",
		"https://182.151.13.56/G6ROKAzNGOqd59e77230759a84ec6a25ed79da39f3cd_20230816200911A298.png",
		"https://182.151.13.56/G7KGsHH0UbbQ2991c5e27bf55e1da2f5cb39a34a8839_20231122212127A489.png",
		"https://182.151.13.56/G8AJT17l1jEZd58afedf384c73bb3939c24627e3442b_20230813224905A273.png",
		// 添加更多图片URL
	}

	ServeHttp(imageURLs)
}
