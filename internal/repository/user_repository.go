package repository

import (
	"errors"
	"learn/internal/models"
	"sync"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
}

// InMemoryUserRepository implements UserRepository with an in-memory store
type InMemoryUserRepository struct {
	users map[string]*models.User
	mutex sync.RWMutex
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*models.User),
	}
}

// Create adds a new user to the repository
func (r *InMemoryUserRepository) Create(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if user with same username or email already exists
	for _, existingUser := range r.users {
		if existingUser.Username == user.Username {
			return ErrUserAlreadyExists
		}
		if existingUser.Email == user.Email {
			return ErrUserAlreadyExists
		}
	}

	r.users[user.ID] = user
	return nil
}

// GetByID retrieves a user by ID
func (r *InMemoryUserRepository) GetByID(id string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetByUsername retrieves a user by username
func (r *InMemoryUserRepository) GetByUsername(username string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// GetByEmail retrieves a user by email
func (r *InMemoryUserRepository) GetByEmail(email string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}
	r.users[user.ID] = user
	return nil
}

// Delete removes a user by ID
func (r *InMemoryUserRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}