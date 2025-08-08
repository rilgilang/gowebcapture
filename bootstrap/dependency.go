package bootstrap

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"os"
	"strconv"
)

type BoostrapClient struct {
	DB      *gorm.DB
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
	DBDialect              string `json:"db_dialect"`
	DBHost                 string `json:"db_host"`
	DBPort                 int    `json:"db_port"`
	DBName                 string `json:"DB_NAME"`
	DBUsername             string `json:"db_username"`
	DBPassword             string `json:"db_password"`
	LinuxBrowserPath       string `json:"linux_browser_path"`
	DarwinBrowserPath      string `json:"darwin_browser_path"`
	FFMPEGFramerate        string `json:"ffmpeg_framerate"`
	FFMPEGVideoSize        string `json:"ffmpeg_video_size"`
	FFMPEGCrop             bool   `json:"ffmpeg_crop"`
	FFMPEGCropSize         string `json:"ffmpeg_crop_size"`
}

func Setup() (client *BoostrapClient, config *Config, err error) {
	config = setupConfig()

	fmt.Println("config --> ", config)

	// S3
	storage, err := minio.New(config.StorageEndpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(config.StorageAccessKey, config.StorageSecretAccessKey, ""),
		Secure:       config.StorageSecure,
		BucketLookup: minio.BucketLookupPath,
	},
	)
	if err != nil {
		//fmt.Println("Error connecting to MinIO:", err)
		panic(err)
	}

	// minio health check
	buckets, err := storage.ListBuckets(context.Background())
	if err != nil {
		fmt.Println("Error get list of bucket MinIO:", err)
	}

	fmt.Println("minio bucket --> ", buckets)

	// Redis
	cache := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	db, err := DatabaseConnection(config)
	if err != nil {
		panic(err)
	}

	return &BoostrapClient{
		DB:      db,
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

	ffmpegCrop, err := strconv.ParseBool(os.Getenv("FFMPEG_CROP"))
	if err != nil {
		panic(err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
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
		DBHost:                 os.Getenv("DB_HOST"),
		DBDialect:              os.Getenv("DB_DIALECT"),
		DBUsername:             os.Getenv("DB_USERNAME"),
		DBName:                 os.Getenv("DB_NAME"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBPort:                 dbPort,
		DarwinBrowserPath:      os.Getenv("DARWIN_BROWSER_PATH"),
		LinuxBrowserPath:       os.Getenv("LINUX_BROWSER_PATH"),
		FFMPEGFramerate:        os.Getenv("FFMPEG_FRAMERATE"),
		FFMPEGVideoSize:        os.Getenv("FFMPEG_VIDEO_SIZE"),
		FFMPEGCrop:             ffmpegCrop,
		FFMPEGCropSize:         os.Getenv("FFMPEG_CROP_SIZE"),
	}
}
