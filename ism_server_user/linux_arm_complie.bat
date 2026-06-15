SET CGO_ENABLED=1
SET GOOS=linux
SET GOARCH=arm
set GOARM=5
set CC=aarch64-linux-gcc
set CXX=aarch64-linux-g++


yum install alsa-lib-devel -y

go build -ldflags "-w -s"

CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOARM=5 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ go build -ldflags "-w -s" -trimpath
