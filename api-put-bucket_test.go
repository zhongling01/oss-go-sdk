package ossClient

import (
	"context"
	"github.com/zhongling01/oss-go-sdk/pkg/credentials"
	"testing"
)

func TestClient_MakeBucketPublicAccess(t *testing.T) {
	c, err := New(EndpointDefault, &Options{
		Creds: credentials.NewStaticV4(AccessKeyIDDefault, SecretAccessKeyDefault, ""),
	})
	if err != nil {
		t.Fatal(err)
	}
	bucketName := "test-bucket-public-access-option"

	err = c.MakeBucket(context.Background(), bucketName, MakeBucketOptions{
		PublicAccess: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.RemoveBucket(context.Background(), bucketName)
}

func TestClient_MakeBucketForceCreate(t *testing.T) {
	c, err := New(EndpointDefault, &Options{
		Creds: credentials.NewStaticV4(AccessKeyIDDefault, SecretAccessKeyDefault, ""),
	})
	if err != nil {
		t.Fatal(err)
	}
	bucketName := "test-bucket-force-create-option"

	err = c.MakeBucket(context.Background(), bucketName, MakeBucketOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = c.MakeBucket(context.Background(), bucketName, MakeBucketOptions{
		ForceCreate: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.RemoveBucket(context.Background(), bucketName)
}

func TestClient_MakeBucketRecycleEnabled(t *testing.T) {
	c, err := New(EndpointDefault, &Options{
		Creds: credentials.NewStaticV4(AccessKeyIDDefault, SecretAccessKeyDefault, ""),
	})
	if err != nil {
		t.Fatal(err)
	}
	bucketName := "test-bucket-recycle-option"

	err = c.MakeBucket(context.Background(), bucketName, MakeBucketOptions{
		RecycleEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.RemoveBucket(context.Background(), bucketName)
	defer c.RemoveBucket(context.Background(), bucketName)
}
