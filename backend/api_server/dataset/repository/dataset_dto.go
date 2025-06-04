package repository

import (
	"time"

	"github.com/montanaflynn/stats"

	"api_server/ent"
)

type DatasetStatistics struct {
	ClassStatics                *ClassStatics                `json:"classStatics,omitempty"`
	MultiLabelClassStatics      *MultiLabelClassStatics      `json:"multiLabelClassStatics,omitempty"`
	ResolutionStatics           *ResolutionStatics           `json:"resolutionStatics,omitempty"`
	Features                    map[string][]string          `json:"features,omitempty"`
	CategoricalFeatureStatics   *CategoricalFeatureStatics   `json:"categoricalFeatureStatics,omitempty"`
	NumericalFeatureStatics     *NumericalFeatureStatics     `json:"numericalFeatureStatics,omitempty"`
	NumericalHeatmap            *NumericalHeatmap            `json:"numericalHeatmap,omitempty"`
	CategoricalHeatmap          *CategoricalHeatmap          `json:"categoricalHeatmap,omitempty"`
	CategoricalNumericalHeatmap *CategoricalNumericalHeatmap `json:"categoricalNumericalHeatmap,omitempty"`
	NoneTypeStat                *NoneTypeStat                `json:"noneTypeStat,omitempty"`
}

type ClassStatics struct {
	Class map[string]map[string]int `json:"class,omitempty"`
	Count map[string]int            `json:"count,omitempty"`
}

type MultiLabelClassStatics struct {
	Class map[string]map[int]int `json:"class,omitempty"`
	Count map[string]int         `json:"count,omitempty"`
}

type ResolutionStatics struct {
	Resolution map[string]map[int]map[int]int `json:"resolution,omitempty"`
	Count      map[string]int                 `json:"count,omitempty"`
}

type CategoricalFeatureStatics struct {
	FeatureStat   []map[string]string            `json:"feature_stat,omitempty"`   // table
	FeatureDetail map[string][]map[string]string `json:"feature_detail,omitempty"` // detail table
}

type NumericalFeatureStatics struct {
	FeatureStat   []map[string]string             `json:"feature_stat,omitempty"`   // table
	FeatureDetail map[string][]map[string]string  `json:"feature_detail,omitempty"` // detail table
	BoxPlot       map[string]map[string][]float64 `json:"box_plot,omitempty"`       // box plot
	PDF           map[string]map[string]DataPDF   `json:"pdf,omitempty"`            // gaussian
}

type NumericalHeatmap struct {
	Feature     []string                                 `json:"feature"`
	Correlation map[string]map[string]map[string]float64 `json:"correlation"`
}

type CategoricalHeatmap struct {
	Feature []string                                 `json:"feature"`
	Heatmap map[string]map[string]map[string]float64 `json:"heatmap"`
}

type CompareNumericalFeaturesStatics struct {
	CompareResult map[string][]stats.Coordinate `json:"compareResult"`
	Regression    map[string]stats.Series       `json:"regression"`
}

type CompareFeaturesStatics struct {
	DatasetId int    `json:"id"`
	Feature1  string `json:"feature1"`
	Feature2  string `json:"feature2"`
}
type CompareCategoricalFeaturesStatics struct {
	CategoricalFeatures       []string                             `json:"categoricalFeatures"`
	CategoricalDetailFeatures map[string][]string                  `json:"categoricalDetailFeatures"`
	CompareResult             map[string]map[string]map[string]int `json:"compareResult"`
}

type CompareCategoricalNumericalFeaturesStatics struct {
	Feature                   []string                      `json:"feature"`
	CategoricalDetailFeatures map[string][]string           `json:"categoricalDetailFeatures"`
	PDF                       map[string]map[string]DataPDF `json:"pdfs"`
}
type CategoricalNumericalHeatmap struct {
	Feature map[string][]string `json:"feature"`
	Heatmap HeatmapDTO          `json:"heatmap"`
}

type HeatmapDTO struct {
	Train any `json:"train"`
	Valid any `json:"valid"`
	Test  any `json:"test"`
}

type DataPDF struct {
	XData []float64 `json:"xData"`
	YData []float64 `json:"yData"`
}

type DataPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type NoneTypeStat struct {
	ImageStat   map[string]*ImageTypeStat          `json:"imageStat,omitempty"`   // image count per directories
	TabularStat map[string]map[string]*TabularStat `json:"tabularStat,omitempty"` // data count per directory, file
}

type ImageTypeStat struct {
	Count             int                `json:"count"`
	ResolutionStatics *ResolutionStatics `json:"resolution"`
}

type TabularStat struct {
	Count    int                 `json:"count"`
	Features map[string][]string `json:"features,omitempty"`
}

type DatasetDTO struct {
	Name        string        `json:"name,omitempty"`
	ID          int           `json:"id,omitempty"`
	ParentID    int           `json:"parent_id"`
	Description string        `json:"description,omitempty"`
	Path        string        `json:"path,omitempty"`
	IsDeleted   bool          `json:"is_deleted,omitempty"`
	IsLeaf      bool          `json:"is_leaf,omitempty"`
	IsValid     bool          `json:"is_valid,omitempty"`
	IsTrainable bool          `json:"is_trainable,omitempty"`
	IsTestable  bool          `json:"is_testable,omitempty"`
	IsUse       bool          `json:"is_use,omitempty"`
	DataType    string        `json:"data_type,omitempty"`
	Engine      []string      `json:"engine,omitempty"`
	Stat        []string      `json:"stat,omitempty"`
	StatPath    string        `json:"stat_path,omitempty"`
	CreatedAt   time.Time     `json:"created_at,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
	DeletedAt   time.Time     `json:"deleted_at,omitempty"`
	Childs      []*DatasetDTO `json:"dirs,omitempty"`
	DRID        int           `json:"dataset_root_datasets"`
}

type FeatureType string

func ConvertDatasetEntToDTO(entity *ent.Dataset) *DatasetDTO {
	return &DatasetDTO{
		ID:          entity.ID,
		Name:        entity.Name,
		ParentID:    entity.ParentID,
		Description: entity.Description,
		Path:        entity.Path,
		IsValid:     entity.IsValid,
		IsTrainable: entity.IsTrainable,
		IsTestable:  entity.IsTestable,
		IsLeaf:      entity.IsLeaf,
		IsDeleted:   entity.IsDeleted,
		IsUse:       entity.IsUse,
		Stat:        entity.Stat,
		StatPath:    entity.StatPath,
		Engine:      entity.Engine,
		DataType:    entity.DataType,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DeletedAt:   entity.DeletedAt,
	}
}

func ConvertDatasetEntsToDTOs(ents []*ent.Dataset) []*DatasetDTO {
	if ents == nil {
		return nil
	}

	dDTOs := []*DatasetDTO{}

	for _, d := range ents {
		dDTOs = append(dDTOs, ConvertDatasetEntToDTO(d))
	}

	return dDTOs
}

type GetDatasetsDTO struct {
	Datasets  []*DatasetDTO `json:"datasets"`
	TotalPage int           `json:"totalPage"`
	HasMore   bool          `json:"hasMore"`
	NextPage  int           `json:"nextPage"`
}

func ConvertDatasetEntsToGetDatasetsDTOs(ents []*ent.Dataset, pageCount int, hasMore bool, nextPage int) *GetDatasetsDTO {
	return &GetDatasetsDTO{
		Datasets:  ConvertDatasetEntsToDTOs(ents),
		TotalPage: pageCount,
		HasMore:   hasMore,
		NextPage:  nextPage,
	}
}
