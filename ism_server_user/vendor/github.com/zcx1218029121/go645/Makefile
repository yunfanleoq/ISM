BUILD_TIME=`date +%FT%T%z`
GITTAG=1.0
OUTPUT=./out/zoom

ARM=arm
MAC= mac
WIN =win
OUTPUT=./g645
all: ${ARM} ${MAC} ${WIN}

${ARM}:
	@echo "============ [${ARM}] 开始编译  =========="
	GOOS=linux GOARCH=arm go build -o ${OUTPUT}_${ARM}   ./example/gaea/main.go
	@echo "============[${ARM}] 编译结束============"
${MAC}:
	@echo "============ [${MAC}] 开始编译  ============"
	GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT}_${MAC}_amd64  ./main.go
	@echo "============ [${MAC}] 编译结束============"
${WIN}:
	@echo "============ [${WIN}] 开始编译  ============"
	GOOS=windows GOARCH=amd64 go build -o ${OUTPUT}_${WIN}_amd64 ./main.go
	@echo "============ [${WIN}] 编译结束  ============"