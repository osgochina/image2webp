#!/usr/bin/env bash

# 使用镜像
export  GOPROXY=https://proxy.golang.com.cn,direct


# 如果go bin 不存在，则去环境变量中查找
if [ ! -x "$goBin" ]; then
    goBin=$(which go)
fi
if [ ! -x "$goBin" ]; then
    echo "No goBin found."
    exit 2
fi

# 编译时间
build_date=$(date +"%Y-%m-%d %H:%M:%S")
# 编译时候当前git的commit id
build_git=$(git rev-parse --short HEAD)

# 编译的golang版本
go_version=$(${goBin} version)

#编译版本
if [ -z "$build_version" ]; then
    build_version="$build_git"
fi



ldflags=()

# 链接时设置变量值
ldflags+=("-X" "\"main.BuildVersion=${build_version}\"")
ldflags+=("-X" "\"github.com/osgochina/dmicro/dserver.BuildVersion=${build_version}\"")
ldflags+=("-X" "\"github.com/osgochina/dmicro/dserver.BuildGoVersion=${go_version}\"")
ldflags+=("-X" "\"github.com/osgochina/dmicro/dserver.BuildGitCommitId=${build_git}\"")
ldflags+=("-X" "\"github.com/osgochina/dmicro/dserver.BuildTime=${build_date}\"")

CGO_ENABLED=0

if [ "${1}" == "windows" ]
then
    echo "The Windows operating system is not supported!"
    exit;
elif [ "${1}" == "osx" ]
then
    go build -v -ldflags "${ldflags[*]} -s -w"  -o image2webp-darwin-${2}
else
    go build -v -ldflags "${ldflags[*]} -s -w"  -o image2webp-${1}-${2}
fi

echo "build done!"
