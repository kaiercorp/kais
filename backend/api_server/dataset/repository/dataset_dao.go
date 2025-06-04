package repository

import (
	"context"
	"math"
	"slices"
	"time"

	"api_server/ent"
	"api_server/ent/dataset"
	"api_server/ent/datasetroot"
	"api_server/logger"
	"api_server/utils"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
)

type DatasetDAOInterface interface {
	SelectDatasets(ctx context.Context, datasetType []string, page int) ([]*ent.Dataset, int, bool, int, error)
	SelectTestableDatasets(ctx context.Context, datasetType []string, page int) ([]*ent.Dataset, int, *logger.Report)
	SelectDataPathByDataSetId(ctx context.Context, id int) (string, *logger.Report)
	SelectDataSetByParentID(ctx context.Context, parent_id int) ([]*ent.Dataset, *logger.Report)
	SelectDataSetByID(ctx context.Context, id int) ([]*ent.Dataset, *logger.Report)
	SelectDatasetForAPI(ctx context.Context, parent_id int, data_type string) ([]*ent.Dataset, *logger.Report)
	SelectDataSetByName(ctx context.Context, name string) ([]*ent.Dataset, *logger.Report)
	SelectByPath(ctx context.Context, path string) (*ent.Dataset, *logger.Report)
	SelectDatasetsByDRID(ctx context.Context, dr_id int) []*ent.Dataset
	SelectStatistics(ctx context.Context, id int) (*ent.Dataset, *logger.Report)
	InsertOne(ctx context.Context, ds DatasetDTO) *ent.Dataset
	UpdateDatasetDeleted(ctx context.Context, dataset_id int)
	UpdateDatasetExist(ctx context.Context, dataset_id int)
	UpdateValidation(ctx context.Context, ds DatasetDTO)
	UpdateStat(ctx context.Context, id int, stat string) *logger.Report
	UpdateStatPath(ctx context.Context, id int, stat string) *logger.Report
	DeleteDataset(ctx context.Context, id int) *logger.Report
	DeleteDatasetByDRID(ctx context.Context, dr_id int) *logger.Report
}

type DatasetDAO struct {
	entClient *ent.Client
}

var datasetDAOInstance *DatasetDAO

func NewDatasetDAO() *DatasetDAO {
	if datasetDAOInstance == nil {
		datasetDAOInstance = &DatasetDAO{
			entClient: utils.GetEntClient(),
		}
	}

	return datasetDAOInstance
}

func buildDatasetFilter(datasetTypes []string) func(s *sql.Selector) {
	return func(s *sql.Selector) {
		// "all"이 포함되어 있는지 확인
		hasAll := false
		for _, datasetType := range datasetTypes {
			if datasetType == "all" {
				hasAll = true
				break
			}
		}

		dr := sql.Table(datasetroot.Table)
		s.Join(dr).On(s.C(dataset.FieldDrID), dr.C(datasetroot.FieldID))

		// 기본 조건: 사용 중이고 삭제되지 않은 데이터셋
		baseCondition := sql.And(
			sql.IsTrue(dr.C(datasetroot.FieldIsUse)),
			sql.IsFalse(dataset.FieldIsDeleted),
			sql.EQ(s.C(dataset.FieldParentID), 0),
		)

		// "all"이 포함되어 있으면 datasetType 조건을 추가하지 않음
		if hasAll {
			s.Where(baseCondition)
		} else {
			// datasetType 조건 추가
			orConditions := make([]*sql.Predicate, 0, len(datasetTypes))
			for _, datasetType := range datasetTypes {
				orConditions = append(orConditions, sqljson.ValueContains(dataset.FieldEngine, datasetType))
			}

			s.Where(sql.And(
				baseCondition,
				sql.Or(orConditions...),
			))
		}
	}
}

func (dao *DatasetDAO) SelectDatasets(ctx context.Context, datasetTypes []string, page int) ([]*ent.Dataset, int, bool, int, error) {
	filterFunc := buildDatasetFilter(datasetTypes)
	const pageSize = 25

	// 데이터셋 쿼리
	datasets, err := dao.entClient.Dataset.
		Query().
		Select(
			dataset.FieldID,
			dataset.FieldCreatedAt,
			dataset.FieldEngine,
			dataset.FieldName,
			dataset.FieldDescription,
		).
		Where(filterFunc).
		Order(dataset.ByID()).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		All(ctx)

	if err != nil {
		return nil, 0, false, 0, err
	}

	// 카운트 쿼리
	count := dao.entClient.Dataset.
		Query().
		Where(filterFunc).
		CountX(ctx)

	curPage := int(math.Ceil(float64(count) / float64(pageSize)))
	hasNextPage := curPage > page
	nextPage := page + 1

	return datasets, curPage, hasNextPage, nextPage, nil
}

