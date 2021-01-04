package cache

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/gouniverse/uid"
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
func (c *Cache) BeforeCreate(scope *gorm.Scope) (err error) {
	uuid := uid.NanoUid()
	scope.SetColumn("ID", uuid)
	return nil
}

// CacheFindByKey finds a cache by key
func CacheFindByKey(key string) *Cache {
	cache := &Cache{}
	if GetDb().Where("`key` = ?", key).First(&cache).RecordNotFound() {
		return nil
	}

	return cache
}

// CacheGet gets a key from cache
func CacheGet(key string, valueDefault string) string {
	cache := CacheFindByKey(key)

	if cache != nil {
		return cache.Value
	}

	return valueDefault
}

// CacheSet sets a key in cache
func CacheSet(key string, value string, seconds int64) bool {
	cache := CacheFindByKey(key)
	expiresAt := time.Now().Add(time.Second * time.Duration(seconds))

	if cache != nil {
		cache.Value = value
		cache.ExpiresAt = &expiresAt
		//dbResult := GetDb().Table(User).Where("`key` = ?", key).Update(&cache)
		dbResult := GetDb().Save(&cache)
		if dbResult != nil {
			return false
		}
		return true
	}

	var newCache = Cache{Key: key, Value: value, ExpiresAt: &expiresAt}

	dbResult := GetDb().Create(&newCache)

	if dbResult.Error != nil {
		return false
	}

	return true
}

// CacheExpireJobGoroutine - soft deletes expired cache
func CacheExpireJobGoroutine() {
	i := 0
	for {
		i++
		fmt.Println("Cleaning expired cache...")
		GetDb().Where("`expires_at` < ?", time.Now()).Delete(Cache{})
		time.Sleep(60 * time.Second) // Every minute
	}
}
