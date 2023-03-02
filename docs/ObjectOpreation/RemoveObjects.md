### RemoveObjects

### (ctx context.Context, bucketName string, objectsCh <-chan ObjectInfo, opts RemoveObjectsOptions) 

### <-chan RemoveObjectError

从一个input channel里删除一个对象集合。一次发送到服务端的删除请求最多可删除1000个对象。通过error channel返回的错误信息。

__参数__


| 参数         | 类型                             | 描述                  |
| :----------- | :------------------------------- | :-------------------- |
| `ctx`        | _context.Context_                | 上下文控制            |
| `bucketName` | _string_                         | 存储桶名称            |
| `objectsCh`  | _chan string_                    | 要删除的对象的channel |
| `opts`       | _ossClient.RemoveObjectsOptions_ | 允许用户设置删除选项  |

__ossClient.RemoveObjectsOptions__

| __ossClient.RemoveObjectsOptions__ | 类型   | 描述                       |
| ---------------------------------- | ------ | -------------------------- |
| `GovernanceBypass`                 | _bool_ | 是否绕过对象保留的治理模式 |

**返回**

| 参数      | 类型                       | 描述                                       |
| --------- | -------------------------- | ------------------------------------------ |
| `errorCh` | _<-chan RemoveObjectError_ | 删除时观察到的错误的Receive-only channel。 |

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


__示例1:删除一个桶中的所有对象的所有版本,并绕过治理模式的对象保留__


```go
//删除一个桶中的所有对象的所有版本,并绕过治理模式的对象保留
objectCh := client.ListObjects(context.Background(), bucketName, ossClient.ListObjectsOptions{
    Recursive:    true,
    WithVersions: true,
})

opt := ossClient.RemoveObjectsOptions{
    GovernanceBypass: true,
}
chErr := client.RemoveObjects(context.Background(), bucketName, objectCh, opt)
if chErr != nil {
    for err := range chErr {
        fmt.Println(err)
    }
    return
}
```

__示例2:删除一个桶中的所有对象的最新版本__


```go
//删除一个桶中的所有对象的最新版本
objectCh := client.ListObjects(context.Background(), bucketName, ossClient.ListObjectsOptions{
    Recursive: true,
})

opt := ossClient.RemoveObjectsOptions{
    GovernanceBypass: true,
}
chErr := client.RemoveObjects(context.Background(), bucketName, objectCh, opt)
if chErr != nil {
    for err := range chErr {
        fmt.Println(err)
    }
    return
}
```

