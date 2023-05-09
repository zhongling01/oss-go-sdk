package ossClient

import (
	"context"
	"fmt"
	"github.com/trinet2005/oss-go-sdk/pkg/credentials"
	"strings"
	"testing"
	"time"
)

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
