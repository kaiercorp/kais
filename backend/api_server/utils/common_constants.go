package utils

/*
	constants
*/

const (
	API_BASE_URL_V1 = "/api/v1"
	WS_BASE_URL_v1  = "/ws/v1"
	SW_VERSION      = "v2.0.1"
)

const (
	DATABASE_DEFAULT_USER = "kaierdba"
	DATABASE_DEFAULT_PWD  = "kaier"
)

const (
	VIS_CLS_TGUI    = "VIS_CLS_TGUI"
	VIS_CLS_SL_TGUI = "VIS_CLS_SL_TGUI"
	VIS_CLS_ML_TGUI = "VIS_CLS_ML_TGUI"
	VIS_OD_TGUI     = "VIS_OD_TGUI"
	VIS_OCR_TGUI    = "VIS_OCR_TGUI"
	VIS_SEG_TGUI    = "VIS_SEG_TGUI"
	VIS_AD_TGUI     = "VIS_AD_TGUI"
	TAB_CLS_TGUI    = "TAB_CLS_TGUI"
	TAB_REG_TGUI    = "TAB_REG_TGUI"
	TS_AD_TGUI      = "TS_AD_TGUI"
	TS_DF_TGUI      = "TS_DF_TGUI"
)

const (
	JOB_TYPE_VISION_CLS_SL = "vcls-sl"
	JOB_TYPE_VISION_CLS_ML = "vcls-ml"
	JOB_TYPE_VISION_SEG    = "vseg"
	JOB_TYPE_VISION_OD     = "vod"
	JOB_TYPE_VISION_OCR    = "vocr"
	JOB_TYPE_VISION_AD     = "vad"
	JOB_TYPE_TABLE_CLS     = "tcls"
	JOB_TYPE_TABLE_REG     = "treg"
	JOB_TYPE_TS_AD         = "tsad"
	JOB_TYPE_TS_DF         = "tsdf"
	JOB_TYPE_INVALID       = "invalid"

	TAPI_JOB_TYPE_VISION_CLS_SL = "vcls_sl"
	TAPI_JOB_TYPE_VISION_CLS_ML = "vcls_ml"
)

const (
	TASK_STEP_IDLE = "idle"

	MODELING_TYPE_INITIAL    = "initial"
	MODELING_TYPE_UPDATE     = "update"
	MODELING_TYPE_EVALUATION = "evaluation"
	MODELING_TYPE_BLIND      = "blind"

	MODELING_STEP_IDLE     = "idle"
	MODELING_STEP_RUN      = "run"
	MODELING_STEP_FINISH   = "finish"   // 모델링 완료
	MODELING_STEP_COMPLETE = "complete" // 모델 평가 완료
	MODELING_STEP_CANCEL   = "cancel"
	MODELING_STEP_FAIL     = "fail"
	MODELING_STEP_REQUEST  = "request" // sended to engine

	GPU_STATE_IDLE       = "idle"
	GPU_STATE_MODEING    = "modeling"
	GPU_STATE_EVALUATION = "inference"
	GPU_STATE_COMPLETE   = "complete"
)

const (
	EXT_CSV  = ".csv"
	EXT_XLS  = ".xls"
	EXT_XLSX = ".xlsx"
	EXT_PNG  = ".png"
	EXT_JPG  = ".jpg"
	EXT_JPEG = ".jpeg"
	EXT_GIF  = ".gif"
	EXT_BMP  = ".bmp"
	EXT_PPM  = ".ppm"
	EXT_PGM  = ".pgm"
	EXT_TIF  = ".tif"
	EXT_TIFF = ".tiff"
	EXT_WEBP = ".webp"
	EXT_NPY  = ".npy"
)

const (
	DATA_TYPE_INVALID       = "invalid"
	DATA_TYPE_TABLE         = "table"
	DATA_TYPE_IMG           = "image"
	DATA_TYPE_TIMESERIES    = "ts"
	DATA_FORMAT_KAIER_TVT   = "kaierTVT"
	DATA_FORMAT_KAIER_TV    = "kaierTV"
	DATA_FORMAT_KAIER_TT    = "kaierTT"
	DATA_FORMAT_KAIER_TRAIN = "kaierT"
	DATA_FORMAT_NONE        = "none"
)

const (
	CONFIG_TYPE_SYSTEM = "SYSTEM"
	CONFIG_TYPE_USER   = "USER"
)

const (
	DIR_TRAIN = "train"
	DIR_VALID = "valid"
	DIR_TEST  = "test"
)
