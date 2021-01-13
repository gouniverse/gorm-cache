package cache

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gouniverse/uid"
	"gorm.io/gorm"
)

// Cache type
type Cache struct {
	ID        string     `gorm:"type:varchar(40);column:id;primary_key;"`
	Key       string     `gorm:"type:varchar(40);column:key;DEFAULT NULL;"`
	Value     string     `gorm:"type:longtext;column:value;"`
	ExpiresAt *time.Time `gorm:"type:datetime;olumn:expores_at;DEFAULT NULL;"`
	CreatedAt time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt *time.Time `gorm:"type:datetime;olumn:deleted_at;DEFAULT NULL;"`
}

// TableName the name of the Cache table
func (Cache) TableName() string {
	return "snv_caches_cache"
}

// BeforeCreate adds UID to model
func (c *Cache) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.HumanUid()
	c.ID = uuid
	return nil
}

// FindByKey finds a cache by key
func FindByKey(db *gorm.DB, key string) *Cache {
	cache := &Cache{}

	result := db.Where("`key` = ?", key).First(&cache)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		log.Panic(result.Error)
	}

	return cache
}

// Get gets a key from cache
func Get(db *gorm.DB, key string, valueDefault string) string {
	cache := FindByKey(db, key)

	if cache != nil {
		return cache.Value
	}

	return valueDefault
}

// Set sets a key in cache
func Set(db *gorm.DB, key string, value string, seconds int64) bool {
	cache := FindByKey(db, key)
	expiresAt := time.Now().Add(time.Second * time.Duration(seconds))

	if cache != nil {
		cache.Value = value
		cache.ExpiresAt = &expiresAt
		//dbResult := GetDb().Table(User).Where("`key` = ?", key).Update(&cache)
		dbResult := db.Save(&cache)
		if dbResult != nil {
			return false
		}
		return true
	}

	var newCache = Cache{Key: key, Value: value, ExpiresAt: &expiresAt}

	dbResult := db.Create(&newCache)

	if dbResult.Error != nil {
		return false
	}

	return true
}

// ExpireCacheGoroutine - soft deletes expired cache
func ExpireCacheGoroutine(db gorm.DB) {
	i := 0
	for {
		i++
		fmt.Println("Cleaning expired cache...")
		db.Where("`expires_at` < ?", time.Now()).Delete(Cache{})
		time.Sleep(60 * time.Second) // Every minute
	}
}
