package playwithredis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.ClusterClient
}

func NewClient(ctx context.Context, addrs []string) (*Client, error) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{Addrs: addrs})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{rdb}, nil
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) Clear(ctx context.Context) error {
	return c.rdb.FlushAll(ctx).Err()
}

func (c *Client) StoreReview(ctx context.Context, r Review) (ID, error) {
	id, err := c.rdb.Incr(ctx, "reviews:sequence").Result()
	if err != nil {
		return 0, fmt.Errorf("store review: %w", err)
	}
	key := fmt.Sprintf("review:%d", id)

	if err := c.rdb.HSet(ctx, key, "text", r.Text).Err(); err != nil {
		return 0, fmt.Errorf("store review: %w", err)
	}
	if err := c.rdb.HSet(ctx, key, "score", r.Score).Err(); err != nil {
		return 0, fmt.Errorf("store review: %w", err)
	}
	if err := c.rdb.HSet(ctx, key, "type", r.Type).Err(); err != nil {
		return 0, fmt.Errorf("store review: %w", err)
	}

	return ID(id), nil
}

func (c *Client) LoadReview(ctx context.Context, id ID) (Review, error) {
	key := fmt.Sprintf("review:%d", id)

	text, err := c.rdb.HGet(ctx, key, "text").Result()
	if err != nil {
		return Review{}, fmt.Errorf("load review: %w", err)
	}
	score, err := c.rdb.HGet(ctx, key, "score").Int()
	if err != nil {
		return Review{}, fmt.Errorf("load review: %w", err)
	}
	rtype, err := c.rdb.HGet(ctx, key, "type").Result()
	if err != nil {
		return Review{}, fmt.Errorf("load review: %w", err)
	}

	return Review{
		Text:  text,
		Score: score,
		Type:  rtype,
	}, nil
}

func (c *Client) StoreMovie(ctx context.Context, m Movie) error {
	key := fmt.Sprintf("movie-reviews:%d", m.ID)

	for _, review := range m.Reviews {
		id, err := c.StoreReview(ctx, review)
		if err != nil {
			return fmt.Errorf("store movie: %w", err)
		}
		if err := c.rdb.RPush(ctx, key, int64(id)).Err(); err != nil {
			return fmt.Errorf("store movie: %w", err)
		}
	}

	return nil
}

func (c *Client) LoadMovie(ctx context.Context, id ID) (Movie, error) {
	key := fmt.Sprintf("movie-reviews:%d", id)

	reviewsCount, err := c.rdb.LLen(ctx, key).Result()
	if err != nil {
		return Movie{}, fmt.Errorf("load movie: %w", err)
	}
	reviews := make([]Review, 0, reviewsCount)

	reviewIDs, err := c.rdb.LRange(ctx, key, 0, reviewsCount-1).Result()
	if err != nil {
		return Movie{}, err
	}

	for _, idStr := range reviewIDs {
		parsedID, _ := strconv.Atoi(idStr)
		review, err := c.LoadReview(ctx, ID(parsedID))
		if err != nil {
			return Movie{}, fmt.Errorf("load movie: %w", err)
		}

		reviews = append(reviews, review)
	}

	// for i := int64(0); i < reviewsCount; i++ {
	// 	reviewID, err := c.rdb.LIndex(ctx, key, i).Int64()
	// 	if err != nil {
	// 		return Movie{}, fmt.Errorf("load movie: %w", err)
	// 	}

	// 	review, err := c.LoadReview(ctx, ID(reviewID))
	// 	if err != nil {
	// 		return Movie{}, fmt.Errorf("load movie: %w", err)
	// 	}

	// 	reviews = append(reviews, review)
	// }

	return Movie{ID: id, Reviews: reviews}, nil
}
