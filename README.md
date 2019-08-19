# pinyin

一个简单的中文转拼音工具。本工具是根据 https://github.com/chain-zhang/pinyin 的代码修改而来。主要的改动如下：

1. 将拼音码文件直接放到代码中，防止因为网络问题导致的读取时间不可控问题；
2. 对中文转换进行调整，保留原有标点符号、数字等非中文字符。


## 使用方式

### 需要自定义模式
```go
// 定义一个pinyin结构体
py := New(cnStr)
// 设置模式
py.Mode(ModeTone)
// 转换， pyStr为中文转换后的拼音
pyStr, err := py.Convert()
```

### 使用默认的转换参数（拼音不含声调）
```go
// 转换， pyStr为中文转换后的拼音
pyStr, err := pinyin.Convert("这里放入需要转换的中文")
```