func (dao *DatasetDAO) SelectTestableDatasets(ctx context.Context, datasetType []string, page int) ([]*ent.Dataset, int, *logger.Report) {
	dataType := "table"
	if slices.Contains(datasetType, utils.JOB_TYPE_VISION_AD) || slices.Contains(datasetType, utils.JOB_TYPE_VISION_CLS_SL) || slices.Contains(datasetType, utils.JOB_TYPE_VISION_CLS_ML) {
		dataType = "image"
	}

	datasets, err := dao.entClient.Dataset.
		Query().
		Select(dataset.FieldID, dataset.FieldCreatedAt, dataset.FieldEngine, dataset.FieldName, dataset.FieldDescription).
		Where(
			func(s *sql.Selector) {
				dr := sql.Table(datasetroot.Table)
				s.Join(dr).On(s.C(dataset.FieldDrID), dr.C(datasetroot.FieldID))

				s.Where(
					sql.And(
						sql.IsTrue(dr.C(datasetroot.FieldIsUse)),
						sql.IsFalse(dataset.FieldIsDeleted),
						sql.EQ(dataset.FieldDataType, dataType),
						sql.EQ(dataset.FieldParentID, 0),
					),
				)
			},
		).
		Order(dataset.ByID()).
		Offset((page - 1) * 25).
		Limit(25).
		All(ctx)

	count := dao.entClient.Dataset.
		Query().
		Where(
			func(s *sql.Selector) {
				dr := sql.Table(datasetroot.Table)
				s.Join(dr).On(s.C(dataset.FieldDrID), dr.C(datasetroot.FieldID))

				s.Where(
					sql.And(
						sql.IsTrue(dr.C(datasetroot.FieldIsUse)),
						sql.IsFalse(dataset.FieldIsDeleted),
						sql.EQ(dataset.FieldDataType, dataType),
						sql.EQ(dataset.FieldParentID, 0),
					),
				)
			},
		).
		CountX(ctx)

	if err != nil {
		return nil, 0, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	return datasets, int(math.Ceil(float64(count) / float64(25))), nil
}

func (dao *DatasetDAO) SelectDataPathByDataSetId(ctx context.Context, id int) (string, *logger.Report) {
	ds, err := dao.entClient.Dataset.
		Query().
		Where(dataset.ID(id)).
		Only(ctx)

	if err != nil {
		return "", logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	return ds.Path, nil
}

func (dao *DatasetDAO) SelectDataSetByParentID(ctx context.Context, parent_id int) ([]*ent.Dataset, *logger.Report) {
	dss, err := dao.entClient.Dataset.
		Query().
		Where(dataset.And(
			dataset.ParentID(parent_id),
			dataset.IsDeleted(false),
		)).
		Order(dataset.ByName(sql.OrderAsc())).
		All(ctx)

	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	return dss, nil
}

func (dao *DatasetDAO) SelectDataSetByID(ctx context.Context, id int) ([]*ent.Dataset, *logger.Report) {
	dss, err := dao.entClient.Dataset.
		Query().
		Where(dataset.ID(id)).
		Order(dataset.ByName(sql.OrderAsc())).
		All(ctx)

	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	return dss, nil
}

func (dao *DatasetDAO) SelectDatasetForAPI(ctx context.Context, parent_id int, data_type string) ([]*ent.Dataset, *logger.Report) {
	dss, err := dao.entClient.Dataset.
		Query().
		Where(func(s *sql.Selector) {
			dr := sql.Table(datasetroot.Table)
			s.Join(dr).On(s.C(dataset.FieldDrID), dr.C(datasetroot.FieldID))
			s.Where(sql.And(
				sql.EQ(s.C(dataset.FieldParentID), parent_id),
				sql.EQ(s.C(dataset.FieldDataType), data_type),
				sql.EQ(s.C(dataset.FieldIsDeleted), sql.False()),
			))
		}).
		Order(dataset.ByName(sql.OrderAsc())).
		All(ctx)

	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	return dss, nil
}

func (dao *DatasetDAO) SelectByPath(ctx context.Context, path string) (*ent.Dataset, *logger.Report) {
	ds, _ := dao.entClient.Dataset.
		Query().
		Where(dataset.Path(path)).
		Only(ctx)

	return ds, nil
}

func (dao *DatasetDAO) SelectDatasetsByDRID(ctx context.Context, dr_id int) []*ent.Dataset {
	dss, err := dao.entClient.Dataset.Query().
		Where(dataset.And(dataset.DrID(dr_id), dataset.ParentID(0), dataset.IsDeleted(false))).
		Order(dataset.ByName(sql.OrderAsc())).
		All(ctx)

	if err != nil {
		logger.CreateReport(&logger.CODE_DB_SELECT, err)
		return nil
	}

	return dss
}

func (dao *DatasetDAO) SelectStatistics(ctx context.Context, id int) (*ent.Dataset, *logger.Report) {
	ds, err := dao.entClient.Dataset.Query().Where(dataset.ID(id)).Only(ctx)

	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	return ds, nil
}

func (dao *DatasetDAO) InsertOne(ctx context.Context, ds DatasetDTO) *ent.Dataset {
	inserted, err := dao.entClient.Dataset.Create().
		SetName(ds.Name).
		SetParentID(ds.ParentID).
		SetDescription(ds.Description).
		SetPath(ds.Path).
		SetIsValid(ds.IsValid).
		SetIsTrainable(ds.IsTrainable).
		SetIsTestable(ds.IsTestable).
		SetIsLeaf(ds.IsLeaf).
		SetIsUse(ds.IsUse).
		SetIsDeleted(ds.IsDeleted).
		SetStat(ds.Stat).
		SetEngine(ds.Engine).
		SetDataType(ds.DataType).
		SetDrID(ds.DRID).
		Save(ctx)

	if err != nil {
		logger.CreateReport(&logger.CODE_DB_INSERT, err)
	}

	return inserted
}

func (dao *DatasetDAO) UpdateDatasetDeleted(ctx context.Context, dataset_id int) {
	err := dao.entClient.Dataset.Update().
		Where(dataset.ID(dataset_id)).
		SetIsDeleted(true).
		SetDeletedAt(time.Now()).
		Exec(ctx)

	if err != nil {
		logger.CreateReport(&logger.CODE_DB_DELETE, err)
	}
}

func (dao *DatasetDAO) UpdateDatasetExist(ctx context.Context, dataset_id int) {
	err := dao.entClient.Dataset.Update().
		Where(dataset.ID(dataset_id)).
		SetIsDeleted(false).
		SetUpdatedAt(time.Now()).
		Exec(ctx)

	if err != nil {
		logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}
}

func (dao *DatasetDAO) UpdateValidation(ctx context.Context, ds DatasetDTO) {
	err := dao.entClient.Dataset.Update().
		Where(dataset.ID(ds.ID)).
		SetName(ds.Name).
		SetUpdatedAt(time.Now()).
		SetIsValid(ds.IsValid).
		SetIsTrainable(ds.IsTrainable).
		SetIsTestable(ds.IsTestable).
		SetDataType(ds.DataType).
		SetEngine(ds.Engine).
		SetDescription(ds.Description).
		Exec(ctx)

	if err != nil {
		logger.CreateReport(&logger.CODE_DB_UPDATE, err)
		return
	}
}

func (dao *DatasetDAO) UpdateStat(ctx context.Context, id int, stat string) *logger.Report {
	err := dao.entClient.Dataset.Update().
		Where(dataset.ID(id)).
		SetStat([]string{stat}).
		SetUpdatedAt(time.Now()).
		Exec(ctx)

	if err != nil {
		return logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}

	return nil
}

func (dao *DatasetDAO) DeleteDataset(ctx context.Context, id int) *logger.Report {
	err := dao.entClient.Dataset.DeleteOneID(id).Exec(ctx)

	if err != nil {
		return logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}

	return nil
}

func (dao *DatasetDAO) DeleteDatasetByDRID(ctx context.Context, dr_id int) *logger.Report {
	_, err := dao.entClient.Dataset.Delete().Where(dataset.DrID(dr_id)).Exec(ctx)

	if err != nil {
		return logger.CreateReport(&logger.CODE_DB_DELETE, err)
	}

	return nil
}

func (dao *DatasetDAO) UpdateStatPath(ctx context.Context, id int, stat_path string) *logger.Report {
	err := dao.entClient.Dataset.Update().
		Where(dataset.ID(id)).
		SetStatPath(stat_path).
		SetUpdatedAt(time.Now()).
		Exec(ctx)

	if err != nil {
		return logger.CreateReport(&logger.CODE_DB_UPDATE, err)
	}

	return nil
}

func (dao *DatasetDAO) SelectDataSetByName(ctx context.Context,
	name string) ([]*ent.Dataset, *logger.Report) {
	dss, err := dao.entClient.Dataset.
		Query().
		Where(dataset.Name(name)).
		Order(dataset.ByName(sql.OrderAsc())).
		All(ctx)

	if err != nil {
		return nil, logger.CreateReport(&logger.CODE_DB_SELECT, err)
	}

	return dss, nil
}
