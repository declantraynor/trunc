package storage

import (
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// Store describes a simple key/value store.
type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

// RedisStore is an implementation of Store backed by redis.
type RedisStore struct {
	conn redis.Conn
}

// Get returns the value stored at the specified key, or an error if the
// specified key does not exist.
func (r *RedisStore) Get(key string) (string, error) {
	return redis.String(r.conn.Do("GET", key))
}

// Set stores the specified key/value pair, returning any error returned
// by redis.
func (r *RedisStore) Set(key, value string) error {
	if _, err := r.conn.Do("SET", key, value); err != nil {
		return err
	}
	return nil
}

// Disconnect can be used to close the underlying redis connection used by
// a RedisStore.
func (r *RedisStore) Disconnect() {
	r.conn.Close()
}

// NewRedisStore returns a RedisStore with an initialised connection to redis,
// or an error if there is some problem connecting.
func NewRedisStore(host string, port int) (*RedisStore, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		return &RedisStore{}, errors.New(fmt.Sprintf("unable to connect to redis at %s", address))
	}
	return &RedisStore{conn: conn}, nil
}
