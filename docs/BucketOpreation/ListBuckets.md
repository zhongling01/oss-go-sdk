### ListBuckets

### (ctx context.Context, listRecycle bool) 

### ([]BucketInfo, error)

列出回收站内的所有桶或列出回收站外的所有桶

__参数__

| 参数          | 类型              | 描述                                                         |
| ------------- | ----------------- | ------------------------------------------------------------ |
| `ctx`         | _context.Context_ | 上下文控制                                                   |
| `listRecycle` | _bool_            | 是否为显示回收站中的桶，置为true时只显示回收站中的桶，置为false时只显示回收站之外的桶 |



| BucketInfo            | 类型        | 描述             |
| --------------------- | ----------- | ---------------- |
| `bucket.Name`         | _string_    | 存储桶名称       |
| `bucket.CreationDate` | _time.Time_ | 存储桶的创建时间 |



__示例：列出所有回收站中的桶__


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

//列出所有回收站中的桶
buckets, err := client.ListBuckets(context.Background(), true)
if err != nil {
    fmt.Println(err)
    return
}
for _, bucket := range buckets {
    fmt.Printf("桶名称:%s \n创建日期:%s\n", bucket.Name, bucket.CreationDate)
}
```

