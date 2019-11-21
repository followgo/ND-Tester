# configurator 配置装载器

配置文件加载、保存、监视。

支持的文件类型：

- yaml
- toml
- json

## Example

```go
func main() {
	// obj 对象
	obj := struct {
		Name string
		Rank int
	}{"abc", 9}

	// filename 配置文件路径
	filename := "abc.yaml"

	// 配置管理器
	c := configurator.NewConfigurator(filename, &obj)

	// 保存配置
	if err := c.Save("测试文件"); err != nil {
		log.Fatalln(err)
	}
	defer os.Remove(filename)

	// 加载配置
	obj.Name = "123"
	obj.Rank = 0
	if err := c.Load(); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", obj)
}
```