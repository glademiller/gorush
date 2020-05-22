package badger

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/netscale-technologies/gorush/config"
	"github.com/netscale-technologies/gorush/storage"

	"github.com/appleboy/com/convert"
	"github.com/dgraph-io/badger/v2"
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
	opts   badger.Options
	name   string
	db     *badger.DB
}

// Init client storage.
func (s *Storage) Init() error {
	var err error
	s.name = "badger"
	dbPath := s.config.Stat.BadgerDB.Path
	if dbPath == "" {
		dbPath = os.TempDir() + "badger"
	}
	s.opts = badger.DefaultOptions(dbPath)

	s.db, err = badger.Open(s.opts)

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
	s.setBadger(storage.TotalCountKey, 0)
	s.setBadger(storage.IosSuccessKey, 0)
	s.setBadger(storage.IosErrorKey, 0)
	s.setBadger(storage.AndroidSuccessKey, 0)
	s.setBadger(storage.AndroidErrorKey, 0)
}

func (s *Storage) setBadger(key string, count int64) {
	err := s.db.Update(func(txn *badger.Txn) error {
		value := convert.ToString(count).(string)
		return txn.Set([]byte(key), []byte(value))
	})

	if err != nil {
		log.Println(s.name, "update error:", err.Error())
	}
}

func (s *Storage) getBadger(key string, count *int64) {
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		dst := []byte{}
		val, err := item.ValueCopy(dst)
		if err != nil {
			return err
		}

		i, err := strconv.ParseInt(fmt.Sprintf("%s", val), 10, 64)
		if err != nil {
			return err
		}

		*count = i

		return nil
	})

	if err != nil {
		log.Println(s.name, "get error:", err.Error())
	}
}

// AddTotalCount record push notification count.
func (s *Storage) AddTotalCount(count int64) {
	total := s.GetTotalCount() + count
	s.setBadger(storage.TotalCountKey, total)
}

// AddIosSuccess record counts of success iOS push notification.
func (s *Storage) AddIosSuccess(count int64) {
	total := s.GetIosSuccess() + count
	s.setBadger(storage.IosSuccessKey, total)
}

// AddIosError record counts of error iOS push notification.
func (s *Storage) AddIosError(count int64) {
	total := s.GetIosError() + count
	s.setBadger(storage.IosErrorKey, total)
}

// AddAndroidSuccess record counts of success Android push notification.
func (s *Storage) AddAndroidSuccess(count int64) {
	total := s.GetAndroidSuccess() + count
	s.setBadger(storage.AndroidSuccessKey, total)
}

// AddAndroidError record counts of error Android push notification.
func (s *Storage) AddAndroidError(count int64) {
	total := s.GetAndroidError() + count
	s.setBadger(storage.AndroidErrorKey, total)
}

// AddWebSuccess record counts of success web push notification.
func (s *Storage) AddWebSuccess(count int64) {
	total := s.GetWebSuccess() + count
	s.setBadger(storage.WebSuccessKey, total)
}

// AddWebError record counts of error web push notification.
func (s *Storage) AddWebError(count int64) {
	total := s.GetWebError() + count
	s.setBadger(storage.WebErrorKey, total)
}

// GetTotalCount show counts of all notification.
func (s *Storage) GetTotalCount() int64 {
	var count int64
	s.getBadger(storage.TotalCountKey, &count)

	return count
}

// GetIosSuccess show success counts of iOS notification.
func (s *Storage) GetIosSuccess() int64 {
	var count int64
	s.getBadger(storage.IosSuccessKey, &count)

	return count
}

// GetIosError show error counts of iOS notification.
func (s *Storage) GetIosError() int64 {
	var count int64
	s.getBadger(storage.IosErrorKey, &count)

	return count
}

// GetAndroidSuccess show success counts of Android notification.
func (s *Storage) GetAndroidSuccess() int64 {
	var count int64
	s.getBadger(storage.AndroidSuccessKey, &count)

	return count
}

// GetAndroidError show error counts of Android notification.
func (s *Storage) GetAndroidError() int64 {
	var count int64
	s.getBadger(storage.AndroidErrorKey, &count)

	return count
}

// GetWebSuccess show success counts of web notification.
func (s *Storage) GetWebSuccess() int64 {
	var count int64
	s.getBadger(storage.WebSuccessKey, &count)

	return count
}

// GetWebError show error counts of web notification.
func (s *Storage) GetWebError() int64 {
	var count int64
	s.getBadger(storage.WebErrorKey, &count)

	return count
}
