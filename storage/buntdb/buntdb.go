package buntdb

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jaraxasoftware/gorush/config"
	"github.com/jaraxasoftware/gorush/storage"

	"github.com/tidwall/buntdb"
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
	db     *buntdb.DB
}

// Init client storage.
func (s *Storage) Init() error {
	var err error
	s.db, err = buntdb.Open(s.config.Stat.BuntDB.Path)
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
	s.setBuntDB(storage.TotalCountKey, 0)
	s.setBuntDB(storage.IosSuccessKey, 0)
	s.setBuntDB(storage.IosErrorKey, 0)
	s.setBuntDB(storage.AndroidSuccessKey, 0)
	s.setBuntDB(storage.AndroidErrorKey, 0)
	s.setBuntDB(storage.WebSuccessKey, 0)
	s.setBuntDB(storage.WebErrorKey, 0)
}

func (s *Storage) setBuntDB(key string, count int64) {
	err := s.db.Update(func(tx *buntdb.Tx) error {
		if _, _, err := tx.Set(key, fmt.Sprintf("%d", count), nil); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Println("BuntDB update error:", err.Error())
	}
}

func (s *Storage) getBuntDB(key string, count *int64) {
	err := s.db.View(func(tx *buntdb.Tx) error {
		val, _ := tx.Get(key)
		*count, _ = strconv.ParseInt(val, 10, 64)
		return nil
	})

	if err != nil {
		log.Println("BuntDB get error:", err.Error())
	}
}

// AddTotalCount record push notification count.
func (s *Storage) AddTotalCount(count int64) {
	total := s.GetTotalCount() + count
	s.setBuntDB(storage.TotalCountKey, total)
}

// AddIosSuccess record counts of success iOS push notification.
func (s *Storage) AddIosSuccess(count int64) {
	total := s.GetIosSuccess() + count
	s.setBuntDB(storage.IosSuccessKey, total)
}

// AddIosError record counts of error iOS push notification.
func (s *Storage) AddIosError(count int64) {
	total := s.GetIosError() + count
	s.setBuntDB(storage.IosErrorKey, total)
}

// AddAndroidSuccess record counts of success Android push notification.
func (s *Storage) AddAndroidSuccess(count int64) {
	total := s.GetAndroidSuccess() + count
	s.setBuntDB(storage.AndroidSuccessKey, total)
}

// AddAndroidError record counts of error Android push notification.
func (s *Storage) AddAndroidError(count int64) {
	total := s.GetAndroidError() + count
	s.setBuntDB(storage.AndroidErrorKey, total)
}

// AddWebSuccess record counts of success Web push notification.
func (s *Storage) AddWebSuccess(count int64) {
	total := s.GetWebSuccess() + count
	s.setBuntDB(storage.WebSuccessKey, total)
}

// AddWebError record counts of error Web push notification.
func (s *Storage) AddWebError(count int64) {
	total := s.GetWebError() + count
	s.setBuntDB(storage.WebErrorKey, total)
}

// GetTotalCount show counts of all notification.
func (s *Storage) GetTotalCount() int64 {
	var count int64
	s.getBuntDB(storage.TotalCountKey, &count)

	return count
}

// GetIosSuccess show success counts of iOS notification.
func (s *Storage) GetIosSuccess() int64 {
	var count int64
	s.getBuntDB(storage.IosSuccessKey, &count)

	return count
}

// GetIosError show error counts of iOS notification.
func (s *Storage) GetIosError() int64 {
	var count int64
	s.getBuntDB(storage.IosErrorKey, &count)

	return count
}

// GetAndroidSuccess show success counts of Android notification.
func (s *Storage) GetAndroidSuccess() int64 {
	var count int64
	s.getBuntDB(storage.AndroidSuccessKey, &count)

	return count
}

// GetAndroidError show error counts of Android notification.
func (s *Storage) GetAndroidError() int64 {
	var count int64
	s.getBuntDB(storage.AndroidErrorKey, &count)

	return count
}

// GetWebSuccess show success counts of Web notification.
func (s *Storage) GetWebSuccess() int64 {
	var count int64
	s.getBuntDB(storage.WebSuccessKey, &count)

	return count
}

// GetWebError show error counts of Web notification.
func (s *Storage) GetWebError() int64 {
	var count int64
	s.getBuntDB(storage.WebErrorKey, &count)

	return count
}
