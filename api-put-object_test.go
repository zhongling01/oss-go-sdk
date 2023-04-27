/*
 * MinIO Go Library for Amazon S3 Compatible Cloud Storage
 * Copyright 2015-2017 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package ossClient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/trinet2005/oss-go-sdk/pkg/credentials"
	"io"
	"strings"
	"testing"
)

func TestPutObjectOptionsValidate(t *testing.T) {
	testCases := []struct {
		name, value string
		shouldPass  bool
	}{
		// Invalid cases.
		{"X-Amz-Matdesc", "blah", false},
		{"x-amz-meta-X-Amz-Iv", "blah", false},
		{"x-amz-meta-X-Amz-Key", "blah", false},
		{"x-amz-meta-X-Amz-Matdesc", "blah", false},
		{"It has spaces", "v", false},
		{"It,has@illegal=characters", "v", false},
		{"X-Amz-Iv", "blah", false},
		{"X-Amz-Key", "blah", false},
		{"X-Amz-Key-prefixed-header", "blah", false},
		{"Content-Type", "custom/content-type", false},
		{"content-type", "custom/content-type", false},
		{"Content-Encoding", "gzip", false},
		{"Cache-Control", "blah", false},
		{"Content-Disposition", "something", false},
		{"Content-Language", "somelanguage", false},

		// Valid metadata names.
		{"my-custom-header", "blah", true},
		{"custom-X-Amz-Key-middle", "blah", true},
		{"my-custom-header-X-Amz-Key", "blah", true},
		{"blah-X-Amz-Matdesc", "blah", true},
		{"X-Amz-MatDesc-suffix", "blah", true},
		{"It-Is-Fine", "v", true},
		{"Numbers-098987987-Should-Work", "v", true},
		{"Crazy-!#$%&'*+-.^_`|~-Should-193832-Be-Fine", "v", true},
	}
	for i, testCase := range testCases {
		err := PutObjectOptions{UserMetadata: map[string]string{
			testCase.name: testCase.value,
		}}.validate()
		if testCase.shouldPass && err != nil {
			t.Errorf("Test %d - output did not match with reference results, %s", i+1, err)
		}
	}
}

/* trinet */
func testPartialUpdate(originData []byte, mode string, offset int64, newData io.Reader, originSize, bodySize int64, expect string) error {
	opts := &Options{
		Creds: credentials.NewStaticV4(AccessKeyIDDefault, SecretAccessKeyDefault, ""),
	}
	client, err := New(EndpointDefault, opts)
	if err != nil {
		return err
	}
	bucket := "test-bucket"
	object := "test-partial-obj"
	err = client.MakeBucket(context.Background(), bucket, MakeBucketOptions{ForceCreate: true})
	if err != nil {
		return err
	}
	defer client.RemoveBucketWithOptions(context.Background(), bucket, RemoveBucketOptions{ForceDelete: true})

	// 上传一个初始的对象
	_, err = client.PutObject(context.Background(), bucket, object, bytes.NewReader(originData), originSize, PutObjectOptions{})
	if err != nil {
		return err
	}
	defer client.RemoveObject(context.Background(), bucket, object, RemoveObjectOptions{})

	// 验证局部更新
	_, err = client.UpdateObject(context.Background(), bucket, object, mode, int(offset), newData, bodySize)
	if err != nil {
		return err
	}
	gr, err := client.GetObject(context.Background(), bucket, object, GetObjectOptions{})

	data, err := io.ReadAll(gr)
	if err != nil {
		return err
	}

	//println(expect)
	if string(data) != expect {
		return errors.New(fmt.Sprintf("expect: %s, but get:%s\n", expect, string(data)))
	}

	return nil
}

// 测试局部更新Insert模式
func TestPartialUpdateInsert(t *testing.T) {
	var offset, size int64

	origin := "12345"
	newData := "678"

	originData := []byte(origin)
	originSize := int64(len(originData))
	size = int64(len(newData))

	offset = 0
	expect := origin[:offset] + newData + origin[offset:]
	err := testPartialUpdate(originData, PartialUpdateInsertMode, offset, bytes.NewReader([]byte(newData)), originSize, size, expect)
	if err != nil {
		t.Fatal(err)
	}

	offset = 1
	expect = origin[:offset] + newData + origin[offset:]
	err = testPartialUpdate(originData, PartialUpdateInsertMode, offset, bytes.NewReader([]byte(newData)), originSize, size, expect)
	if err != nil {
		t.Fatal(err)
	}

	offset = originSize
	expect = origin[:offset] + newData + origin[offset:]
	err = testPartialUpdate(originData, PartialUpdateInsertMode, offset, bytes.NewReader([]byte(newData)), originSize, size, expect)
	if err != nil {
		t.Fatal(err)
	}

	offset = originSize + 1
	expect = "test error case"
	err = testPartialUpdate(originData, PartialUpdateInsertMode, offset, bytes.NewReader([]byte(newData)), originSize, size, expect)
	if err == nil {
		t.Fatal("want error")
	} else {
		t.Log(err)
	}
}

