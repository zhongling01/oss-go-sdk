### RecycleBucket

### (ctx context.Context, bucketName string) 

### (bool, error)

将桶从回收站中恢复

__参数__

| 参数          | 类型              | 描述       |
| ------------- | ----------------- | ---------- |
| `ctx`         | _context.Context_ | 上下文控制 |
| `bucketName ` | _string_          | 存储桶名称 |



| 返回  | 类型  | 描述      |
| ----- | ----- | --------- |
| `err` | _err_ | 标准Error |



__示例__


```go
//初始化客户端以调用sdk
client, err := clientInit(AccessKeyIDDefault, SecretAccessKeyDefault)
if err != nil {
    return err
}
//创建一个桶名称为bucketName的值的桶并打开桶的回收站功能
_, err = BucketCreate(bucketName, false, true, false)
if err != nil {
    return err
}
//对桶进行删除操作
_, err = BucketDeleteWithOptions(bucketName, false)
if err != nil {
   return err
}

//将回收站的桶恢复
err = client.RecycleBucket(context.Background(), bucketName)
if err != nil {
    return err
}
```

