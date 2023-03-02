### FGetObject

### (ctx context.Context,bucketName, objectName, filePath string, opts GetObjectOptions)

###  error

下载并将文件保存到本地文件系统。

__参数__


| 参数         | 类型                         | 描述                          |
| :----------- | :--------------------------- | :---------------------------- |
| `ctx`        | _context.Context_            | 请求上下文（Request context） |
| `bucketName` | _string_                     | 存储桶名称                    |
| `objectName` | _string_                     | 对象的名称                    |
| `filePath`   | _string_                     | 下载后保存的路径              |
| `opts`       | _ossClient.GetObjectOptions_ | GET请求的一些额外参数         |



| _ossClient.GetObjectOptions_ | 类型   | 描述                                                 |
| ---------------------------- | ------ | ---------------------------------------------------- |
| `VersionID`                  | string | 在此填写目标对象的某个版本号，可下载对象的某一特定版 |



```go
//初始化客户端
opts := &ossClient.Options{
	Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
}
client, err := ossClient.New("127.0.0.1:19000", opts)
if err != nil {
	fmt.Println(err)
	return
}
```


__示例1:普通下载文件__


```go
//普通下载文件
opt := ossClient.GetObjectOptions{}

err = client.FGetObject(context.Background(), bucketName, objectName, filePath, opt)
if err != nil {
    fmt.Println(err)
    return
}
```

__示例2:下载对象的特定版本__


```go
//下载对象的特定版本
opt := ossClient.GetObjectOptions{
    VersionID: "e9a4678d-56ef-472e-a76f-3c04363f58a6",
}

err = client.FGetObject(context.Background(), bucketName, objectName, filePath, opt)
if err != nil {
    fmt.Println(err)
    return
}
```

