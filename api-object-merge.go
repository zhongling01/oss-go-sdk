package ossClient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"
	"fmt"
	uuid2 "github.com/google/uuid"
	"io"
)

const (
	MergeDir   = ".internal.merge.objects/"
	IdxPrefix  = ".idx."
	DataPrefix = ".data."
)

type ObjectIndex struct {
	Valid  bool  `json:"valid"`
	Offset int64 `json:"offset"`
	Size   int64 `json:"size"`
}

type ObjectIndexInfo struct {
	VacancySize int64                   `json:"vacancySize"`
	TotalSize   int64                   `json:"totalSize"`
	ObjectNum   int                     `json:"objectNum"`
	Info        map[string]*ObjectIndex `json:"objInfo"`
}

type PutObjectMerge struct {
	ID         string
	bucketName string
	client     *Client
	meta       *ObjectIndexInfo
	reader     *PutObjectMergeReader
}

type PutObjectMergeReader struct {
	data [][]byte
	i, j int64
}

func (r *PutObjectMergeReader) Read(buf []byte) (int, error) {
	if r.i >= int64(len(r.data)) {
		return 0, io.EOF
	}

	if r.j >= int64(len(r.data[r.i])) {
		r.i++
		r.j = 0
		return 0, nil
	}

	n := copy(buf, r.data[r.i][r.j:])
	r.j += int64(n)
	return n, nil
}

func checkBucket(c *Client, bucketName string) error {
	info, err := c.GetBucketVersioning(context.Background(), bucketName)
	if err != nil {
		return err
	}
	if info.Enabled() {
		return errors.New("Merge operations are not supported for buckets with versioning turned on")
	}
	//if info.Suspended() {
	//	return errors.New("Bucket is suspended")
	//}

	return nil
}

func (c *Client) InitMergePartUpload(id, bucketName string) (*PutObjectMerge, error) {
	err := checkBucket(c, bucketName)
	if err != nil {
		return nil, err
	}

	if id == "" {
		uuid, err := uuid2.NewRandom()
		if err != nil {
			return nil, err
		}
		id = fmt.Sprintf("%s-%d",uuid.String(), time.Now().UnixNano())
	}

	return &PutObjectMerge{
		ID:         id,
		bucketName: bucketName,
		client:     c,
		meta: &ObjectIndexInfo{
			VacancySize: 0,
			TotalSize:   0,
			ObjectNum:   0,
			Info:        make(map[string]*ObjectIndex, 0),
		},
		reader: &PutObjectMergeReader{
			data: make([][]byte, 0),
		},
	}, nil
}

func (p *PutObjectMerge) UploadMergePart(objectName string, reader io.Reader) (*ObjectIndexInfo, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	dataSize := int64(len(data))
	if dataSize == 0 {
		return nil, errors.New("no data given")
	}

	if p.meta.TotalSize + dataSize > maxSinglePutObjectSize{
		return nil, errors.New("the merged file is too large, please execute CompleteMergePartUpload first")
	}

	p.reader.data = append(p.reader.data, data)
	if _, ok := p.meta.Info[objectName]; ok {
		p.meta.VacancySize += p.meta.Info[objectName].Size
	} else {
		p.meta.ObjectNum++
	}
	p.meta.Info[objectName] = &ObjectIndex{
		Valid:  true,
		Offset: p.meta.TotalSize,
		Size:   dataSize,
	}
	p.meta.TotalSize += dataSize

	return p.meta, nil
}

func (p *PutObjectMerge) CompleteMergePartUpload(ctx context.Context) error {
	err := checkBucket(p.client, p.bucketName)
	if err != nil {
		return err
	}

	objectIndexInfo, err := json.Marshal(p.meta)
	if err != nil {
		return err
	}

	_, err = p.client.PutObject(ctx, p.bucketName, MergeDir+IdxPrefix+p.ID, bytes.NewReader(objectIndexInfo), int64(len(objectIndexInfo)), PutObjectOptions{})
	if err != nil {
		return err
	}

	_, err = p.client.PutObject(ctx, p.bucketName, MergeDir+DataPrefix+p.ID, p.reader, p.meta.TotalSize, PutObjectOptions{})
	if err != nil {
		p.client.removeObject(ctx, p.bucketName, MergeDir+IdxPrefix+p.ID, RemoveObjectOptions{GovernanceBypass: true})
		return err
	}

	return nil
}

func (c *Client) GetObjectWithID(ctx context.Context, id, bucketName, objectName string) (*Object, *ObjectIndexInfo, error) {
	meta, err := c.GetObjectIndexInfo(ctx, id, bucketName)
	if err != nil {
		return nil, nil, err
	}

	data, err := c.GetObjectWithIndex(ctx, id, bucketName, objectName, meta)
	if err != nil {
		return nil, nil, err
	}

	return data, meta, nil
}

func (c *Client) GetObjectIndexInfo(ctx context.Context, id, bucketName string) (*ObjectIndexInfo, error) {
	meta := &ObjectIndexInfo{
		Info: make(map[string]*ObjectIndex, 0),
	}

	metaData, err := c.GetObject(ctx, bucketName, MergeDir+IdxPrefix+id, GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer metaData.Close()

	buf, err := io.ReadAll(metaData)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, meta)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (c *Client) GetObjectWithIndex(ctx context.Context, id, bucketName, objectName string, meta *ObjectIndexInfo) (*Object, error) {
	if _, ok := meta.Info[objectName]; !ok {
		return nil, errors.New("object not found")
	} else if !meta.Info[objectName].Valid {
		return nil, errors.New("object invalid")
	}

	opts := GetObjectOptions{}
	err := opts.SetRange(meta.Info[objectName].Offset, meta.Info[objectName].Offset+meta.Info[objectName].Size-1)
	if err != nil {
		return nil, err
	}

	data, err := c.GetObject(ctx, bucketName, MergeDir+DataPrefix+id, opts)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) DeleteMergeID(ctx context.Context, id, bucketName string) error {
	err := checkBucket(c, bucketName)
	if err != nil {
		return err
	}

	err = c.RemoveObject(ctx, bucketName, MergeDir+DataPrefix+id, RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return c.RemoveObject(ctx, bucketName, MergeDir+IdxPrefix+id, RemoveObjectOptions{})
}

func (c *Client) DeleteObjectWithId(ctx context.Context, id, bucketName, objectName string) error {
	err := checkBucket(c, bucketName)
	if err != nil {
		return err
	}

	meta, err := c.GetObjectIndexInfo(ctx, id, bucketName)
	if err != nil {
		return err
	}

	if _, ok := meta.Info[objectName]; !ok {
		return errors.New("object not found")
	} else if !meta.Info[objectName].Valid {
		return errors.New("object already invalid")
	}

	meta.Info[objectName].Valid = false
	meta.VacancySize += meta.Info[objectName].Size

	objectIndexInfo, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	_, err = c.PutObject(ctx, bucketName, MergeDir+IdxPrefix+id, bytes.NewReader(objectIndexInfo), int64(len(objectIndexInfo)), PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
