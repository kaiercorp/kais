package router

import (
	"api_server/logger"
	repo "api_server/project/repository"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockService struct {
	mock.Mock
}

func (msvc *MockService) Create(req repo.ProjectDTO) (*repo.ProjectPages, *logger.Report) {
	args := msvc.Called(req)

	return args.Get(0).(*repo.ProjectPages), nil
}

func (msvc *MockService) Read(page int) (*repo.ProjectPages, *logger.Report) {
	args := msvc.Called(page)

	return args.Get(0).(*repo.ProjectPages), nil
}

func (msvc *MockService) Edit(req repo.ProjectDTO) (*repo.ProjectDTO, *logger.Report) {
	args := msvc.Called(req)

	return args.Get(0).(*repo.ProjectDTO), nil
}

func (msvc *MockService) Delete(project_id int) (*repo.ProjectPages, *logger.Report) {
	args := msvc.Called(project_id)

	return args.Get(0).(*repo.ProjectPages), nil
}

type ControllerTestSuite struct {
	suite.Suite
}

func (suite *ControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	logger.InitLogger("", "/workspace/log/test.log")
}

func (suite *ControllerTestSuite) TestGetProjects() {
	// TODO
	// var mockProjects repo.ProjectPages

	// mockService := new(MockService)
	// mockService.On("Read", 0).Return(mockProjects, nil)

	// mockController := New(mockService)

	// router := gin.Default()
	// router.GET("/project", mockController.GetByPages)

	/*
		var mockMenu []dto.MenuDTO
		mockMenu = append(mockMenu, dto.MenuDTO{
			Key: "key",
		})

		mockUserService := new(MockUserService)
		mockUserService.On("GetUserGroup", "").Return(2)

		mockMenuService := new(MockMenuService)
		mockMenuService.On("ViewMenus", 2).Return(mockMenu, nil)

		menuController := NewMenuController(mockMenuService, mockUserService)

		router := gin.Default()
		router.GET("/menu", menuController.GetMenus)

		req, _ := http.NewRequest("GET", "/menu", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)
		jsonMap := make(map[string]interface{})
		_ = json.Unmarshal(w.Body.Bytes(), &jsonMap)
		assert.Equal(suite.T(), mockMenu[0].Key, jsonMap["data"].([]interface{})[0].(map[string]interface{})["key"])

		mockMenuService.AssertExpectations(suite.T())
		mockUserService.AssertExpectations(suite.T())
	*/
}

func TestControllerTestSuite(t *testing.T) {
	// TODO
	// suite.Run(t, new(ControllerTestSuite))
}
