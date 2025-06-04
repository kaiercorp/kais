package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	docs "api_server/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	auth_router "api_server/auth/router"
	config_router "api_server/configuration/router"
	dataset_router "api_server/dataset/router"
	device_router "api_server/device/router"
	download_router "api_server/download/router"
	project_router "api_server/project/router"
	tapi_router "api_server/tapi/router"
	task_router "api_server/task/router"
	user_router "api_server/user/router"
	ws_router "api_server/websocket/router"

	// modeling_router "api_server/modeling/router"
	// system_monitor_router "api_server/system/router"

	config_service "api_server/configuration/service"
	dataset_module "api_server/dataset/modules"
	dataset_repo "api_server/dataset/repository"
	device_service "api_server/device/service"
	"api_server/logger"
	task_service "api_server/task/service"
	"api_server/utils"
)

// @title KAI.S API document
// @version 2.0.0
// @BasePath /api
func main() {
	fileconfig := utils.ReadConfigFromFile()

	println("Starting KAI.S...")
	utils.InitDBMS()

	config_init := config_service.NewInit()
	config_init.Init(utils.SW_VERSION)

	cf := config_service.NewStatic()

	log_file := cf.Get("PATH_LOG_DIR") + "/app.log"
	logger.InitLogger(cf.Get("LOG_LEVEL"), log_file)

	logger.Info("PORT:", fileconfig.KAISPORT, " DB:", fileconfig.DBIP, ":", fileconfig.DBPORT, "/", fileconfig.DB, " BROKER:", fileconfig.BROKERIP)

	// Watch Datasets
	datasetDAO := dataset_repo.NewDatasetDAO()
	datasetWatcher := dataset_module.NewDatasetWatcher(dataset_module.NewDatasetValidator(datasetDAO), dataset_module.NewDatasetAnalyzer(datasetDAO), datasetDAO, dataset_repo.NewDatasetRootDAO(datasetDAO))
	go datasetWatcher.WatchDataset()
	// Check GPU Nodes
	device_service.GatherDevicesInfo()
	// change trial state to cancel
	// trial_service.InitializeTrials()
	// Watch Task Queue
	scheduler := task_service.NewTaskScheduler()
	go scheduler.WatchTasks()

	gomode := os.Args[0]
	if strings.Contains(gomode, "Temp\\go-build") || strings.Contains(gomode, "tmp/go-build") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	docs.SwaggerInfo.BasePath = utils.API_BASE_URL_V1
	r := gin.New()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Use(gin.Recovery())
	r.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers, Authorization"},
			MaxAge:       12 * time.Hour,
		}))

	r.Use(requestid.New())
	r.Use(static.ServeRoot("/static", cf.Get("PATH_STATIC_TEST")))
	curPath, _ := os.Getwd()
	r.Use(static.Serve("/", static.LocalFile(filepath.Join(curPath, "frontend"), true)))

	auth_router.InitAuthRouter(r)
	config_router.InitRouter(r)
	dataset_router.InitRouter(r)
	project_router.InitRouter(r)
	task_router.InitRouter(r)
	device_router.InitRouter(r)
	download_router.InitRouter(r)
	ws_router.InitWebSocketRouter(r)
	user_router.InitRouter(r)
	tapi_router.InitRouter(r)

	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	logger.Info("Startup KAI.S " + utils.SW_VERSION)
	_ = r.Run(fileconfig.KAISPORT)
}
