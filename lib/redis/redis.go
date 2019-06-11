package redis

import (
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

// Client ...
type Client struct {
	conn redis.Conn
}

// newPool redis pool
func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			redisMasterServer := fmt.Sprintf("%s:%s", os.Getenv("ELASTIC_CACHE_HOST"), os.Getenv("ELASTIC_CACHE_PORT"))
			c, err := redis.Dial("tcp", redisMasterServer)
			if err != nil {
				fmt.Println("err", err.Error())
			}

			return c, err
		},
	}
}

// Init newpool ...
func Init() {
	pool = newPool()
}

// MasterConnect ...
func MasterConnect() (conn *Client) {
	if pool == nil {
		pool = newPool()
	}
	instanceRedisCli := new(Client)
	instanceRedisCli.conn = pool.Get()
	return instanceRedisCli

}

// Close redis connection
func (redisCli *Client) Close() {
	redisCli.conn.Close()
}

// Exists ...
func (redisCli *Client) Exists(key string) (bool, error) {
	isExists, err := redis.Bool(redisCli.conn.Do("EXISTS", key))
	return isExists, err
}

// SetValue string set
func (redisCli *Client) SetValue(key string, value string, expireTime int) error {
	_, err := redisCli.conn.Do("SET", key, value)

	if err == nil && expireTime > 0 {
		redisCli.conn.Do("EXPIRE", key, expireTime)
	}

	if err != nil {
		return err
	}
	return nil
}

// GetValue ...
func (redisCli *Client) GetValue(key string) (string, error) {
	return redis.String(redisCli.conn.Do("GET", key))
}

// DelValue ...
func (redisCli *Client) DelValue(key string) error {
	_, err := redisCli.conn.Do("DEL", key)
	if err != nil {
		return err
	}
	return nil
}
