# 介绍

这里是《GO语言编程》一书中示例源代码。
tq增加了人保接口测试，增加了mssql存储过程调用，增加了日志(2019.4.6)
# 环境配置

1) 安装 go1

2) 打开 ~/.bashrc，加入： 

    export $GOBOOK=~/gobook  #假设代码在 ~/gobook 下
    source $GOBOOK/env.sh

3) 保存 ~/.bashrc，并 source 之

# 运行代码

对于单文件程序，如sample.go，直接运行go run sample.go即可。

对于多文件的复杂样例，直接 go run 主程序文件即可。

2019.10.14 ----
---------------------------
# 设置环境变量：
GOPATH=C:\mygo;C:\mygo\src\mygobook\chapter5
假设你已经使用了SS客户端，本地socks5代理为127.0.0.1:1080
在CMD窗口输入如下指令设置代理：

set http_proxy=socks5://127.0.0.1:1080
set https_proxy=socks5://127.0.0.1:1080
set ftp_proxy=socks5://127.0.0.1:1080
测试 curl https://www.facebook.com 能得到返回结果。

然后下载包：
go get github.com/fsnotify/fsnotify
go get github.com/spf13/viper
go get github.com/lexkong/log


取消代理命令：

set http_proxy=
set https_proxy=
set ftp_proxy=
*设置代理后只对当前命令行窗口生效，重新打开CDM需要再次设置。
