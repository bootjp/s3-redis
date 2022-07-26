package app

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"

	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/minio/minio-go/v7"
)

type S3 struct {
	Client *minio.Client
	prefix string
	bucket string
}

func NewS3(config Config) *S3 {
	c, err := minio.New(config.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3.AccessKey, config.S3.SecretKey, ""),
		Secure: true,
		Region: config.S3.Region,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &S3{
		Client: c,
		bucket: config.S3.Bucket,
	}
}

func (s *S3) Put(ctx context.Context, key string, value []byte) error {
	r := bytes.NewReader(value)
	_, err := s.Client.PutObject(ctx, s.bucket, key, r, int64(len(value)), minio.PutObjectOptions{})

	return err
}
func (s *S3) Get(ctx context.Context, key string) ([]byte, error) {
	o, err := s.Client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(o)
}

func (s *S3) Delete(ctx context.Context, key string) error {
	return s.Client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}

func (s *S3) List(ctx context.Context, key string) ([]string, error) {
	objectCh := s.Client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    s.prefix,
		Recursive: true,
	})

	var result []string
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		result = append(result, object.Key)

	}

	return result, nil
}
