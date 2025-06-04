package logger

type State struct {
	Code    string
	Message string
}

var (
	CODE_SUCCESS = State{Code: "SUCCESS", Message: "Success"}
	CODE_FAILE   = State{Code: "9999", Message: "Unknown Error"}

	CODE_DB_INSERT    = State{Code: "DB0011", Message: "Insert query error"}
	CODE_DB_SELECT    = State{Code: "DB0021", Message: "Select query error"}
	CODE_DB_UPDATE    = State{Code: "DB0031", Message: "Update query error"}
	CODE_DB_DELETE    = State{Code: "DB0041", Message: "Delete query error"}
	CODE_DB_DUPLICATE = State{Code: "DB0012", Message: "Duplicate entry error"}

	CODE_REQUEST          = State{Code: "RQ0001", Message: "Invalid request"}
	CODE_JSON_MARSHAL     = State{Code: "RQ0011", Message: "Json Fail"}
	CODE_JSON_UNMARSHAL   = State{Code: "RQ0012", Message: "Json Fail"}
	CODE_API_PARAM_ENGINE = State{Code: "RQ1001", Message: "Invalid engine type"}

	CODE_LOGIN_PARAMS  = State{Code: "AU0001", Message: "Invalid login info"}
	CODE_LOGIN_FAILED  = State{Code: "AU0002", Message: "Failed to login"}
	CODE_LOGOUT_FAILED = State{Code: "AU0003", Message: "Failed to logout"}
	CODE_TOKEN_EXPIRED = State{Code: "AU0004", Message: "Expired login"}

	CODE_REMOTE_SELECT          = State{Code: "RM0001", Message: "Failed to call engine"}
	CODE_REMOTE_CREATE_REQ      = State{Code: "RM0002", Message: "Can't create request"}
	CODE_REMOTE_REQUEST         = State{Code: "RM0003", Message: "Can't connect to server"}
	CODE_REMOTE_RESPONSE        = State{Code: "RM0004", Message: "Can't parse response"}
	CODE_REMOTE_NOT_FOUND_MODEL = State{Code: "RM0005", Message: "Not found model"}

	CODE_FILE_NOT_EXIST   = State{Code: "FL0001", Message: "File Not Exist"}
	CODE_FILE_OPEN        = State{Code: "FL0002", Message: "File Open Error"}
	CODE_FILE_READ        = State{Code: "FL0003", Message: "File Read Error"}
	CODE_DIR_NOT_EXIST    = State{Code: "FL0005", Message: "Directory Not Exist"}
	CODE_TESTID_NOT_EXIST = State{Code: "FL0005", Message: "Test ID Not Exist"}
	CODE_INVALID_METRIC   = State{Code: "FL0006", Message: "Invalid Metric"}

	CODE_DATA_TABLE_TYPE  = State{Code: "5101", Message: "It is not Table data"}
	CODE_DATA_IMAGE_TYPE  = State{Code: "5201", Message: "It is not Image data"}
	CODE_DATA_IMAGE_CLASS = State{Code: "5202", Message: "Not exist class folders"}

	CODE_EXECUTE    = State{Code: "EX001", Message: "Failed to execute code"}
	CODE_CHANGE_DIR = State{Code: "CH001", Message: "Fsiled to change directory"}

	CODE_ADD_ZIP   = State{Code: "ZF001", Message: "Failed to add zip file"}
	CODE_ENTRY_ZIP = State{Code: "ZF002", Message: "Failed to create zip entry"}
	CODE_COPY_ZIP  = State{Code: "ZF003", Message: "Failed to copy to zip"}

	CODE_MODELING_IN_PROGRESS      = State{Code: "MR001", Message: "There's something in progress"}
	CODE_MODELING_DEVICE_NOT_EXIST = State{Code: "MR002", Message: "There's nothing usable devices"}

	CODE_LICENSE_FAIL         = State{Code: "LI0001", Message: "Invalid license"}
	CODE_LICENSE_EXPIRED      = State{Code: "LI0002", Message: "Expired"}
	CODE_LICENSE_UNAUTHORIZED = State{Code: "LI0003", Message: "Not allowed"}
	CODE_LICENSE_UNKNOWN      = State{Code: "LI9999", Message: "lincense error"}

	ERROR_CODE_SSH_ERROR = State{Code: "7001", Message: "Failed to find gpu"}
	CODE_TAPI_SUCCESS    = State{Code: "0000", Message: "Success"}
)
