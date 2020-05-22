package boltdb

import (
	"log"

	"github.com/netscale-technologies/gorush/config"
	"github.com/netscale-technologies/gorush/storage"

	"github.com/asdine/storm/v3"
)

// New func implements the storage interface for gorush (https://github.com/appleboy/gorush)
func New(config config.ConfYaml) *Storage {
	return &Storage{
		config: config,
	}
}

// Storage is interface structure
type Storage struct {
	config config.ConfYaml
	db     *storm.DB
}

// Init client storage.
func (s *Storage) Init() error {
	var err error
	s.db, err = storm.Open(s.config.Stat.BoltDB.Path)
	return err
}

// Close the storage connection
func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}

	return s.db.Close()
}

// Reset Client storage.
func (s *Storage) Reset() {
	s.setBoltDB(storage.TotalCountKey, 0)
	s.setBoltDB(storage.IosSuccessKey, 0)
	s.setBoltDB(storage.IosErrorKey, 0)
	s.setBoltDB(storage.AndroidSuccessKey, 0)
	s.setBoltDB(storage.AndroidErrorKey, 0)
	s.setBoltDB(storage.WebSuccessKey, 0)
	s.setBoltDB(storage.WebErrorKey, 0)
}

func (s *Storage) setBoltDB(key string, count int64) {
	err := s.db.Set(s.config.Stat.BoltDB.Bucket, key, count)
	if err != nil {
		log.Println("BoltDB set error:", err.Error())
	}
}

func (s *Storage) getBoltDB(key string, count *int64) {
	err := s.db.Get(s.config.Stat.BoltDB.Bucket, key, count)
	if err != nil {
		log.Println("BoltDB get error:", err.Error())
	}
}

// AddTotalCount record push notification count.
func (s *Storage) AddTotalCount(count int64) {
	total := s.GetTotalCount() + count
	s.setBoltDB(storage.TotalCountKey, total)
}

// AddIosSuccess record counts of success iOS push notification.
func (s *Storage) AddIosSuccess(count int64) {
	total := s.GetIosSuccess() + count
	s.setBoltDB(storage.IosSuccessKey, total)
}

// AddIosError record counts of error iOS push notification.
func (s *Storage) AddIosError(count int64) {
	total := s.GetIosError() + count
	s.setBoltDB(storage.IosErrorKey, total)
}

// AddAndroidSuccess record counts of success Android push notification.
func (s *Storage) AddAndroidSuccess(count int64) {
	total := s.GetAndroidSuccess() + count
	s.setBoltDB(storage.AndroidSuccessKey, total)
}

// AddAndroidError record counts of error Android push notification.
func (s *Storage) AddAndroidError(count int64) {
	total := s.GetAndroidError() + count
	s.setBoltDB(storage.AndroidErrorKey, total)
}

// AddWebSuccess record counts of success Web push notification.
func (s *Storage) AddWebSuccess(count int64) {
	total := s.GetWebSuccess() + count
	s.setBoltDB(storage.WebSuccessKey, total)
}

// AddWebError record counts of error Web push notification.
func (s *Storage) AddWebError(count int64) {
	total := s.GetWebError() + count
	s.setBoltDB(storage.WebErrorKey, total)
}

// GetTotalCount show counts of all notification.
func (s *Storage) GetTotalCount() int64 {
	var count int64
	s.getBoltDB(storage.TotalCountKey, &count)

	return count
}

// GetIosSuccess show success counts of iOS notification.
func (s *Storage) GetIosSuccess() int64 {
	var count int64
	s.getBoltDB(storage.IosSuccessKey, &count)

	return count
}

// GetIosError show error counts of iOS notification.
func (s *Storage) GetIosError() int64 {
	var count int64
	s.getBoltDB(storage.IosErrorKey, &count)

	return count
}

// GetAndroidSuccess show success counts of Android notification.
func (s *Storage) GetAndroidSuccess() int64 {
	var count int64
	s.getBoltDB(storage.AndroidSuccessKey, &count)

	return count
}

// GetAndroidError show error counts of Android notification.
func (s *Storage) GetAndroidError() int64 {
	var count int64
	s.getBoltDB(storage.AndroidErrorKey, &count)

	return count
}

// GetWebSuccess show success counts of Web notification.
func (s *Storage) GetWebSuccess() int64 {
	var count int64
	s.getBoltDB(storage.WebSuccessKey, &count)

	return count
}

// GetWebError show error counts of Web notification.
func (s *Storage) GetWebError() int64 {
	var count int64
	s.getBoltDB(storage.WebErrorKey, &count)

	return count
}
