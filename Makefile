# Makefile

# 设置可执行文件名（就叫`s3v`，简单点）
EXECUTABLE := ./s3v

# 设置Go源码文件所在目录
SRC := cmd/main.go

# 设置编译时的ldflagscd
LDFLAGS := -w -s

# 编译目标
all: $(EXECUTABLE)

$(EXECUTABLE): $(SRC)
	@echo "Building $(EXECUTABLE)..."
	go build -ldflags "$(LDFLAGS)" -o $(EXECUTABLE) $(SRC)

# 清理目标，用于删除已编译的可执行文件
clean:
	@echo "Cleaning up..."
	rm -f $(EXECUTABLE)

# 运行目标
run: $(EXECUTABLE)
	@echo "Running $(EXECUTABLE)..."
	./$(EXECUTABLE)

# 安装必要的依赖或进行其他预处理（如有需要）
setup:
	@echo "Setting up environment..."
	go mod tidy

# 默认操作
.PHONY: all clean run setup
.DEFAULT_GOAL := all