// 测试局部更新Replace模式
func TestPartialUpdateReplace(t *testing.T) {
	var offset, size int64
	var expect string

	origin := "12345"
	newData := "678"

	originData := []byte(origin)
	originSize := int64(len(originData))
	size = int64(len(newData))

	offset = 0
	if offset+size < originSize {
		expect = origin[:offset] + newData + origin[offset+size:]
	} else {
		expect = origin[:offset] + newData
	}
	err := testPartialUpdate(originData, PartialUpdateReplaceMode, offset, bytes.NewReader([]byte(newData)), originSize, size, expect)
	if err != nil {
		t.Fatal(err)
	}

	offset = 1
	if offset+size < originSize {
		expect = origin[:offset] + newData + origin[offset+size:]
	} else {
		expect = origin[:offset] + newData
	}
	err = testPartialUpdate(originData, PartialUpdateReplaceMode, offset, bytes.NewReader([]byte(newData)), originSize, size, expect)
	if err != nil {
		t.Fatal(err)
	}

	offset = originSize
	if offset+size < originSize {
		expect = origin[:offset] + newData + origin[offset+size:]
	} else {
		expect = origin[:offset] + newData
	}
	err = testPartialUpdate(originData, PartialUpdateReplaceMode, offset, bytes.NewReader([]byte(newData)), originSize, size, expect)
	if err != nil {
		t.Fatal(err)
	}

	offset = originSize + 1
	expect = "test error case"
	err = testPartialUpdate(originData, PartialUpdateReplaceMode, offset, bytes.NewReader([]byte(newData)), originSize, size, expect)
	if err == nil {
		t.Fatal("want error")
	} else {
		t.Log(err)
	}
}

func testAppend(originData []byte, newData io.Reader, originSize, bodySize int64, expect string) error {
	opts := &Options{
		Creds: credentials.NewStaticV4(AccessKeyIDDefault, SecretAccessKeyDefault, ""),
	}
	client, err := New(EndpointDefault, opts)
	if err != nil {
		return err
	}
	bucket := "test-bucket"
	object := "test-append-obj"
	err = client.MakeBucket(context.Background(), bucket, MakeBucketOptions{ForceCreate: true})
	if err != nil {
		return err
	}
	defer client.RemoveBucketWithOptions(context.Background(), bucket, RemoveBucketOptions{ForceDelete: true})

	// 上传一个初始的对象
	_, err = client.PutObject(context.Background(), bucket, object, bytes.NewReader(originData), originSize, PutObjectOptions{})
	if err != nil {
		return err
	}
	defer client.RemoveObject(context.Background(), bucket, object, RemoveObjectOptions{})

	// 验证局部更新
	_, err = client.AppendObject(context.Background(), bucket, object, newData, bodySize)
	if err != nil {
		return err
	}
	gr, err := client.GetObject(context.Background(), bucket, object, GetObjectOptions{})

	data, err := io.ReadAll(gr)
	if err != nil {
		return err
	}

	//println(expect)
	if string(data) != expect {
		return errors.New(fmt.Sprintf("expect: %s, but get:%s\n", expect, string(data)))
	}

	return nil
}

// 测试追加
func TestAppendObject(t *testing.T) {
	var size int64

	origin := "12345"
	newData := "678"

	originData := []byte(origin)
	originSize := int64(len(originData))
	size = int64(len(newData))

	expect := origin[:] + newData
	err := testAppend(originData, bytes.NewReader([]byte(newData)), originSize, size, expect)
	if err != nil {
		t.Fatal(err)
	}
}

// 测试写入指定存储引擎池
func TestPreferredEnginePool(t *testing.T) {
	opts := &Options{
		Creds: credentials.NewStaticV4(AccessKeyIDDefault, SecretAccessKeyDefault, ""),
	}
	client, err := New(EndpointDefault, opts)
	if err != nil {
		t.Fatal(err.Error())
	}

	//  ====== 测试基础的引擎 ======
	bucket := "test-pool-engine-bucket"
	object := "test-obj"
	err = client.MakeBucket(context.Background(), bucket, MakeBucketOptions{ForceCreate: true})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer client.RemoveBucketWithOptions(context.Background(), bucket, RemoveBucketOptions{ForceDelete: true})

	data := "test"
	size := int64(len(data))
	for _, engine := range []ErasurePoolEngine{DefaultEngine, HDD, SSD} {
		// 使用debug，去服务端看是否正确写入存储池
		_, err = client.PutObject(context.Background(), bucket, object, strings.NewReader(data), size,
			PutObjectOptions{PreferredEnginePool: engine})
		if err != nil {
			t.Fatal(err.Error())
		}
		err = client.RemoveObject(context.Background(), bucket, object, RemoveObjectOptions{})
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	_, err = client.PutObject(context.Background(), bucket, object, strings.NewReader(data), size,
		PutObjectOptions{PreferredEnginePool: HDD})
	if err != nil {
		t.Fatal(err.Error())
	}

	// ====== 测试 CopyObject ======
	src := CopySrcOptions{
		Bucket: bucket,
		Object: object,
	}
	dstBucket := "test-pool-engine-bucket-dst"
	err = client.MakeBucket(context.Background(), dstBucket, MakeBucketOptions{ForceCreate: true})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer client.RemoveBucketWithOptions(context.Background(), dstBucket, RemoveBucketOptions{ForceDelete: true})
	dst := CopyDestOptions{
		Bucket: dstBucket,
		Object: object,
		Size:   size,
	}
	_, err = client.CopyObject(context.Background(), dst, src)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = client.RemoveObject(context.Background(), dstBucket, object, RemoveObjectOptions{})
	if err != nil {
		t.Fatal(err.Error())
	}

	// ====== 测试multipart模式 ======
	// 见 TestClient_MultipartUploadPreferredEnginePool

	// TODO: 测试decommissioning
}

/* trinet */
