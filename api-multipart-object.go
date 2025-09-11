/* trinet
 * support multipart sdk:
 * - support breakpoint transfer;
 * - support update the uploaded part;
 * - support for reading uploaded part;
 * - support for merging parts on completion.
 */

package ossClient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/zhongling01/oss-go-sdk/pkg/s3utils"
	"hash"
	"hash/crc32"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

type MultipartUploader struct {
	BucketName string
	ObjectName string
	UploadID   string
	c          *Client
	opts       *PutObjectOptions
	partsInfo  map[int]ObjectPart
	buf        []byte
	// Create checksums
	// CRC32C is ~50% faster on AMD64 @ 30GB/s
	crcBytes  map[int][]byte
	crc       hash.Hash32
	eof       bool // The sequential upload has been completed, now can UpdatePart or CompleteMultipartUpload
	completed bool
}

// getUploadID - fetch upload id if already present for an object name
// or initiate a new request to fetch a new upload id.
func (c *Client) NewUploadID(ctx context.Context, bucketName, objectName string, opts *PutObjectOptions) (*MultipartUploader, error) {
	if opts.DisableMultipart {
		return nil, errors.New("multipart disabled")
	}

	var objectSize int64 = -1
	//if objectSize == 0 {
	//	return nil, errors.New("objectSize is illegal")
	//}
	//if objectSize < 0 {
	//	objectSize = -1
	//}

	_, _, _, err := OptimalPartInfo(objectSize, opts.PartSize)
	if err != nil {
		return nil, err
	}

	if !opts.SendContentMd5 {
		if opts.UserMetadata == nil {
			opts.UserMetadata = make(map[string]string, 1)
		}
		opts.UserMetadata["X-Amz-Checksum-Algorithm"] = "CRC32C"
	}

	// Input validation.
	if err = s3utils.CheckValidBucketName(bucketName); err != nil {
		return nil, err
	}
	if err = s3utils.CheckValidObjectName(objectName); err != nil {
		return nil, err
	}

	// Initiate multipart upload for an object.
	initMultipartUploadResult, err := c.initiateMultipartUpload(ctx, bucketName, objectName, *opts)
	if err != nil {
		return nil, err
	}

	delete(opts.UserMetadata, "X-Amz-Checksum-Algorithm")

	return &MultipartUploader{
		BucketName: bucketName,
		ObjectName: objectName,
		UploadID:   initMultipartUploadResult.UploadID,
		c:          c,
		opts:       opts,
		partsInfo:  make(map[int]ObjectPart),    // Initialize parts uploaded map
		buf:        make([]byte, opts.PartSize), // Create a buffer
		crcBytes:   make(map[int][]byte),
		crc:        crc32.New(crc32.MakeTable(crc32.Castagnoli)),
	}, nil
}

func (m *MultipartUploader) uploadPart(ctx context.Context, buf []byte, length, partNumber int) error {
	var md5Base64 string
	customHeader := make(http.Header)
	if m.opts.SendContentMd5 {
		// Calculate md5sum.
		hash := m.c.md5Hasher()
		hash.Write(buf[:length])
		md5Base64 = base64.StdEncoding.EncodeToString(hash.Sum(nil))
		hash.Close()
	} else {
		m.crc.Reset()
		m.crc.Write(buf[:length])
		cSum := m.crc.Sum(nil)
		customHeader.Set("x-amz-checksum-crc32c", base64.StdEncoding.EncodeToString(cSum))
		m.crcBytes[partNumber] = cSum
	}

	// Update progress reader appropriately to the latest offset
	// as we read from the source.
	rd := newHook(bytes.NewReader(buf[:length]), m.opts.Progress)

	// Proceed to upload the part.
	p := uploadPartParams{bucketName: m.BucketName, objectName: m.ObjectName, uploadID: m.UploadID,
		reader: rd, partNumber: partNumber, md5Base64: md5Base64,
		size: int64(length), sse: m.opts.ServerSideEncryption, streamSha256: !m.opts.DisableContentSha256,
		customHeader: customHeader}

	objPart, err := m.c.uploadPart(ctx, p)
	if err != nil {
		return err
	}

	// Save successfully uploaded part metadata.
	m.partsInfo[partNumber] = objPart

	return nil
}

// UploadPart - Uploads a part in a multipart upload.
func (m *MultipartUploader) UploadPart(ctx context.Context, data io.Reader, partNumber int) error {
	length, rerr := readFull(data, m.buf)
	// For unknown size, Read EOF we break away.
	// We do not have to upload till totalPartsCount.
	if rerr == io.EOF {
		m.eof = true
		return io.EOF
	}

	if rerr != nil && rerr != io.ErrUnexpectedEOF && rerr != io.EOF {
		return rerr
	}

	return m.uploadPart(ctx, m.buf, length, partNumber)
}

