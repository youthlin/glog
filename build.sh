#! /bin/bash
set -ex        # 回显命令

PRE_PWD=$(pwd) # 保存最初的工作目录

# 脚本所在目录
PWD=$(
  cd "$(dirname "$0")" || exit
  pwd
)
cd "$PWD"         # 进入脚本目录

rm -rf output     # 删除输出目录
mkdir -p output   # 创建输出目录
cp -r conf output # 拷贝配置文件
go env            # 显示一下 GO 坏境

# 执行编译
go build -v -race -o output/glog

cd "$PRE_PWD" #回到最初工作目录
