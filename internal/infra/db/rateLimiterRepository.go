package db

type Cache struct {
	client *redis.Client
}

func NewCache() *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &Cache{client: client}
)