// UpdatePart - Update the uploaded part.
func (m *MultipartUploader) UpdatePart(ctx context.Context, data io.Reader, partNumber int, configuredPartSize int) error {
	if m.completed {
		return errors.New("upload is completed")
	}

	if partNumber < len(m.partsInfo) {
		if configuredPartSize < absMinPartSize {
			return errInvalidArgument("Input part size is smaller than allowed minimum of 5MiB.")
		}
		if configuredPartSize > maxPartSize {
			return errInvalidArgument("Input part size is bigger than allowed maximum of 5GiB.")
		}
	} else if partNumber > len(m.partsInfo) {
		return errors.New("partNumber is illegal")
	}

	buf := make([]byte, configuredPartSize)
	length, rerr := readFull(data, buf)
	if length != configuredPartSize {
		return errors.New(fmt.Sprintf("read %d bytes from data is smaller than %d", length, configuredPartSize))
	}

	if rerr != nil && rerr != io.ErrUnexpectedEOF && rerr != io.EOF {
		return rerr
	}

	return m.uploadPart(ctx, buf, length, partNumber)
}

func (m *MultipartUploader) GetPart(ctx context.Context, partNumber int) (io.ReadCloser, ObjectInfo, error) {
	return m.getPart(ctx, partNumber, GetObjectOptions{})
}

func (m *MultipartUploader) getPart(ctx context.Context, partNumber int, opts GetObjectOptions) (io.ReadCloser, ObjectInfo, error) {
	if m.completed {
		return nil, ObjectInfo{}, errors.New("upload is completed")
	}

	if partNumber > len(m.partsInfo) {
		return nil, ObjectInfo{}, errors.New("partNumber is illegal")
	}
	urlValues := make(url.Values)
	urlValues.Set("partNumber", strconv.Itoa(partNumber))
	urlValues.Set("uploadingID", m.UploadID)

	// Execute GET on objectName.
	resp, err := m.c.executeMethod(ctx, http.MethodGet, requestMetadata{
		bucketName:       m.BucketName,
		objectName:       m.ObjectName,
		queryValues:      urlValues,
		customHeader:     opts.Header(),
		contentSHA256Hex: emptySHA256Hex,
	})
	if err != nil {
		return nil, ObjectInfo{}, err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
			return nil, ObjectInfo{}, errors.New(fmt.Sprintf("http request %s", resp.Status))
		}
	}

	objectStat, err := ToObjectInfo(m.BucketName, m.ObjectName, resp.Header)
	if err != nil {
		closeResponse(resp)
		return nil, ObjectInfo{}, err
	}

	return resp.Body, objectStat, nil
}

