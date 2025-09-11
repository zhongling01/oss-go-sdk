package ossClient

import (
	"context"
	"fmt"
	"github.com/zhongling01/oss-go-sdk/pkg/credentials"
	"strings"
	"testing"
	"time"
)

// 基础测试
func TestClient_BaseRemoveObject(t *testing.T) {
	opts := &Options{
		Creds: credentials.NewStaticV4(AccessKeyIDDefault, SecretAccessKeyDefault, ""),
	}
	c, err := New(EndpointDefault, opts)
	if err != nil {
		t.Fatal(err.Error())
	}

	bucket := "test-remove-bucket"
	err = c.MakeBucket(context.Background(), bucket, MakeBucketOptions{ForceCreate: true})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer c.RemoveBucketWithOptions(context.Background(), bucket, RemoveBucketOptions{
		ForceDelete: true,
	})

	smallData := "0123456789"
	smallDataSize := int64(len(smallData))
	bigData := strings.Repeat(smallData, 1024*1024)
	bigDataSize := int64(len(bigData))

	// 测试基本的删除
	obj := "test"
	_, err = c.PutObject(context.Background(), bucket, obj, strings.NewReader(smallData), smallDataSize, PutObjectOptions{})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = c.RemoveObject(context.Background(), bucket, obj, RemoveObjectOptions{})
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = c.PutObject(context.Background(), bucket, obj, strings.NewReader(bigData), bigDataSize, PutObjectOptions{})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = c.RemoveObject(context.Background(), bucket, obj, RemoveObjectOptions{})
	if err != nil {
		t.Fatal(err.Error())
	}

	// 测试子目录
	obj = "subdir/obj"
	_, err = c.PutObject(context.Background(), bucket, obj, strings.NewReader(smallData), smallDataSize, PutObjectOptions{})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = c.RemoveObject(context.Background(), bucket, obj, RemoveObjectOptions{})
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = c.PutObject(context.Background(), bucket, obj, strings.NewReader(smallData), smallDataSize, PutObjectOptions{})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = c.RemoveObject(context.Background(), bucket, "subdir", RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		t.Fatal(err.Error())
	}
}

// 测试 ParallelDrives
func TestClient_RemoveBucketWithOptions(t *testing.T) {
	opts := &Options{
		Creds: credentials.NewStaticV4(AccessKeyIDDefault, SecretAccessKeyDefault, ""),
	}
	// 当删除数目过多使用强制删除会很耗时，因此要提高超时时间来等待
	transport, err := DefaultTransport(opts.Secure)
	if err != nil {
		t.Fatal(err.Error())
	}
	transport.ResponseHeaderTimeout = 10 * time.Minute
	opts.Transport = transport
	c, err := New(EndpointDefault, opts)
	if err != nil {
		t.Fatal(err.Error())
	}

	bucket := "test-remove-bucket"
	err = c.MakeBucket(context.Background(), bucket, MakeBucketOptions{ForceCreate: true})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = c.RemoveBucketWithOptions(context.Background(), bucket, RemoveBucketOptions{
		ForceDelete: true,
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	err = c.MakeBucket(context.Background(), bucket, MakeBucketOptions{ForceCreate: true})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = c.RemoveBucketWithOptions(context.Background(), bucket, RemoveBucketOptions{
		ForceDelete: true,
		Internal: AdvancedDeleteBucketOptions{
			ParallelDrives: 4,
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}

}

// 测试 DeletePrefixParallelDrives
func TestClient_RemoveObject(t *testing.T) {
	opts := &Options{
		Creds: credentials.NewStaticV4(AccessKeyIDDefault, SecretAccessKeyDefault, ""),
	}
	// 当删除数目过多使用强制删除会很耗时，因此要提高超时时间来等待
	transport, err := DefaultTransport(opts.Secure)
	if err != nil {
		t.Fatal(err.Error())
	}
	transport.ResponseHeaderTimeout = 10 * time.Minute
	opts.Transport = transport
	c, err := New(EndpointDefault, opts)
	if err != nil {
		t.Fatal(err.Error())
	}

	bucket := "test-remove-bucket"
	err = c.MakeBucket(context.Background(), bucket, MakeBucketOptions{ForceCreate: true})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer c.RemoveBucketWithOptions(context.Background(), bucket, RemoveBucketOptions{
		ForceDelete: true,
	})

	for i := 0; i < 10; i++ {
		_, err = c.PutObject(context.Background(), bucket, fmt.Sprintf("subdir/obj%d", i), strings.NewReader("1234"), int64(4), PutObjectOptions{})
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	err = c.RemoveObject(context.Background(), bucket, "subdir", RemoveObjectOptions{
		ForceDelete: true,
		Internal: AdvancedRemoveOptions{
			DeletePrefixParallelDrives: 4,
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}
