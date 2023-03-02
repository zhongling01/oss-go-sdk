### BucketExists

### (ctx context.Context, bucketName string) 

### (bool, error)

验证桶是否存在，前提是你有权限访问这个桶

__参数__

| 参数          | 类型              | 描述                     |
| ------------- | ----------------- | ------------------------ |
| `ctx`         | _context.Context_ | 上下文控制               |
| `bucketName ` | _string_          | 需要验证存在的存储桶名称 |



| 返回     | 类型   | 描述           |
| -------- | ------ | -------------- |
| `exists` | _bool_ | 存储桶是否存在 |
| `err`    | _err_  | 标准Error      |



__示例：验证桶是否存在__


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

//验证桶是否存在
exists, err := client.BucketExists(context.Background(), bucketName)
if err != nil {
    fmt.Println(err)
    return
}
if exists {
    fmt.Println("此桶存在！")
} else {
    fmt.Println("此桶不存在！")
}
```

