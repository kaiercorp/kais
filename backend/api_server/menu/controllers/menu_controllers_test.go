package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"api_server/logger"
	"api_server/menu/dto"
)

type MockMenuService struct {
	mock.Mock
}

func (this *MockMenuService) ViewMenus(group int) ([]dto.MenuDTO, *logger.Report) {
	args := this.Called(group)

	return args.Get(0).([]dto.MenuDTO), nil
}

func (this *MockMenuService) EditMenus(menus []dto.MenuDTO) *logger.Report {
	args := this.Called(menus)

	return args.Get(0).(*logger.Report)
}

type MockUserService struct {
	mock.Mock
}

func (this *MockUserService) GetUserGroup(token string) int {
	args := this.Called(token)

	return args.Int(0)
}

type MenuControllerTestSuite struct {
	suite.Suite
}

func (suite *MenuControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	logger.InitLogger("", "/workspace/log/test.log")
}

func (suite *MenuControllerTestSuite) TestGetMenus() {
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
}

//func(suite *MenuControllerTestSuite) TestSetMenus(t *testing.T) {
//
//}

func TestMenuControllerTestSuite(t *testing.T) {
	suite.Run(t, new(MenuControllerTestSuite))
}
