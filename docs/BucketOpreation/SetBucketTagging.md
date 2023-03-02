### SetBucketTagging

### (ctx context.Context, bucketName string, tags *tags.Tags) 

### error

设置存储桶的标签

__参数__

| 参数         | 类型              | 描述       |
| ------------ | ----------------- | ---------- |
| `ctx`        | _context.Context_ | 上下文控制 |
| `bucketName` | _string_          | 存储桶名称 |
| `tags `      | _*tags.Tags_      | 标签结构体 |

__示例：设置一个桶的标签__


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

//设置新的桶标签
tagMap := make(map[string]string)
tagMap["key"] = "value"

newTag, err := tags.NewTags(tagMap, false)
if err != nil {
    fmt.Println(err)
    return
}

err = client.SetBucketTagging(context.Background(), bucketName, newTag)
if err != nil {
    fmt.Println(err)
    return
}

res, err := client.GetBucketTagging(context.Background(), bucketName)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(res.String())
```

