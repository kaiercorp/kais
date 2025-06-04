package repository

import (
	"context"

	"api_server/ent"
	"api_server/ent/dataset"
	"api_server/ent/datasetroot"
	"api_server/utils"
)

type DatasetRootDAOInterface interface {
	SelectActive(ctx context.Context) ([]*ent.DatasetRoot, error)
	SelectAll(ctx context.Context) ([]*ent.DatasetRoot, error)
	InsertOne(ctx context.Context, dr DatasetRootDTO) (*ent.DatasetRoot, error)
	UpdateOne(ctx context.Context, dr DatasetRootDTO) (*ent.DatasetRoot, error)
	DeleteOne(ctx context.Context, id int) error
}

type DatasetRootDAO struct {
	entClient  *ent.Client
	datasetDAO DatasetDAOInterface
}

var datasetRootDAOInstance *DatasetRootDAO

func NewDatasetRootDAO(datasetDAO DatasetDAOInterface) *DatasetRootDAO {
	if datasetRootDAOInstance == nil {
		datasetRootDAOInstance = &DatasetRootDAO{
			entClient:  utils.GetEntClient(),
			datasetDAO: datasetDAO,
		}
	}

	return datasetRootDAOInstance
}

func (dao *DatasetRootDAO) SelectActive(ctx context.Context) ([]*ent.DatasetRoot, error) {
	return dao.entClient.DatasetRoot.
		Query().
		Where(datasetroot.IsUse(true)).
		All(ctx)
}

func (dao *DatasetRootDAO) SelectAll(ctx context.Context) ([]*ent.DatasetRoot, error) {
	if drs, err := dao.entClient.DatasetRoot.
		Query().
		Order(datasetroot.ByID()).
		All(ctx); err != nil {
		return nil, err
	} else {
		for _, dr := range drs {
			dr.Edges.Datasets = dao.datasetDAO.SelectDatasetsByDRID(ctx, dr.ID)
		}

		return drs, nil
	}
}

func (dao *DatasetRootDAO) InsertOne(ctx context.Context, req DatasetRootDTO) (*ent.DatasetRoot, error) {
	return dao.entClient.DatasetRoot.Create().
		SetName(req.Name).
		SetPath(req.Path).
		SetIsUse(req.IsUse).
		Save(ctx)
}

func (dao *DatasetRootDAO) UpdateOne(ctx context.Context, req DatasetRootDTO) (*ent.DatasetRoot, error) {
	_, _ = dao.entClient.Dataset.
		Delete().
		Where(
			dataset.DrID(req.ID),
		).
		Exec(ctx)

	return dao.entClient.DatasetRoot.UpdateOneID(req.ID).
		Where(datasetroot.ID(req.ID)).
		SetName(req.Name).
		SetPath(req.Path).
		SetIsUse(req.IsUse).
		Save(ctx)
}

func (dao *DatasetRootDAO) DeleteOne(ctx context.Context, id int) error {
	return dao.entClient.DatasetRoot.DeleteOneID(id).Exec(ctx)
}
