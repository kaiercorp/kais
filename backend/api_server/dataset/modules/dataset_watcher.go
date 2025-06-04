package modules

import (
	"context"
	"path/filepath"
	"time"

	repo "api_server/dataset/repository"
	"api_server/ent"
	"api_server/logger"
	"api_server/utils"
)

type DatasetWatcherInterface interface {
	WatchDataset()
	DetectDatasetModification()
}

type DatasetWatcher struct {
	ctx              context.Context
	datasetValidator DatasetValidatorInterface
	datasetAnalyzer  DatasetAnalyzerInterface
	datasetDAO       repo.DatasetDAOInterface
	datasetRootDAO   repo.DatasetRootDAOInterface
	dbDatasets       []*ent.Dataset
	diskDatasets     []*repo.DatasetDTO
}

var datasetWatcher *DatasetWatcher

// NewDatasetWatcher returns singleton datasetWatcher object
func NewDatasetWatcher(
	datasetValidator DatasetValidatorInterface,
	datasetAnalyzer DatasetAnalyzerInterface,
	datasetDAO repo.DatasetDAOInterface,
	datasetRootDAO repo.DatasetRootDAOInterface) *DatasetWatcher {

	if datasetWatcher == nil {
		datasetWatcher = &DatasetWatcher{
			ctx:              context.Background(),
			datasetValidator: datasetValidator,
			datasetAnalyzer:  datasetAnalyzer,
			datasetDAO:       datasetDAO,
			datasetRootDAO:   datasetRootDAO,
		}
	}

	return datasetWatcher
}

// WatchDataset detects the dataset every 20 seconds
// and update dataset table in the dbms
func (w *DatasetWatcher) WatchDataset() {
	duration, _ := time.ParseDuration("20s")
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	w.DetectDatasetModification()
	for range ticker.C {
		w.DetectDatasetModification()
	}
}

func (w *DatasetWatcher) DetectDatasetModification() {
	if datasetRoots, r := w.datasetRootDAO.SelectActive(w.ctx); r == nil {
		for _, datasetRoot := range datasetRoots {
			w.dbDatasets = w.datasetDAO.SelectDatasetsByDRID(w.ctx, datasetRoot.ID)
			w.diskDatasets = []*repo.DatasetDTO{}
			w.findDirs(datasetRoot.Path, nil)
			w.updateDatasets()
			w.addToDBMS(datasetRoot.ID)

			w.datasetValidator.Validate(datasetRoot.ID)
			w.datasetAnalyzer.Analyze(datasetRoot.ID)
		}
	}
}

// read all directories from dataset_root
func (w *DatasetWatcher) findDirs(path string, parent *repo.DatasetDTO) {
	if dirs, err := utils.ReadDirs(path); err != nil {
		logger.Error("Read dirs from disk : ", err.Error())
	} else {
		for _, d := range dirs {
			path1 := filepath.Join(path, d.Name())
			dataset := &repo.DatasetDTO{
				Name: d.Name(),
				Path: path1,
			}
			if parent != nil {
				parent.Childs = append(parent.Childs, dataset)
			} else {
				w.diskDatasets = append(w.diskDatasets, dataset)
			}

			w.findDirs(path1, dataset)
		}
	}
}

// update dbms isDeleted column
func (w *DatasetWatcher) updateDatasets() {
	for _, dsDB := range w.dbDatasets {
		isExist := false
		for _, dsDisk := range w.diskDatasets {
			if dsDB.Path == dsDisk.Path {
				dsDisk.ID = dsDB.ID
				w.datasetDAO.UpdateDatasetExist(w.ctx, dsDB.ID)
				isExist = true
				break
			}
		}
		if !isExist {
			w.datasetDAO.UpdateDatasetDeleted(w.ctx, dsDB.ID)
		}
	}
}

// add new dataset to dbms
func (w *DatasetWatcher) addToDBMS(dr_id int) {
	for _, ds := range w.diskDatasets {
		ds.DRID = dr_id
		w.addDataset(ds)
	}
}

func (w *DatasetWatcher) addDataset(dataset *repo.DatasetDTO) {
	if len(dataset.Childs) < 1 {
		dataset.IsLeaf = true
	}

	dsID := 0
	if dataset.ID < 1 {
		if exist, _ := w.datasetDAO.SelectByPath(w.ctx, dataset.Path); exist == nil || exist.ID < 1 {
			inserted := w.datasetDAO.InsertOne(w.ctx, *dataset)
			dsID = inserted.ID
		} else {
			dsID = exist.ID
		}
	}

	for _, child := range dataset.Childs {
		child.ParentID = dsID
		child.DRID = dataset.DRID
		w.addDataset(child)
	}
}
