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

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetById(id string) (*models.Role, error) {
	args := m.Called(id)
	if res := args.Get(0); res != nil {
		return res.(*models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRoleRepository) GetByName(name string) (*models.Role, error) {
	args := m.Called(name)
	if res := args.Get(0); res != nil {
		return res.(*models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRoleRepository) Create(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func TestRoleService_GetRole(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	expectedRole := &models.Role{
		RoleID:      types.UUID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		Name:        "admin",
		Description: "Admin role",
	}

	mockRepo.On("GetById", "role-123").Return(expectedRole, nil)

	role, err := roleService.GetRole("role-123")
	assert.NoError(t, err)
	assert.Equal(t, expectedRole, role)
	mockRepo.AssertExpectations(t)
}

func TestRoleService_GetRole_Error(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	mockRepo.On("GetById", "not-found").Return(nil, errors.New("not found"))

	role, err := roleService.GetRole("not-found")
	assert.Error(t, err)
	assert.Nil(t, role)
	mockRepo.AssertExpectations(t)
}

func TestRoleService_CreateRole(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	roleService := services.NewRoleService(mockRepo)

	input := dto.CreateRoleRequestBody{
		Name:        "member",
		Description: "Member role",
	}

	mockRepo.On("Create", mock.MatchedBy(func(r *models.Role) bool {
		return r.Name == "member" && r.Description == "Member role"
	})).Return(nil)

	role, err := roleService.CreateRole(input)
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "member", role.Name)
	assert.Equal(t, "Member role", role.Description)
	mockRepo.AssertExpectations(t)
}
