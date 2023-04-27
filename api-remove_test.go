package ossClient

import (
	"context"
	"github.com/trinet2005/oss-go-sdk/pkg/credentials"
	"testing"
)

func TestClient_RemoveBucketWithOptions(t *testing.T) {
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
