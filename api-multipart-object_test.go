package ossClient

import (
	"context"
	"fmt"
	uuid2 "github.com/google/uuid"
	"github.com/trinet2005/oss-go-sdk/pkg/credentials"
	"io"
	"strings"
	"testing"
)

type TestMultipart struct {
	size  int
	point int
	data  string
}

func (t *TestMultipart) Read(p []byte) (n int, err error) {
	if t.point > t.size {
		return 0, io.EOF
	}
	n = copy(p, t.data)
	t.point += n
	return
}

func TestClient_MultipartUpload(t *testing.T) {
	c, err := New("127.0.0.1:19000", &Options{
		Creds: credentials.NewStaticV4("minioadmin", "minioadmin", ""),
	})
	if err != nil {
		t.Fatal(err)
	}
	bucketName := "test-bucket" + uuid2.New().String()
	objectName := "test-object"
	c.makeBucket(context.Background(), bucketName, MakeBucketOptions{})
	defer c.RemoveBucket(context.Background(), bucketName)
	defer c.removeObject(context.Background(), bucketName, objectName, RemoveObjectOptions{})

	reader := &TestMultipart{
		size: 1024 * 1024 * 16,
		data: "0123456789",
	}

	opts := &PutObjectOptions{
		DisableMultipart: false,
		MergeMultipart:   true,
		PartSize:         absMinPartSize,
	}
	// Initiate a new multipart upload.
	// TODO：实现断点续传
	m, err := c.NewUploadID(context.Background(), bucketName, objectName, opts)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err != nil {
			m.AbortMultipartUpload(context.Background())
		}
	}()

	// Part number always starts with '1'
	partNumber := 1
	for partNumber <= maxPartsCount {
		reader.data = strings.Repeat(fmt.Sprintf("%d", partNumber), 10)
		uerr := m.UploadPart(context.Background(), reader, partNumber)
		if uerr != nil {
			if uerr == io.EOF {
				break
			}
			t.Fatal(uerr)
		}
		go func(partNumber int) {
			// 读取已上传分段
			r, _, err := m.GetPart(context.Background(), partNumber)
			if err != nil {
				t.Fatal(err)
			}
			tmpBuf, err := io.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(tmpBuf[:20]))
			r.Close()
		}(partNumber)

		partNumber++
	}

	// 更新已上传分片
	reader = &TestMultipart{
		size: 1024 * 1024 * 10,
		data: "abcdefg",
	}
	changePartNum := 2
	uerr := m.UpdatePart(context.Background(), reader, changePartNum, reader.size)
	if uerr != nil && uerr != io.EOF {
		t.Fatal(uerr)
	}
	// 读取更改的片段
	r, _, err := m.GetPart(context.Background(), changePartNum)
	if err != nil {
		t.Fatal(err)
	}
	tmpBuf, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(tmpBuf[:20]))
	r.Close()

	// 测试最后一个分段可以小于限制大小
	reader = &TestMultipart{
		size: 1024 * 1024 * 1,
		data: "QWERTYUIOP",
	}
	uerr = m.UpdatePart(context.Background(), reader, len(m.partsInfo), reader.size)
	if uerr != nil && uerr != io.EOF {
		t.Fatal(uerr)
	}
	// 读取更改的片段
	r, _, err = m.GetPart(context.Background(), len(m.partsInfo))
	if err != nil {
		t.Fatal(err)
	}
	tmpBuf, err = io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(tmpBuf[:20]))
	r.Close()

	uploadInfo, err := m.CompleteMultipartUpload(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(uploadInfo)

}
