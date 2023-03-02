### MakeBucket

### (ctx context.Context,bucketName, opts MakeBucketOptions) 

### error

创建一个存储桶。

__参数__

| 参数         | 类型                          | 描述       |
| ------------ | ----------------------------- | ---------- |
| `ctx`        | _context.Context_             | 上下文控制 |
| `bucketName` | _string_                      | 存储桶名称 |
| `opts `      | _ossClient.MakeBucketOptions_ | 新建桶选项 |



| ossClient.MakeBucketOptions | 类型     | 描述                                                         |
| --------------------------- | -------- | ------------------------------------------------------------ |
| `Region`                    | _string_ | 存储桶被创建的region(地区)，默认是us-east-1(美国东一区)，下面列举的是其它合法的值。 |
| `ObjectLocking`             | _bool_   | 是否开启桶的对象锁定功能，桶的对象锁定功能只有在创建桶时可以配置，创建完成后不可修改 |
| `RecycleEnabled`            | _bool_   | 是否开启桶的回收站功能，桶的回收站功能只有在创建桶时可以配置，创建完成后不可修改 |



| Region可用值 |                |                |
| ------------ | -------------- | -------------- |
| us-east-1    | eu-central-1   | me-south-1     |
| us-east-2    | eu-north-1     | sa-east-1      |
| us-west-1    | ap-east-1      | us-gov-west-1  |
| us-west-2    | ap-south-1     | us-gov-east-1  |
| ca-central-1 | ap-southeast-1 | cn-north-1     |
| eu-west-1    | ap-southeast-2 | cn-northwest-1 |
| eu-west-2    | ap-northeast-1 |                |
| eu-west-3    | ap-northeast-2 |                |
|              | ap-northeast-3 |                |

__示例：创建一个开启对象锁定和回收站功能的桶__


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

//创建桶
err = client.MakeBucket(context.Background(), "examplebucket", ossClient.MakeBucketOptions{
    Region:         "",   //默认为'us-east-1'
    ObjectLocking:  true, //启用对象锁定
    RecycleEnabled: true, //启用回收桶功能
})
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println("make bucket successful")
```

