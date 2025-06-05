package dto

type MenuDTO struct {
	Key       string    `json:"key"`
	Label     string    `json:"label"`
	Icon      string    `json:"icon"`
	Url       string    `json:"url"`
	IsUse     bool      `json:"isUse"`
	IsTitle   bool      `json:"isTitle"`
	ParentKey string    `json:"parentKey"`
	MenuOrder int		`json:"menuOrder"`
	Group     int       `json:"group"`
	Children  []MenuDTO `json:"children"`
}
