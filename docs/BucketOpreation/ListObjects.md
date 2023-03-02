### ListObjects

### (ctx context.Context, bucketName string, opts ListObjectsOptions)

###  <-chan ObjectInfo

列举存储桶里的对象。

__参数__


| 参数         | 类型                 | 描述       |
| :----------- | :------------------- | :--------- |
| `ctx`        | _context.Context_    | 上下文控制 |
| `bucketName` | _string_             | 存储桶名称 |
| `opts`       | _ListObjectsOptions_ | 列举选项   |



| ListObjectsOptions | 类型     | 描述                                                         |
| ------------------ | -------- | ------------------------------------------------------------ |
| `Prefix`           | _string_ | 填写`Prefix`将列举出以`Prefix`为前缀的对象                   |
| `recursive`        | _bool_   | `true`代表递归查找，`false`代表类似文件夹查找，以'/'分隔，不查子文件夹。默认为`false` |
| `MaxKeys`          | _int_    | 每次请求的返回最大数量,最大值为1000，设置超出1000则会应用为1000 |
| `WithVersions`     | _bool_   | `true`代表在列出的信息中包含对象的所有版本的版本号，`false`则不会包含，默认为`false` |



__返回值__

| 参数 | 类型                        | 描述                                             |
| :--- | :-------------------------- | :----------------------------------------------- |
| `ch` | _chan ossClient.ObjectInfo_ | 存储桶中所有对象的read channel，对象的格式如下： |

__ossClient.ObjectInfo__

| 属性               | 类型        | 描述               |
| :----------------- | :---------- | :----------------- |
| `obj.Key`          | _string_    | 对象的名称         |
| `obj.Size`         | _int64_     | 对象的大小         |
| `obj.ETag`         | _string_    | 对象的MD5校验码    |
| `obj.LastModified` | _time.Time_ | 对象的最后修改时间 |
| `obj.VersionID`    | _string_    | 对象版本号         |


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

__示例：列出具有特定前缀的对象__

```go
//列出具有特定前缀的对象
prefix := "examplePrefix"

opts := ossClient.ListObjectsOptions{
    Prefix: prefix,
}
ch := client.ListObjects(context.Background(), bucketName, opts)
for obj := range ch {
    fmt.Printf("前缀为%s 的对象有:%s\n", prefix, obj.Key)
}
```


__示例：递归列出对象__

```go
//递归列出对象
opts := ossClient.ListObjectsOptions{
    Recursive: true,
}
ch := client.ListObjects(context.Background(), bucketName, opts)
for obj := range ch {
    fmt.Println(obj.Key)
}
```

__示例：列出对象并携带他们的所有版本的版本号__

```go
//列出对象，并携带他们的所有版本的版本号
opts := ossClient.ListObjectsOptions{
	WithVersions: false,
}
ch := client.ListObjects(context.Background(), bucketName, opts)
for obj := range ch {
	fmt.Println(obj.Key + ":" + obj.VersionID)
}
```

