package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	ctx context.Context
	client *redis.Client
}

func NewRedis(host, pass string, db int) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass,
		DB:       db,
	})
	return &Redis{ctx: context.Background(), client: client}
}

func (r *Redis) Set(key string, value interface{}, expire time.Duration) (err error) {
	err = r.client.Set(r.ctx, key, value, expire).Err()
	return
}

func (r *Redis) Get(key string) (result string, err error) {
	result, err = r.client.Get(r.ctx, key).Result()
	return
}

func (r *Redis) HasKey(key string) bool {
	result := r.client.Exists(r.ctx, key).Val()
	return result == 1
}

func (r *Redis) HSet(key, field string, value interface{}, expire time.Duration) (err error) {
	err = r.client.HSet(r.ctx, key, field, value).Err()
	if err != nil {
		return err
	}

	if expire != 0 {
		err = r.client.Expire(r.ctx, key, expire).Err()
		if err != nil {
			r.client.Del(r.ctx, key)
			return nil
		}
	}
	return
}

func (r *Redis) HGet(key, field string) (result string, err error) {
	result, err = r.client.HGet(r.ctx, key, field).Result()
	return
}

func (r *Redis) HExists(key, field string) bool {
	val := r.client.HExists(r.ctx, key, field).Val()
	return val
}

func (r *Redis) SAdd(key string, value interface{}) (err error) {
	err = r.client.SAdd(r.ctx, key, value).Err()
	return
}

func (r *Redis) SRem(key string, value interface{}) (bool, error) {
	result, err := r.client.SRem(r.ctx, key, value).Result()
	return result == 1, err
}

func (r *Redis) SIsMember(key string, value interface{}) (ok bool, err error) {
	ok, err = r.client.SIsMember(r.ctx, key, value).Result()
	return
}

func (r *Redis) GetBit(key string, offset int64) (int64, error){
	result, err := r.client.GetBit(r.ctx, key, offset).Result()
	return result, err
}

func (r *Redis) SetBit(key string, offset int64, value int) (int64, error){
	result, err := r.client.SetBit(r.ctx, key, offset, value).Result()
	return result, err
}

func (r *Redis) BitCount(key string, start, end int64) (int64, error) {
	result, err := r.client.BitCount(r.ctx, key, &redis.BitCount{start, end}).Result()
	return result, err
}