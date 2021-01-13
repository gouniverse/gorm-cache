# GORM Cache

Cache model for GORM


- Create cache table
```
db().AutoMigrate(cache.Cache{})
```

- Set cache to remove expired records
```
cache.ExpireCacheGoroutine(db())
```

- Set a key-value pair in cache
```
cache.Set(db(), "key", "value")
```

- Get a key's value from cache, or default if not exists
```
key := cache.Get(db(), "key", "default")
```

