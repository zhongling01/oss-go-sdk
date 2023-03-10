### UpdateObject

### (updateOffset int, updateMod, bucketName, objectName string, reader io.Reader, objectSize int64)

### (info UploadInfo, err error)

对存储桶中已存在的对象进行局部更新

局部更新大小不得超过5GB

## 模式

一共有插入和替换两种类型，对应字符串"Insert" 和 "Replace"

- 举例：原数据为 12345，要更新的数据为678

  - insert模式，offset为0，更新后数据为67812345

  - insert模式，offset为1，更新后数据为16782345

  - insert模式，offset为5，更新后数据为12345678

  - insert模式，offset为-1，更新后数据为12345678

  - insert模式，offset为6，返回错误

  - replace模式，offset为0，更新后数据为67845

  - replace模式，offset为1，更新后数据为16785

  - replace模式，offset为4，更新后数据为1234678

  - replace模式，offset为5，更新后数据为12345678

  - replace模式，offset为-1，更新后数据为12345678

  - replace模式，offset为6，返回错误

    


__参数__

| 参数            | 类型              | 描述                                        |
| --------------- | ----------------- | ------------------------------------------- |
| `ctx`           | _context.Context_ | 上下文控制                                  |
| `updateOffset ` | _int_             | 更新的偏移量                                |
| `updateMod`     | _string_          | 更新的模式，对应字符串"Insert" 和 "Replace" |
| `bucketName`    | _string_          | 存储桶名称                                  |
| `reader`        | _io.Reader_       | 任意实现了io.Reader的GO类型                 |
| `objectSize`    | _int64_           | 更新对象的大小，需要填写                    |



**返回**

| UploadInfo  | 类型     | 描述                     |
| ----------- | -------- | ------------------------ |
| `Bucket`    | _string_ | 更新目标所在的存储桶名称 |
| `Key`       | _string_ | 更新的对象名称           |
| `ETag`      | _string_ | 更新的对象的ETag         |
| `Size`      | _int64_  | 更新对象的大小           |
| `VersionID` | _string_ | 更新对象的新版本号       |

__示例：__


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

//准备更新文件
updateMode := PartialUpdateInsertMode

reader, err := os.Open(updatefilePath)
if err != nil {
    fmt.Printf("Open file Error: %v\n", err)
    return
}
defer reader.Close()

fileStat, err := reader.Stat()
if err != nil {
    fmt.Printf("get file stat Error:%v\n", err)
    return
}

objectSize := fileStat.Size()

//更新对象
_, err := client.UpdateObject(2, updateMode, bucketName, objectName, reader, objectSize)
if err != nil {
    fmt.Printf("Update Error:%v\n", err)
    return
}

fmt.Printf("update file success length:[%d] mode:[%s]\n", objectSize, updateMode)
```

