package repository

type ImageClassModel struct {
	Classes     []string `json:"classes"`
	IsTiff      bool     `json:"is_tiff"`
	FrameNumber int      `json:"frame_number"`
}
