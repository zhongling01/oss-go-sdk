### PutObject

### (ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,opts PutObjectOptions) 
### (info UploadInfo, err error)

当对象小于128MiB时，直接在一次PUT请求里进行上传。当大于128MiB时，根据文件的实际大小，PutObject会自动地将对象进行拆分成128MiB一块或更大一些进行上传。对象的最大大小是5TB。

__参数__


| 参数         | 类型                         | 描述                                        |
| :----------- | :--------------------------- | :------------------------------------------ |
| `ctx`        | _context.Context_            | 上下文控制                                  |
| `bucketName` | _string_                     | 存储桶名称                                  |
| `objectName` | _string_                     | 对象的名称                                  |
| `reader`     | _io.Reader_                  | 任意实现了io.Reader的GO类型                 |
| `objectSize` | _int64_                      | 上传的对象的大小，-1代表未知,会采用流式传输 |
| `opts`       | _ossClient.PutObjectOptions_ | 允许用户设置上传选项                        |

__ossClient.PutObjectOptions__

**返回**

| UploadInfo     | 类型        | 描述               |
| -------------- | ----------- | ------------------ |
| `Bucket`       | _string_    | 上传的存储桶名称   |
| `Key`          | _string_    | 上传的新对象名称   |
| `ETag`         | _string_    | 上传的新对象的ETag |
| `Size`         | _int64_     | 上传的新对象大小   |
| `LastModified` | _time.Time_ | 对象最后更新的时间 |
| `VersionID`    | _string_    | 上传对象的新版本号 |


__示例__


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
//上传普通文件
reader, err := os.Open(filePath)
if err != nil {
    fmt.Printf("Error opening file:%v", err)
    return
}

stat, err := os.Stat(filePath)
if err != nil {
    fmt.Printf("Error stat file:%v", err)
    return
}
size := stat.Size()

opt := ossClient.PutObjectOptions{}

info, err := client.PutObject(context.Background(), bucketName, objectName, reader, size, opt)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(info)

//以流式上传
info, err = client.PutObject(context.Background(), bucketName, objectName, reader, -1, opt)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(info)


```
