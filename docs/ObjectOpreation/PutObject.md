### PutObject

### (ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,opts PutObjectOptions) 
### (info UploadInfo, err error)

上传一个对象，分为简单上传与分段上传两种模式，简单上传后对象以单个文件的形式存储，分段上传则以文件的多个分段进行存储

简单上传

1. `DisableMultpart`置为`true` 此时为简单上传 会上传结果为单个文件的对象
2. 当设置的`partSize >= objectSize`时，也会采用简单上传 会上传结果为单个文件的对象
3. 简单上传时上传对象的最大大小为5GiB

分段上传

1. 当`objectSize`为-1时启用流式传输，此时会采用分段上传，不能使`DisableMultpart`置为`true`，否则此时无论设置`partSize`为多少都会返回错误
2. 当`objectSize>=0` 同时`partSize < objectSize`时采用分段上传

当`objectSize`等于0时将上传空文件

上传对象的最大大小是48.8TB。

__参数__


| 参数         | 类型                         | 描述                                        |
| :----------- | :--------------------------- | :------------------------------------------ |
| `ctx`        | _context.Context_            | 上下文控制                                  |
| `bucketName` | _string_                     | 存储桶名称                                  |
| `objectName` | _string_                     | 对象的名称                                  |
| `reader`     | _io.Reader_                  | 任意实现了io.Reader的GO类型                 |
| `objectSize` | _int64_                      | 上传的对象的大小，-1代表未知,会采用流式传输 |
| `opts`       | _ossClient.PutObjectOptions_ | 允许用户设置上传选项                        |



| _ossClient.PutObjectOptions_ | 类型              | 描述                                                         |
| ---------------------------- | ----------------- | ------------------------------------------------------------ |
| `DisableMultipart`           | _bool_            | 是否禁用多段上传                                             |
| `PartSize`                   | _uint64_          | 分段大小，填0时默认为16MiB，同时5MiB<=PartSize，PartSize<=5GiB |
| `LegalHold`                  | _LegalHoldStatus_ | 在桶开启对象锁定的前提下，上传对象并同时设置对象的合法保留   |
| `Mode`                       | _RetentionMode_   | 在桶开启对象锁定的前提下，上传对象并同时设置对象的对象保留的保留模式 |
| `RetainUntilDate`            | _time.Time_       | 在桶开启对象锁定的前提下，上传对象并同时设置对象的对象保留的保留期 |



**返回**

| UploadInfo     | 类型        | 描述               |
| -------------- | ----------- | ------------------ |
| `Bucket`       | _string_    | 上传的存储桶名称   |
| `Key`          | _string_    | 上传的新对象名称   |
| `ETag`         | _string_    | 上传的新对象的ETag |
| `Size`         | _int64_     | 上传的新对象大小   |
| `LastModified` | _time.Time_ | 对象最后更新的时间 |
| `VersionID`    | _string_    | 上传对象的新版本号 |


__示例1：简单上传与基本流式上传__


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

```

```go
//简单上传
opt := ossClient.PutObjectOptions{
    DisableMultipart: true,
}

info, err := client.PutObject(context.Background(), bucketName, objectName, reader, size, opt)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(info)

```

```go
//以自定分片大小进行流式上传 分段大小为最大值5GiB
opt := ossClient.PutObjectOptions{
    PartSize: 1024 * 1024 * 1024 * 5,
}
info, err = client.PutObject(context.Background(), bucketName, objectName, reader, -1, opt)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(info)
```

