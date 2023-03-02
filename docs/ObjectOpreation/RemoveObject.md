### RemoveObject

### (ctx context.Context, bucketName, objectName string, opts ossClient.RemoveObjectOptions) 

### error

删除一个对象

**参数**

| 参数         | 类型                            | 描述                 |
| :----------- | :------------------------------ | :------------------- |
| `ctx`        | *context.Context*               | 上下文控制           |
| `bucketName` | *string*                        | 存储桶名称           |
| `objectName` | *string*                        | 对象名称             |
| `opts`       | *ossClient.RemoveObjectOptions* | 允许用户设置删除选项 |

**ossClient.RemoveObjectOptions**

| 参数                    | 类型     | 描述                               |
| :---------------------- | :------- | :--------------------------------- |
| `opts.GovernanceBypass` | *bool*   | 是否绕过对象保留的治理模式         |
| `opts.VersionID`        | *string* | 填写对象版本号以删除对象的特定版本 |

**示例**

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

//删除一个对象的特定版本并绕过治理模式
opt := ossClient.RemoveObjectOptions{
    GovernanceBypass: true,
    VersionID:        "a4b7dfaa-a8b8-47f0-9ea1-f204c886da5e",
}
err = client.RemoveObject(context.Background(), bucketName, objectName, opt)
if err != nil {
    fmt.Println(err)
    return
}
```