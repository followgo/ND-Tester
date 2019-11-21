# mylogrus

logrus, lumberjack 的包装，使 logrus 支持按日志文件按大小进行滚动

## Example

```go
func main() {
	// 设置 logrus 包中的 std 实例
	mylogrus.SetStdLogrus(mylogrus.DefaultOption)

	// 测试
	logrus.Info("a info log")
	logrus.Info("an error log")
	logrus.WithField("key","123").Println("Hello World")
}
```