// CompleteMultipartUpload - Completes a multipart upload by assembling previously uploaded parts.
func (m *MultipartUploader) CompleteMultipartUpload(ctx context.Context) (UploadInfo, error) {
	if !m.eof {
		return UploadInfo{}, errors.New("UploadPart doesn't end")
	}
	if m.completed {
		return UploadInfo{}, errors.New("upload is completed")
	}
	m.completed = true

	// Complete multipart upload.
	var complete completeMultipartUpload
	// Loop over total uploaded parts to save them in
	// Parts array before completing the multipart request.
	for i := 1; i <= len(m.partsInfo); i++ {
		part, ok := m.partsInfo[i]
		if !ok {
			return UploadInfo{}, errInvalidArgument(fmt.Sprintf("Missing part number %d", i))
		}
		complete.Parts = append(complete.Parts, CompletePart{
			ETag:           part.ETag,
			PartNumber:     part.PartNumber,
			ChecksumCRC32:  part.ChecksumCRC32,
			ChecksumCRC32C: part.ChecksumCRC32C,
			ChecksumSHA1:   part.ChecksumSHA1,
			ChecksumSHA256: part.ChecksumSHA256,
		})
	}

	// Sort all completed parts.
	sort.Sort(completedParts(complete.Parts))

	opts := &PutObjectOptions{
		ServerSideEncryption: m.opts.ServerSideEncryption,
		MergeMultipart:       m.opts.MergeMultipart,
	}
	if len(m.crcBytes) > 0 {
		// Add hash of hashes.
		m.crc.Reset()
		var crcBytes []byte
		for i := 1; i <= len(m.crcBytes); i++ {
			crcBytes = append(crcBytes, m.crcBytes[i]...)
		}
		m.crc.Write(crcBytes)
		opts.UserMetadata = map[string]string{"X-Amz-Checksum-Crc32c": base64.StdEncoding.EncodeToString(m.crc.Sum(nil))}
	}

	// Input validation.
	if err := s3utils.CheckValidBucketName(m.BucketName); err != nil {
		return UploadInfo{}, err
	}
	if err := s3utils.CheckValidObjectName(m.ObjectName); err != nil {
		return UploadInfo{}, err
	}

	// Initialize url queries.
	urlValues := make(url.Values)
	urlValues.Set("uploadId", m.UploadID)
	// Marshal complete multipart body.
	completeMultipartUploadBytes, err := xml.Marshal(complete)
	if err != nil {
		return UploadInfo{}, err
	}

	// Instantiate all the complete multipart buffer.
	completeMultipartUploadBuffer := bytes.NewReader(completeMultipartUploadBytes)
	reqMetadata := requestMetadata{
		bucketName:       m.BucketName,
		objectName:       m.ObjectName,
		queryValues:      urlValues,
		contentBody:      completeMultipartUploadBuffer,
		contentLength:    int64(len(completeMultipartUploadBytes)),
		contentSHA256Hex: sum256Hex(completeMultipartUploadBytes),
		customHeader:     opts.Header(),
	}

	// Execute POST to complete multipart upload for an objectName.
	resp, err := m.c.executeMethod(ctx, http.MethodPost, reqMetadata)
	defer closeResponse(resp)
	if err != nil {
		return UploadInfo{}, err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			return UploadInfo{}, httpRespToErrorResponse(resp, m.BucketName, m.ObjectName)
		}
	}

	// Read resp.Body into a []bytes to parse for Error response inside the body
	var b []byte
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return UploadInfo{}, err
	}
	// Decode completed multipart upload response on success.
	completeMultipartUploadResult := completeMultipartUploadResult{}
	err = xmlDecoder(bytes.NewReader(b), &completeMultipartUploadResult)
	if err != nil {
		// xml parsing failure due to presence an ill-formed xml fragment
		return UploadInfo{}, err
	} else if completeMultipartUploadResult.Bucket == "" {
		// xml's Decode method ignores well-formed xml that don't apply to the type of value supplied.
		// In this case, it would leave completeMultipartUploadResult with the corresponding zero-values
		// of the members.

		// Decode completed multipart upload response on failure
		completeMultipartUploadErr := ErrorResponse{}
		err = xmlDecoder(bytes.NewReader(b), &completeMultipartUploadErr)
		if err != nil {
			// xml parsing failure due to presence an ill-formed xml fragment
			return UploadInfo{}, err
		}
		return UploadInfo{}, completeMultipartUploadErr
	}

	// extract lifecycle expiry date and rule ID
	expTime, ruleID := amzExpirationToExpiryDateRuleID(resp.Header.Get(amzExpiration))

	var totalUploadedSize int64
	for _, partInfo := range m.partsInfo {
		totalUploadedSize += partInfo.Size
	}

	return UploadInfo{
		Bucket:           completeMultipartUploadResult.Bucket,
		Key:              completeMultipartUploadResult.Key,
		ETag:             trimEtag(completeMultipartUploadResult.ETag),
		Size:             totalUploadedSize,
		VersionID:        resp.Header.Get(amzVersionID),
		Location:         completeMultipartUploadResult.Location,
		Expiration:       expTime,
		ExpirationRuleID: ruleID,

		ChecksumSHA256: completeMultipartUploadResult.ChecksumSHA256,
		ChecksumSHA1:   completeMultipartUploadResult.ChecksumSHA1,
		ChecksumCRC32:  completeMultipartUploadResult.ChecksumCRC32,
		ChecksumCRC32C: completeMultipartUploadResult.ChecksumCRC32C,
	}, nil
}

// AbortMultipartUpload aborts a multipart upload for the given
// uploadID, all previously uploaded parts are deleted.
func (m *MultipartUploader) AbortMultipartUpload(ctx context.Context) error {
	bucketName := m.BucketName
	objectName := m.ObjectName
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return err
	}
	if err := s3utils.CheckValidObjectName(objectName); err != nil {
		return err
	}

	// Initialize url queries.
	urlValues := make(url.Values)
	urlValues.Set("uploadId", m.UploadID)

	// Execute DELETE on multipart upload.
	resp, err := m.c.executeMethod(ctx, http.MethodDelete, requestMetadata{
		bucketName:       bucketName,
		objectName:       objectName,
		queryValues:      urlValues,
		contentSHA256Hex: emptySHA256Hex,
	})
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusNoContent {
			// Abort has no response body, handle it for any errors.
			var errorResponse ErrorResponse
			switch resp.StatusCode {
			case http.StatusNotFound:
				// This is needed specifically for abort and it cannot
				// be converged into default case.
				errorResponse = ErrorResponse{
					Code:       "NoSuchUpload",
					Message:    "The specified multipart upload does not exist.",
					BucketName: bucketName,
					Key:        objectName,
					RequestID:  resp.Header.Get("x-amz-request-id"),
					HostID:     resp.Header.Get("x-amz-id-2"),
					Region:     resp.Header.Get("x-amz-bucket-region"),
				}
			default:
				return httpRespToErrorResponse(resp, bucketName, objectName)
			}
			return errorResponse
		}
	}
	return nil
}
