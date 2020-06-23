编译windows测试
参考URL：https://www.cnblogs.com/zsy/p/12008141.html
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go