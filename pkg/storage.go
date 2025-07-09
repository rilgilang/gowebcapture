package pkg

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
)

type Storage interface {
	Put(ctx context.Context, bucket string, fileName string, file []byte, size int64, cachable bool, contentType string) error
	Get(ctx context.Context, bucket string, fileName string) ([]byte, error)
}

type storagePkg struct {
	client *minio.Client
}

func NewStorage(client *minio.Client) Storage {
	return &storagePkg{
		client: client,
	}
}

func (s *storagePkg) Put(ctx context.Context, bucket string, fileName string, file []byte, size int64, cachable bool, contentType string) error {
	cacheControl := "public, max-age=86400"
	if !cachable {
		cacheControl = "no-cache, max-age=0"
	}

	r := bytes.NewReader(file)

	_, err := s.client.PutObject(ctx, bucket, fileName, r, size, minio.PutObjectOptions{
		CacheControl: cacheControl,
		ContentType:  contentType,
	})

	if err != nil {
		fmt.Println("error upload: ", err)
		return err
	}

	return err
}

func (s *storagePkg) Get(ctx context.Context, bucket string, fileName string) (file []byte, err error) {
	object, err := s.client.GetObject(ctx, bucket, fileName, minio.GetObjectOptions{})

	if err != nil {
		return nil, err
	}

	defer object.Close()

	_, err = object.Read(file)
	if err != nil {
		fmt.Println("error get: ", err)
		return nil, err
	}

	return file, nil
}
