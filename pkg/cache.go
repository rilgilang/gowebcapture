package pkg

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Lpush(ctx context.Context, key, url string) (err error)
	BRpop(ctx context.Context, key string) (value string, err error)
}

type cachePkg struct {
	client *redis.Client
}

func NewCache() Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "wXU7z2HQtBP6mWjHbPsfK96rmnqweawg", // no password set
		DB:       0,                                  // use default DB
	})

	return &cachePkg{
		client: client,
	}
}

func (c *cachePkg) Lpush(ctx context.Context, key, url string) error {
	return c.client.LPush(ctx, key, url).Err()
}

func (c *cachePkg) BRpop(ctx context.Context, key string) (value string, err error) {

	result := c.client.BRPop(ctx, 0, key)
	return result.Val()[1], result.Err()

}
