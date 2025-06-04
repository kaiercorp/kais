package router

import (
	"strconv"

	"github.com/gin-gonic/gin"

	svc "api_server/download/service"
	"api_server/logger"
)

type DownloadController struct {
	svc *svc.DownloadService
}

func New(downloadService *svc.DownloadService) *DownloadController {
	return &DownloadController{
		svc: downloadService,
	}
}

func (ctrl *DownloadController) DownloadModel(c *gin.Context) {
	logger.ApiRequest(c)

	modelingID, err := strconv.Atoi(c.Param("modeling_id"))
	if err != nil {
		r := logger.CreateReport(&logger.CODE_REQUEST, err)
		logger.ApiResponse(c, r, nil)
		return
	}

	modelName := c.Param("model_name")
	if modelName == "" {
		logger.Debug("model name is empty")
		return
	}

	zipPath, report := ctrl.svc.GetModelZip(c.Request.Context(), modelingID, modelName)
	if report != nil {
		logger.ApiResponse(c, report, nil)
		return
	}

	logger.ApiResponseWithZipFile(c, zipPath)
}
