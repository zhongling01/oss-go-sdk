package ossClient

import (
	"context"
	"github.com/trinet2005/oss-go-sdk/pkg/credentials"
	"testing"
)

func TestClient_MakeBucket(t *testing.T) {
	c, err := New("127.0.0.1:19000", &Options{
		Creds: credentials.NewStaticV4("minioadmin", "minioadmin", ""),
	})
	if err != nil {
		t.Fatal(err)
	}
	bucketName := "test-bucket-public-access-option"

	err = c.makeBucket(context.Background(), bucketName, MakeBucketOptions{
		PublicAccess: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	//defer c.RemoveBucket(context.Background(), bucketName)
}
