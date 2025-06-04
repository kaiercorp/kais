package router

import (
	"github.com/gin-gonic/gin"

	repo_dataset "api_server/dataset/repository"
	repo_device "api_server/device/repository"
	"api_server/project/repository"
	"api_server/project/service"
	repo_task "api_server/task/repository"
	service_task "api_server/task/service"
	"api_server/utils"
)

func InitRouter(r *gin.Engine) {
	project_dao := repository.New()
	user_proejct_dao := repository.NewUserProject()
	task_dao := repo_task.NewTaskDAO()
	modeling_dao := repo_task.NewModelingDAO()
	dataset_dao := repo_dataset.NewDatasetDAO()
	device_dao := repo_device.New()
	task_service := service_task.NewTaskService(task_dao, modeling_dao, dataset_dao, device_dao)
	usr_proj_svc := service.NewUserProject(user_proejct_dao)
	svc := service.New(project_dao, task_service, usr_proj_svc)
	controller := New(svc)
	controllerUP := NewUserProjectController(usr_proj_svc)

	apiRouter := r.Group(utils.API_BASE_URL_V1 + "/project")
	{
		// apiRouter.Use(utils.JWTAuthMiddleware())
		apiRouter.POST("", utils.JWTAuthMiddleware(), controller.CreateOne)
		apiRouter.GET("", controller.GetByPages)
		apiRouter.GET("/logged", utils.JWTAuthMiddleware(), controller.GetByPagesWithUser)
		apiRouter.PUT("", controller.UpdateById)
		apiRouter.DELETE("/:id", controller.DeleteById)
	}
	apiRouterUP := r.Group(utils.API_BASE_URL_V1 + "/user-project")
	{
		// Create a new user-project relationship
		apiRouterUP.POST("", controllerUP.CreateUserProject)
		// Get all projects for a specific user
		apiRouterUP.GET("/user/projects", controllerUP.GetAllProjectsByUser)
		// Get all users for a specific project
		apiRouterUP.GET("/project/users", controllerUP.GetAllUsersByProject)
		// Delete specific user from a specific project
		apiRouterUP.DELETE("", controllerUP.DeleteUserFromProject)
		// Delete all projects associated with a specific user
		apiRouterUP.DELETE("/user/projects", controllerUP.DeleteAllProjectsByUser)
		// Delete all users associated with a specific project
		apiRouterUP.DELETE("/project/users", controllerUP.DeleteAllUsersByProject)
	}

}
