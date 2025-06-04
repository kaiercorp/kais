package service

import (
	"context"

	"api_server/dataset/modules"
	repo "api_server/dataset/repository"
	"api_server/logger"
)

type IDatasetRootService interface {
	ViewDatasetrootAll() ([]*repo.DatasetRootDTO, *logger.Report)
	ViewDatasetrootAllForAPI() ([]*repo.DatasetRootDTO, *logger.Report)
	ViewDatasetrootActive() ([]*repo.DatasetRootDTO, *logger.Report)
	AddDatasetroot(dr repo.DatasetRootDTO) (*repo.DatasetRootDTO, *logger.Report)
	EditDatasetroot(dr repo.DatasetRootDTO) (*repo.DatasetRootDTO, *logger.Report)
	RemoveDatasetroot(id int) *logger.Report
}

type DatasetRootService struct {
	ctx            context.Context
	datasetWatcher modules.DatasetWatcherInterface
	datasetRootDAO repo.DatasetRootDAOInterface
	datasetDAO     repo.DatasetDAOInterface
}

var datasetRootServiceInstance *DatasetRootService

func NewDatasetRootService(datasetWatcher modules.DatasetWatcherInterface, datasetRootDAO repo.DatasetRootDAOInterface, datasetDAO repo.DatasetDAOInterface) *DatasetRootService {
	if datasetRootServiceInstance == nil {
		datasetRootServiceInstance = &DatasetRootService{
			ctx:            context.Background(),
			datasetWatcher: datasetWatcher,
			datasetRootDAO: datasetRootDAO,
			datasetDAO:     datasetDAO,
		}
	}

	return datasetRootServiceInstance
}

func (svc *DatasetRootService) ViewDatasetrootAll() ([]*repo.DatasetRootDTO, *logger.Report) {
	if drs, err := svc.datasetRootDAO.SelectAll(svc.ctx); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertDatasetrootEntsToDTOs(drs), nil
	}
}

func (svc *DatasetRootService) ViewDatasetrootAllForAPI() ([]*repo.DatasetRootDTO, *logger.Report) {
	if drs, err := svc.datasetRootDAO.SelectAll(svc.ctx); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertDatasetrootEntsToDTOs(drs), nil
	}
}

func (svc *DatasetRootService) ViewDatasetrootActive() ([]*repo.DatasetRootDTO, *logger.Report) {
	if drs, err := svc.datasetRootDAO.SelectActive(svc.ctx); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	} else {
		return repo.ConvertDatasetrootEntsToDTOs(drs), nil
	}
}

func (svc *DatasetRootService) AddDatasetroot(dr repo.DatasetRootDTO) (*repo.DatasetRootDTO, *logger.Report) {
	if inserted, err := svc.datasetRootDAO.InsertOne(svc.ctx, dr); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_INSERT, err)
	} else {
		svc.datasetWatcher.DetectDatasetModification()
		return repo.ConvertDatasetrootEntToDTO(inserted), nil
	}
}

func (svc *DatasetRootService) EditDatasetroot(dr repo.DatasetRootDTO) (*repo.DatasetRootDTO, *logger.Report) {
	if edited, err := svc.datasetRootDAO.UpdateOne(svc.ctx, dr); err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	} else {
		svc.datasetWatcher.DetectDatasetModification()
		return repo.ConvertDatasetrootEntToDTO(edited), nil
	}
}

func (svc *DatasetRootService) RemoveDatasetroot(id int) *logger.Report {
	if err := svc.datasetRootDAO.DeleteOne(svc.ctx, id); err != nil {
		return logger.CreateReport(&logger.CODE_DB_DELETE, err)
	}

	return svc.datasetDAO.DeleteDatasetByDRID(svc.ctx, id)
}
