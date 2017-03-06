package cache

import (
	"bytes"
	"errors"
	"strconv"
	"time"

	"github.com/peterbourgon/diskv"
)

// Cache struct
type Cache struct {
	storage    *diskv.Diskv
	storageDir string
}

// New returns an initialized Cache instance
func New(storageDir string) *Cache {
	c := new(Cache)

	c.storageDir = storageDir

	c.initCache()

	return c
}

// initCache initializes the cache storage
func (c *Cache) initCache() {
	c.storage = diskv.New(diskv.Options{
		BasePath:     c.storageDir,
		Transform:    func(s string) []string { return []string{} },
		CacheSizeMax: 1024 * 1024,
	})
}

// Save saves a byte slice to cache
func (c *Cache) Save(key string, data []byte, ttl int64) error {
	// Calculate expiry time
	n := time.Now()
	expire := n.Unix() + ttl

	// Merge the data and the expiry time
	dataTmp := []byte{}
	dataTmp = append(dataTmp, []byte(strconv.Itoa(int(expire)))...)
	dataTmp = append(dataTmp, ':')
	dataTmp = append(dataTmp, data...)

	return c.storage.Write(key, dataTmp)
}

// Get fetches a byte slice from cache
func (c *Cache) Get(key string) ([]byte, error) {
	data, err := c.storage.Read(key)
	if err != nil {
		return []byte{}, err
	}

	sep := bytes.IndexAny(data, ":")
	if sep == -1 {
		return []byte{}, errors.New("Couldn't read cache entry")
	}

	expire, err := strconv.Atoi(string(data[:sep]))
	if err != nil {
		return []byte{}, errors.New("Couldn't read cache entry")
	}

	if int64(expire) <= time.Now().Unix() {
		return []byte{}, errors.New("Entry expired")
	}

	sep++
	return data[sep:], nil
}
