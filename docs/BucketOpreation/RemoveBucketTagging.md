### RemoveBucketTagging

### (ctx context.Context, bucketName string) 

### error

移除存储桶的所有标签

__参数__

| 参数         | 类型              | 描述       |
| ------------ | ----------------- | ---------- |
| `ctx`        | _context.Context_ | 上下文控制 |
| `bucketName` | _string_          | 存储桶名称 |



__示例：移除桶的所有标签__


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

//移除桶的所有标签
err = client.RemoveBucketTagging(context.Background(), bucketName)
if err != nil {
    fmt.Println(err)
    return
}

//获取桶的标签检验结果
res, err := client.GetBucketTagging(context.Background(), bucketName)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(res.String())
```

