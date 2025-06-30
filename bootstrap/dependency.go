package bootstrap

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

type BoostrapClient struct {
	Storage *minio.Client
	Cache   *redis.Client
}

type Config struct {
	StorageAccessKey       string `json:"storage_access_key"`
	StorageSecretAccessKey string `json:"storage_secret_access_key"`
	StorageEndpoint        string `json:"storage_endpoint"`
	StorageBucket          string `json:"storage_bucket"`
	StorageSecure          bool   `json:"storage_secure"`
	RedisHost              string `json:"redis_host"`
	RedisPassword          string `json:"redis_password"`
	RedisDB                int    `json:"redis_db"`
}

func Setup() (client *BoostrapClient, config *Config, err error) {
	config = setupConfig()

	fmt.Println("config --> ", config)

	// S3
	storage, err := minio.New(config.StorageEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.StorageAccessKey, config.StorageSecretAccessKey, ""),
		Secure: config.StorageSecure},
	)

	if err != nil {
		return nil, nil, err
	}

	// Redis
	cache := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	return &BoostrapClient{
		Storage: storage,
		Cache:   cache,
	}, config, nil
}

func setupConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	storageSecure, err := strconv.ParseBool(os.Getenv("STORAGE_SECURE"))
	if err != nil {
		panic(err)
	}

	return &Config{
		StorageAccessKey:       os.Getenv("STORAGE_ACCESS_KEY"),
		StorageSecretAccessKey: os.Getenv("STORAGE_SECRET_ACCESS_KEY"),
		StorageEndpoint:        os.Getenv("STORAGE_ENDPOINT"),
		StorageBucket:          os.Getenv("STORAGE_BUCKET"),
		StorageSecure:          storageSecure,
		RedisHost:              os.Getenv("REDIS_HOST"),
		RedisPassword:          os.Getenv("REDIS_PASSWORD"),
		RedisDB:                0,
	}
}
