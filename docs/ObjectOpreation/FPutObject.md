### FPutObject

### (ctx context.Context, bucketName,  objectName,  filePath string, opts PutObjectOptions)

###  (lnfo UploadInfo,err error)

将filePath对应的文件内容上传到一个对象中。

当对象小于128MiB时，FPutObject直接在一次PUT请求里进行上传。当大于128MiB时，根据文件的实际大小，FPutObject会自动地将对象进行拆分成128MiB一块或更大一些进行上传。对象的最大大小是5TB。

__参数__


| 参数         | 类型                         | 描述                 |
| :----------- | :--------------------------- | :------------------- |
| `ctx`        | _context.Context_            | 上下文控制           |
| `bucketName` | _string_                     | 存储桶名称           |
| `objectName` | _string_                     | 对象的名称           |
| `filePath`   | _string_                     | 要上传的文件的路径   |
| `opt`        | _ossClient.PutObjectOptions_ | 允许用户设置上传选项 |



| _ossClient.PutObjectOptions_ | 类型              | 描述                                                         |
| ---------------------------- | ----------------- | ------------------------------------------------------------ |
| `DisableMultipart`           | _bool_            | 是否禁用多段上传                                             |
| `LegalHold`                  | _LegalHoldStatus_ | 在桶开启对象锁定的前提下，上传对象并同时设置对象的合法保留   |
| `Mode`                       | _RetentionMode_   | 在桶开启对象锁定的前提下，上传对象并同时设置对象的对象保留的保留模式 |
| `RetainUntilDate`            | _time.Time_       | 在桶开启对象锁定的前提下，上传对象并同时设置对象的对象保留的保留期 |



| LegalHoldStatus可用项         | 类型     | 描述         |
| ----------------------------- | -------- | ------------ |
| `ossCleint.LegalHoldEnabled`  | _string_ | 启用合法保留 |
| `ossCleint.LegalHoldDisabled` | _string_ | 禁用合法保留 |



| RetentionMode可用项    | 类型     | 描述     |
| ---------------------- | -------- | -------- |
| `ossCleint.Governance` | _string_ | 治理模式 |
| `ossCleint.Compliance` | _string_ | 监管模式 |

**返回**

| UploadInfo     | 类型        | 描述               |
| -------------- | ----------- | ------------------ |
| `Bucket`       | _string_    | 上传的存储桶名称   |
| `Key`          | _string_    | 上传的新对象名称   |
| `ETag`         | _string_    | 上传的新对象的ETag |
| `Size`         | _int64_     | 上传的新对象大小   |
| `LastModified` | _time.Time_ | 对象最后更新的时间 |
| `VersionID`    | _string_    | 上传对象的新版本号 |



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

__示例1:简单上传文件__


```go
//普通上传文件
opt := ossClient.PutObjectOptions{}

_, err = client.FPutObject(context.Background(), bucketName, objectName, filePath, opt)
if err != nil {
    fmt.Println(err)
    return
}
```

__示例2:上传文件对象的同时设置对象的合法保留__


```go
//上传文件对象的同时设置对象的合法保留
opt := ossClient.PutObjectOptions{
    LegalHold: ossClient.LegalHoldEnabled,
}

_, err = client.FPutObject(context.Background(), bucketName, objectName, uploadFilePath, opt)
if err != nil {
    fmt.Println(err)
    return
}
```

__示例3:上传文件对象的同时设置对象的对象保留__


```go
//上传文件对象的同时设置对象的对象保留
opt := ossClient.PutObjectOptions{
    Mode:            ossClient.Governance,                       //设置为治理模式
    RetainUntilDate: time.Now().Add(24 * 60 * 60 * time.Second), //保留期为一天
}

_, err = client.FPutObject(context.Background(), bucketName, objectName, uploadFilePath, opt)
if err != nil {
    fmt.Println(err)
    return
}
```

