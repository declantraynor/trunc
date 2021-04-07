package storage

import (
	"errors"
	"testing"

	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
)

var conn = redigomock.NewConn()
var store = &RedisStore{conn: conn}

func TestGet(t *testing.T) {
	conn.Clear()

	conn.Command("GET", "foo").Expect("bar")
	result, err := store.Get("foo")

	assert.Equal(t, "bar", result)
	assert.Nil(t, err)
}

func TestGetError(t *testing.T) {
	conn.Clear()

	expectedErr := errors.New("GET failed")
	conn.Command("GET", "foo").ExpectError(expectedErr)

	result, err := store.Get("foo")
	assert.Equal(t, "", result)
	assert.Equal(t, expectedErr, err)
}

func TestSet(t *testing.T) {
	conn.Clear()
	conn.Command("SET", "foo", "bar").Expect(nil)

	err := store.Set("foo", "bar")
	assert.Nil(t, err)
}

func TestSetError(t *testing.T) {
	conn.Clear()

	expectedErr := errors.New("SET failed")
	conn.Command("SET", "foo", "bar").ExpectError(expectedErr)
	err := store.Set("foo", "bar")
	assert.Equal(t, expectedErr, err)
}
