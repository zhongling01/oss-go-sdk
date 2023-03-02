### GetObject

### (ctx context.Context, bucketName, objectName string, opts GetObjectOptions)

###  (*Object, error)

返回对象数据的流，error是读流时经常抛的那些错。

__参数__


| 参数         | 类型                         | 描述                          |
| :----------- | :--------------------------- | :---------------------------- |
| `ctx`        | _context.Context_            | 请求上下文（Request context） |
| `bucketName` | _string_                     | 存储桶名称                    |
| `objectName` | _string_                     | 对象的名称                    |
| `opts`       | _ossClient.GetObjectOptions_ | GET请求的一些额外参数         |



| _ossClient.GetObjectOptions_ | 类型   | 描述                                                 |
| ---------------------------- | ------ | ---------------------------------------------------- |
| `VersionID`                  | string | 在此填写目标对象的某个版本号，可下载对象的某一特定版 |



**返回**

**Object**

| 参数         | 类型         | 描述       |
| ------------ | ------------ | ---------- |
| `objectInfo` | _ObjectInfo_ | 对象的信息 |

| ObjectInfo       | 类型        | 描述                           |
| ---------------- | ----------- | ------------------------------ |
| `ETag`           | _string_    | 下载下来的对象的ETag           |
| `Key`            | _string_    | 对象名称                       |
| `LastModified`   | _time.Time_ | 对象最后被修改的日期           |
| `Size`           | _int64_     | 对象大小                       |
| `ContentType`    | _string_    | MIME标准描述对象数据格式的类型 |
| `Expires`        | _time.Time_ | 可能存在的对象过期时间         |
| `IsLatest`       | _bool_      | 下载下来的对象是否为最新版本   |
| `IsDeleteMarker` | _bool_      | 对象是否为删除标记             |
| `VersionID`      | _string_    | 下载对象的版本号               |
| `Err`            | _error_     | 可能的错误                     |




__示例:获取桶中对象并写入本地文件__


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

//获取桶中对象并写入本地文件
//打开本地文件
writer, err := os.Create("./test3.txt")
if err != nil {
    fmt.Println(err)
    return
}
defer func(writer *os.File) {
    err := writer.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
}(writer)
//获取对象
opt := ossClient.GetObjectOptions{}

reader, err := client.GetObject(context.Background(), bucketName, objectName, opt)
if err != nil {
    fmt.Println(err)
    return
}

defer func(reader *ossClient.Object) {
    err = reader.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
}(reader)
//写入文件
stat, err := reader.Stat()
if err != nil {
    fmt.Println(err)
    return
}

_, err = io.CopyN(writer, reader, stat.Size)
if err != nil {
    fmt.Println(err)
    return
}
```

