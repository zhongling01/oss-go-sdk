### SelectObjectContent(ctx context.Context, bucketName string, objectsName string, expression string, options SelectObjectOptions) (*SelectResults,error)

Select 功能通过结构化查询语句（SQL）筛选存储在对象存储上的对象，以便检索对象并获取用户所需的数据。通过Select 功能筛选对象数据，您可以减少OSS传输的数据量，这将降低检索此数据所需的成本和延迟。

Select 功能目前支持检索以 CSV和JSON 格式存储的对象，支持检索通过 GZIP 或 BZIP2 压缩的对象（仅对于 CSV、JSON 格式的对象）。此外，Select 功能还支持将结果的格式指定为 CSV 或 JSON，并且可以确定结果中记录的分隔方式。

| 字段         | 类型                  | 描述                     |
| :----------- | :-------------------- | :----------------------- |
| `ctx`        | _context.Context_     | 上下文信息               |
| `bucketName` | _string_              | 桶名称                   |
| `objectName` | _string_              | 对象名称                 |
| `options`    | _SelectObjectOptions_ | select对象的一些可选选项 |

__SelectObjectOptions__

| 字段                  | 类型                            | 描述                          |
| :-------------------- | :------------------------------ | :---------------------------- |
| `Expression`          | string                          | SQL语句                       |
| `ExpressionType`      | QueryExpressionType             | 表达式类型，固定为字符串"SQL" |
| `InputSerialization`  | SelectObjectInputSerialization  | 设置select目标对象的数据格式  |
| `OutputSerialization` | SelectObjectOutputSerialization | 设置输出结果的数据格式        |



__SelectObjectInputSerializationt 目标对象的数据格式__

| 字段              | 类型                 | 描述                                                         |
| :---------------- | :------------------- | :----------------------------------------------------------- |
| `CompressionType` | string               | 设置目标的压缩类型，如"GZIP"，"BZIP2"                        |
| `Parquet`         | *ParquetInputOptions | 如果目标对象为Parquet格式设置相应设置（服务解压Parquet不稳定，服务端默认关闭Parquet的解析） |
| `CSV`             | *CSVInputOptions     | 如果目标对象为CSV格式设置相应设置                            |
| `JSON`            | *JSONInputOptions    | 如果目标对象为JSON格式设置相应设置                           |

**CSVInputOptions**：

- **`FileHeaderInfo`**：一个字符串，表示CSV文件是否包含标题行，可选值为"USE"、"IGNORE"或"NONE"。
- **`RecordDelimiter`**：一个字符串，表示CSV文件中记录之间的分隔符。
- **`FieldDelimiter`**：一个字符串，表示CSV文件中字段之间的分隔符。
- **`QuoteCharacter`**：一个字符串，表示CSV文件中用于引用字段的字符。
- **`QuoteEscapeCharacter`**：一个字符串，表示CSV文件中用于转义引号字符的字符。
- **`Comments`**：一个字符串，表示CSV文件中用于注释的字符。

**JSONInputOptions**：

**`Type`**：一个字符串，表示JSON文件是文档格式还是行格式，可选值为"DOCUMENT"或"LINES"。



__SelectObjectOutputSerialization 输出结果的数据格式__

| 字段   | 类型               | 描述                               |
| :----- | :----------------- | :--------------------------------- |
| `CSV`  | *CSVOutputOptions  | 如果文件输出为CSV格式设置相应设置  |
| `JSON` | *JSONOutputOptions | 如果文件输出为JSON格式设置相应设置 |

__CSVOutputOptions__:

- **`QuoteFields`**：一个字符串，表示CSV文件引用风格，可选值为"Always"或"AsNeeded"。
- **`RecordDelimiter`**：一个字符串，表示CSV文件中记录之间的分隔符。
- **`FieldDelimiter`**：一个字符串，表示CSV文件中字段之间的分隔符。
- **`QuoteCharacter`**：一个字符串，表示CSV文件中用于引用字段的字符。
- **`QuoteEscapeCharacter`**：一个字符串，表示CSV文件中用于转义引号字符的字符。

__JSONOutputOptions__:

**`RecordDelimiter`**：一个字符串，表示JSON文件中记录之间的分隔符。



__返回值__

| 字段            | 类型            | 描述                                                         |
| :-------------- | :-------------- | :----------------------------------------------------------- |
| `SelectResults` | _SelectResults_ | 一个可以直接传递给csv.NewReader进行处理输出的io.ReadCloser对象。 |

__样例__

向桶中上传一个json文件 `line.JSON`

```json
{"name":"111","recycleEnable":true,"objectLocking":true,"versioning":true,"publicAccess":true}

{"name":"222","recycleEnable":true,"objectLocking":true,"versioning":true,"publicAccess":true}
```

运行代码示例结果：

```
{"name":"111"}?
{"name":"222"}?
```



```go
package main

import (
	"context"
	"fmt"
	ossClient "github.com/trinet2005/oss-go-sdk"
	"github.com/trinet2005/oss-go-sdk/pkg/credentials"
	"io"
	"os"
	"storConsole/app/oss/ossBase"
)

func main() {
	//初始化客户端
	accessKeyID := ""     //用户id
	secretAccessKey := "" //用户密码
	Endpoint := ""        //服务终端

	opts := &ossClient.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Region: ossBase.RegionDefault,
	}

	c, err := ossClient.New(Endpoint, opts)

	if err != nil {
		return
	}
	bucketName := "" //桶名称
	objectName := "" //对象名称

	reader, err := c.SelectObjectContent(context.Background(), bucketName, objectName, ossClient.SelectObjectOptions{
		Expression:     "select name from s3object",
		ExpressionType: ossClient.QueryExpressionTypeSQL,
		InputSerialization: ossClient.SelectObjectInputSerialization{
			JSON: &ossClient.JSONInputOptions{
				Type: ossClient.JSONLinesType,
			},
		},
		OutputSerialization: ossClient.SelectObjectOutputSerialization{
			JSON: &ossClient.JSONOutputOptions{
				RecordDelimiter: "?\n",
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	//新建文件
	newR, err := os.Create("tmp.json")
	if err != nil {
		return
	}
	defer os.Remove("tmp.json")
	defer newR.Close()
	// 读取SelectObject结果
	buf := make([]byte, 512)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}
		fmt.Println(string(buf[:n]))
		_, err = newR.Write(buf[:n])
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if err != nil {
		return
	}
	return
}

```

