### ExtractOnline

### (ctx context.Context, bucketName string, reader io.Reader, objectSize int64)

### (info UploadInfo, err error)

上传一个对象并将其在线解压缩

注意：

1. 在使用在线解压缩上传时不会启用分段上传
2. 目前支持的在线解压缩时上传对象大小不超过5GB
3. 支持在线加压缩的压缩格式类型有

 		Gzip S2 Zstd BZ2 LZ4  

4. 其他格式的文件在启用在线解压缩功能时会报错：

​		tar file error: archive/tar: invalid tar header

__参数__

| 参数         | 类型              | 描述                        |
| ------------ | ----------------- | --------------------------- |
| `ctx`        | _context.Context_ | 上下文控制                  |
| `bucketName` | _string_          | 存储桶名称                  |
| `reader`     | _io.Reader_       | 任意实现了io.Reader的GO类型 |
| `objectSize` | _int64_           | 上传的对象的大小            |



**返回**

| UploadInfo | 类型     | 描述             |
| ---------- | -------- | ---------------- |
| `Bucket`   | _string_ | 上传的存储桶名称 |
| `ETag`     | _string_ | 上传的对象的ETag |
| `Size`     | _int64_  | 上传对象的总大小 |

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

// 打开相关文件
fileReader, err := os.Open("123.tar")
if err != nil {
    log.Fatal(err)
}
defer fileReader.Close()

// 获取文件数据
fileStat, err := fileReader.Stat()
if err != nil {
    log.Fatal(err)
}

// 获取文件大小
fileSize := fileStat.Size()

//在线解压缩上传
_, err := client.ExtractOnline(context.Background(), bucketName, fileReader, fileSize)
if err != nil {
    log.Fatal(err)
}
//列出在线解压缩后的内容
ch := client.ListObjects(context.Background(), bucketName, ossClient.ListObjectsOptions{
		Recursive: true,
	})
	for i := range ch {
		fmt.Println(i.Key)
}

```

样例中123.tar解压后的文件组织为

123/
123/1.txt
123/2.txt
123/3/
123/3/3.txt

则上传至桶中的组织结构为

123/
123/1.txt
123/2.txt
123/3/
123/3/3.txt

