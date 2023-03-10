### CopyObject

### (ctx context.Context, dst CopyDestOptions, src CopySrcOptions) 

### (UploadInfo, error)

通过现有对象的服务器端复制来创建或替换对象。它支持条件复制、复制对象的一部分和服务器端的目标加密和源解密。有关详细信息，请参阅`CopySrcOptions`和类型。`DestinationInfo`

要将多个源对象复制到单个目标对象，请参阅`ComposeObject`API。

**参数**

| 参数  | 类型                    | 描述                        |
| :---- | :---------------------- | :-------------------------- |
| `ctx` | *context.上下文*        | 呼叫超时/取消的自定义上下文 |
| `dst` | *minio.CopyDestOptions* | 描述目标对象的参数          |
| `src` | *minio.CopySrc选项*     | 描述源对象的参数            |

**minio.UploadInfo**

| 场地             | 类型   | 描述               |
| :--------------- | :----- | :----------------- |
| `info.ETag`      | *细绳* | 新对象的 ETag      |
| `info.VersionID` | *细绳* | 新对象的版本标识符 |

**例子**

```go
// Use-case 1: Simple copy object with no conditions.
// Source object
srcOpts := minio.CopySrcOptions{
    Bucket: "my-sourcebucketname",
    Object: "my-sourceobjectname",
}

// Destination object
dstOpts := minio.CopyDestOptions{
    Bucket: "my-bucketname",
    Object: "my-objectname",
}

// Copy object call
uploadInfo, err := minioClient.CopyObject(context.Background(), dst, src)
if err != nil {
    fmt.Println(err)
    return
}

fmt.Println("Successfully copied object:", uploadInfo)
// Use-case 2:
// Copy object with copy-conditions, and copying only part of the source object.
// 1. that matches a given ETag
// 2. and modified after 1st April 2014
// 3. but unmodified since 23rd April 2014
// 4. copy only first 1MiB of object.

// Source object
srcOpts := minio.CopySrcOptions{
    Bucket: "my-sourcebucketname",
    Object: "my-sourceobjectname",
    MatchETag: "31624deb84149d2f8ef9c385918b653a",
    MatchModifiedSince: time.Date(2014, time.April, 1, 0, 0, 0, 0, time.UTC),
    MatchUnmodifiedSince: time.Date(2014, time.April, 23, 0, 0, 0, 0, time.UTC),
    Start: 0,
    End: 1024*1024-1,
}


// Destination object
dstOpts := minio.CopyDestOptions{
    Bucket: "my-bucketname",
    Object: "my-objectname",
}

// Copy object call
_, err = minioClient.CopyObject(context.Background(), dst, src)
if err != nil {
    fmt.Println(err)
    return
}

fmt.Println("Successfully copied object:", uploadInfo)
```