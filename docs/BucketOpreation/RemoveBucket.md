### RemoveBucket

### (ctx context.Context, bucketName string)

###  error

删除一个存储桶，存储桶必须为空才能被成功删除。

__参数__


| 参数         | 类型              | 描述       |
| :----------- | :---------------- | :--------- |
| `ctx`        | _context.Context_ | 上下文控制 |
| `bucketName` | _string_          | 存储桶名称 |

__示例：删除一个空的存储桶__


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

//删除桶
err = client.RemoveBucket(context.Background(), bucketName)
if err != nil {
    fmt.Println(err)
    return
}

fmt.Println("delete bucket successful")

```

