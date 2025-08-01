package pkg

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"os"
	"path/filepath"
)

type MinioFile struct {
	Bytes []byte
	Size  int64
}

type Storage interface {
	FileToBytes(filename string) (*MinioFile, error)
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
func (s *storagePkg) FileToBytes(filename string) (*MinioFile, error) {
	dir, _ := os.Getwd() // Used for output file
	filePath := filepath.Join(dir, filename)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("err open file: ", err)
		return nil, err

	}

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Error reading file stats: ", err)
		return nil, err
	}
	stat.Size()

	defer file.Close() // Ensure the file is closed after function exits

	// Read all content from the file into a byte slice
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return nil, err
	}

	return &MinioFile{
		data,
		stat.Size(),
	}, err
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
