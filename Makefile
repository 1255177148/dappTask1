APP_NAME := myapp
MAIN := ./main.go
BIN_DIR := ./bin
# 跨平台判断
ifeq ($(OS),Windows_NT)
	# 如果是windows，加exe后缀
    EXT := .exe
    TMP_DIR := $(subst \,/,$(USERPROFILE))/AppData/Local/Temp
else
    EXT :=
    TMP_DIR := /tmp
endif
# 生成的可执行文件目录和文件名
BIN_FILE := $(BIN_DIR)/$(APP_NAME)$(EXT)

# 默认环境(dev/test/prod)
ENV ?= dev
# 设置 APP_ENV 环境变量
export APP_ENV=$(ENV)

.DEFAULT_GOAL := run

ifeq ($(OS),Windows_NT)
	# Windows 系统
	TMP_DIR := $(subst \,/,$(USERPROFILE))/AppData/Local/Temp
else
	# Linux / macOS
	TMP_DIR := /tmp
endif

# =========================
# 设置 build 命令变量
# =========================
ifeq ($(OS),Windows_NT)
	BUILD_CMD := set TMP=$(TMP_DIR) && set TEMP=$(TMP_DIR) && set GOTMPDIR=$(TMP_DIR) && set APP_ENV=$(ENV) && go build -o $(BIN_FILE) $(MAIN)
else
	# Linux / macOS
	BUILD_CMD := TMP=$(TMP_DIR) TEMP=$(TMP_DIR) GOTMPDIR=$(TMP_DIR) APP_ENV=$(ENV) go build -o $(BIN_FILE) $(MAIN)
endif

# 初始化依赖
.PHONY: init
init:
	go mod tidy
	go mod download

# 运行项目
.PHONY: run
run:
	@echo "🚀 运行环境: $(ENV)"
	go run $(MAIN)

# 构建二进制
.PHONY: build
build:
	@echo "🔨 构建 $(APP_NAME) 环境: $(ENV)"
	mkdir -p $(BIN_DIR)
	@$(BUILD_CMD)

# 清理构建产物
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

# 运行测试
.PHONY: test
test:
	APP_ENV=$(ENV) go test ./... -v

# 整理依赖
.PHONY: tidy
tidy:
	go mod tidy

# 代码检查
.PHONY: lint
lint:
	go vet ./...

.PHONY: abigen
abigen:
	abigen --abi=./contract/count.json --pkg=contract --out=./contract/count.go