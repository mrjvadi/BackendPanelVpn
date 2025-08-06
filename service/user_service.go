// Filename: service/user_service.go
package service

//
//type UserService struct {
//	userRepo  storage.UserRepository
//	userCache *cache.TypedCache[models.User]
//	bus       *events.Bus
//}
//
//func NewUserService(userRepo storage.UserRepository, userCache *cache.TypedCache[models.User], bus *events.Bus) *UserService {
//	return &UserService{
//		userRepo:  userRepo,
//		userCache: userCache,
//		bus:       bus,
//	}
//}
//
//func (s *UserService) CreateUser(username string) (*models.User, error) {
//	user := &models.User{
//		Username: username,
//		UUID:     uuid.New(),
//		Status:   "active",
//	}
//
//	if err := s.userRepo.Create(user); err != nil {
//		return nil, err
//	}
//
//	s.bus.Publish(events.Event{Name: "service:user_created", Payload: *user})
//
//	cacheKey := fmt.Sprintf("user:%d", user.ID)
//	_ = s.userCache.Set(context.Background(), cacheKey, *user, 1*time.Hour)
//
//	return user, nil
//}
//
//func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
//	cacheKey := fmt.Sprintf("user:%d", id)
//
//	cachedUser, err := s.userCache.Get(ctx, cacheKey)
//	if err == nil {
//		return &cachedUser, nil
//	}
//
//	if err != redis.Nil {
//		fmt.Printf("[Service] Cache error: %v. Fetching from DB.\n", err)
//	}
//
//	user, err := s.userRepo.FindByID(id)
//	if err != nil {
//		return nil, err
//	}
//
//	_ = s.userCache.Set(ctx, cacheKey, *user, 1*time.Hour)
//
//	return user, nil
//}
