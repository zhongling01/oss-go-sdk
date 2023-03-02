### GetBucketTagging

### (ctx context.Context, bucketName string) 

### (*tags.Tags,error)

获取存储桶的标签

__参数__

| 参数         | 类型              | 描述       |
| ------------ | ----------------- | ---------- |
| `ctx`        | _context.Context_ | 上下文控制 |
| `bucketName` | _string_          | 存储桶名称 |



| 返回    | 类型         | 描述       |
| ------- | ------------ | ---------- |
| `tags ` | _*tags.Tags_ | 标签结构体 |



__示例：获取一个桶的标签__


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

//获取存储桶的标签
res, err := client.GetBucketTagging(context.Background(), bucketName)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(res.String())
```

