# logsvr 日志服务器

用于接收DUT的 syslog，会读取 `config.Dut` 的配置，只接收 DUT 的日志。

## 使用

启动服务器函数：

- logsvr.RunWithSessionFile(file string)
- logsvr.RunWithSessionWriter(w io.WriteCloser)
- logsvr.Run()

停止服务器函数 `logsvr.Stop()`

## 怎么使用？

先启动服务，然后操作DUT设备触发日志的产生，再调用 `PourOutLastMessages() []string` 获取服务器接收到的最新日志，比对信息是否符合预期。

>每次调用 `PourOutLastMessages()[]string` 都会清空暂存区信息。
