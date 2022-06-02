package redis

import (
	"strconv"

	"github.com/netscale-technologies/gorush/config"
	"github.com/netscale-technologies/gorush/storage"

	"github.com/go-redis/redis/v7"
)

// New func implements the storage interface for gorush (https://github.com/appleboy/gorush)
func New(config config.ConfYaml) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) getInt64(key string, count *int64) {
	val, _ := s.client.Get(key).Result()
	*count, _ = strconv.ParseInt(val, 10, 64)
}

// Storage is interface structure
type Storage struct {
	config config.ConfYaml
	client *redis.Client
}

// Init client storage.
func (s *Storage) Init() error {
	s.client = redis.NewClient(&redis.Options{
		Addr:     s.config.Stat.Redis.Addr,
		Password: s.config.Stat.Redis.Password,
		DB:       s.config.Stat.Redis.DB,
	})
	_, err := s.client.Ping().Result()

	return err
}

// Close the storage connection
func (s *Storage) Close() error {
	if s.client == nil {
		return nil
	}

	return s.client.Close()
}

// Reset Client storage.
func (s *Storage) Reset() {
	s.client.Set(storage.TotalCountKey, int64(0), 0)
	s.client.Set(storage.IosSuccessKey, int64(0), 0)
	s.client.Set(storage.IosErrorKey, int64(0), 0)
	s.client.Set(storage.AndroidSuccessKey, int64(0), 0)
	s.client.Set(storage.AndroidErrorKey, int64(0), 0)
	s.client.Set(storage.WebSuccessKey, int64(0), 0)
	s.client.Set(storage.WebErrorKey, int64(0), 0)
}

// AddTotalCount record push notification count.
func (s *Storage) AddTotalCount(count int64) {
	s.client.IncrBy(storage.TotalCountKey, count)
}

// AddIosSuccess record counts of success iOS push notification.
func (s *Storage) AddIosSuccess(count int64) {
	s.client.IncrBy(storage.IosSuccessKey, count)
}

// AddIosError record counts of error iOS push notification.
func (s *Storage) AddIosError(count int64) {
	s.client.IncrBy(storage.IosErrorKey, count)
}

// AddAndroidSuccess record counts of success Android push notification.
func (s *Storage) AddAndroidSuccess(count int64) {
	s.client.IncrBy(storage.AndroidSuccessKey, count)
}

// AddAndroidError record counts of error Android push notification.
func (s *Storage) AddAndroidError(count int64) {
	s.client.IncrBy(storage.AndroidErrorKey, count)
}

// AddWebSuccess record counts of success Web push notification.
func (s *Storage) AddWebSuccess(count int64) {
	total := s.GetWebSuccess() + count
	s.client.Set(storage.WebSuccessKey, strconv.Itoa(int(total)), 0)
}

// AddWebError record counts of error Web push notification.
func (s *Storage) AddWebError(count int64) {
	total := s.GetWebError() + count
	s.client.Set(storage.WebErrorKey, strconv.Itoa(int(total)), 0)
}

// GetTotalCount show counts of all notification.
func (s *Storage) GetTotalCount() int64 {
	var count int64
	s.getInt64(storage.TotalCountKey, &count)

	return count
}

// GetIosSuccess show success counts of iOS notification.
func (s *Storage) GetIosSuccess() int64 {
	var count int64
	s.getInt64(storage.IosSuccessKey, &count)

	return count
}

// GetIosError show error counts of iOS notification.
func (s *Storage) GetIosError() int64 {
	var count int64
	s.getInt64(storage.IosErrorKey, &count)

	return count
}

// GetAndroidSuccess show success counts of Android notification.
func (s *Storage) GetAndroidSuccess() int64 {
	var count int64
	s.getInt64(storage.AndroidSuccessKey, &count)

	return count
}

// GetAndroidError show error counts of Android notification.
func (s *Storage) GetAndroidError() int64 {
	var count int64
	s.getInt64(storage.AndroidErrorKey, &count)

	return count
}

// GetWebSuccess show success counts of Web notification.
func (s *Storage) GetWebSuccess() int64 {
	var count int64
	s.getInt64(storage.WebSuccessKey, &count)

	return count
}

// GetWebError show error counts of Web notification.
func (s *Storage) GetWebError() int64 {
	var count int64
	s.getInt64(storage.WebErrorKey, &count)

	return count
}
