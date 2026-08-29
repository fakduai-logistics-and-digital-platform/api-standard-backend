package services_test

import (
	"errors"
	"testing"
	"uuid"

	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/models"
	"fdlp-standard-api/internal/services"
	"fdlp-standard-api/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetById(id string) (*models.User, error) {
	args := m.Called(id)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindAllByFilterAndPage(filter map[string]any, page, pageSize int) ([]models.User, int64, int, error) {
	args := m.Called(filter, page, pageSize)
	if res := args.Get(0); res != nil {
		return res.([]models.User), args.Get(1).(int64), args.Int(2), args.Error(3)
	}
	return nil, 0, 0, args.Error(3)
}

func TestUserService_GetUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, "secret", nil)

	mockUser := &models.User{
		UserID:   types.UUID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		Username: "testuser",
		Email:    "test@example.com",
		Role: models.Role{
			Name:        "admin",
			Description: "Admin role",
		},
	}

	mockUserRepo.On("GetById", "user-1").Return(mockUser, nil)

	res, err := userService.GetUser("user-1")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", res.Username)
	assert.Equal(t, "test@example.com", res.Email)
	assert.Equal(t, "admin", res.RoleName)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, "secret", nil)

	mockUserRepo.On("GetById", "not-found").Return(nil, errors.New("user not found"))

	res, err := userService.GetUser("not-found")
	assert.Error(t, err)
	assert.Nil(t, res)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetUsersWithPagination(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, "secret", nil)

	filter := map[string]any{"username": "test"}
	expectedUsers := []models.User{
		{
			Username: "test",
			Email:    "test@example.com",
		},
	}

	mockUserRepo.On("FindAllByFilterAndPage", filter, 1, 10).Return(expectedUsers, int64(1), 1, nil)

	users, totalRows, totalPages, err := userService.GetUsersWithPagination(filter, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, int64(1), totalRows)
	assert.Equal(t, 1, totalPages)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_LoginUser_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, "secret-key-12345", nil)

	mockUser := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "$2a$10$abcdefghijklmnopqrstuuNOPQRSTUVWXYZ1234567890abcdef", // will test login
		Role: models.Role{
			Name:        "admin",
			Description: "Admin role",
		},
	}

	mockUserRepo.On("GetByEmail", "notfound@example.com").Return(nil, errors.New("user not found"))

	token, err := userService.LoginUser(dto.LoginUserRequestBody{
		Email:    "notfound@example.com",
		Password: "password123",
	})
	assert.Error(t, err)
	assert.Nil(t, token)
	mockUserRepo.AssertExpectations(t)
	_ = mockUser
}
