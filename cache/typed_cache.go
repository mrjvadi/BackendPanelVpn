// Filename: cache/typed_cache.go
package cache

//type TypedCache[T any] struct {
//	store CacheInterface
//}
//
//func NewTypedCache[T any](store CacheInterface) *TypedCache[T] {
//	return &TypedCache[T]{store: store}
//}
//
//func (c *TypedCache[T]) Get(ctx context.Context, key string) (T, error) {
//	var zero T
//	data, err := c.store.Get(ctx, key)
//	if err != nil {
//		return zero, err
//	}
//	var value T
//	if err := json.Unmarshal(data, &value); err != nil {
//		return zero, err
//	}
//	return value, nil
//}
//
//func (c *TypedCache[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
//	data, err := json.Marshal(value)
//	if err != nil {
//		return err
//	}
//	return c.store.Set(ctx, key, data, ttl)
//}
