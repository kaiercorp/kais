package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type Report struct {
	Code    string
	Message string
	Error   error
}

func CreateReport(s *State, e error) *Report {
	if s != nil {
		ReportError(s.Code, s.Message, e)
		return &Report{Code: s.Code, Message: s.Message, Error: e}
	} else {
		ReportError(CODE_FAILE.Code, CODE_FAILE.Message, e)
		return &Report{Code: CODE_FAILE.Code, Message: CODE_FAILE.Message, Error: e}
	}
}

func (r *Report) SetCode(c string) *Report {
	r.Code = c
	return r
}

func (r *Report) SetMessage(m string) *Report {
	r.Message = m
	return r
}

func (r *Report) SetError(e error) *Report {
	r.Error = e
	return r
}

func ApiRequest(c *gin.Context) {
	message := fmt.Sprintf(`[%s] REQ (%s:%s)`, requestid.Get(c), c.Request.Method, c.Request.URL.Path)
	Info(message)
}

func ApiResponse(c *gin.Context, r *Report, data any) {
	message := fmt.Sprintf(`[%s] RES (%s:%s)`, requestid.Get(c), c.Request.Method, c.Request.URL.Path)
	if r != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    r.Code,
			"message": r.Message,
		})
		Error(message)
	} else if data != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"code":    CODE_SUCCESS.Code,
			"message": CODE_SUCCESS.Message,
			"data":    data,
		})
		Info(message)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"code":    CODE_SUCCESS.Code,
			"message": CODE_SUCCESS.Message,
		})
		Info(message)
	}
}

func TApiResponse(c *gin.Context, r *Report, data any) {
	message := fmt.Sprintf(`[%s] RES (%s:%s)`, requestid.Get(c), c.Request.Method, c.Request.URL.Path)
	if r != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    r.Code,
			"message": r.Message,
		})
		Error(message)
	} else if data != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"code":    CODE_TAPI_SUCCESS.Code,
			"message": CODE_TAPI_SUCCESS.Message,
			"data":    data,
		})
		Info(message)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"code":    CODE_TAPI_SUCCESS.Code,
			"message": CODE_TAPI_SUCCESS.Message,
		})
		Info(message)
	}
}

func ApiResponseWithZipFile(c *gin.Context, filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		r := CreateReport(&CODE_FILE_OPEN, err)
		ApiResponse(c, r, nil)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		r := CreateReport(&CODE_FILE_READ, err)
		ApiResponse(c, r, nil)
		return
	}

	// Set Content-Disposition in extraHeader
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf("attachment; filename=%s", fileInfo.Name()),
	}

	Info(fmt.Sprintf(`[%s] RES (%s:%s)`, requestid.Get(c), c.Request.Method, c.Request.URL.Path))
	c.DataFromReader(
		http.StatusOK,
		fileInfo.Size(),
		"application/zip",
		file,
		extraHeaders, // Pass extraHeaders here
	)
}

// ApiResponseWithJsonFile 함수는 JSON 파일을 읽어 클라이언트에게 반환하는 API 응답을 처리합니다.
// 클라이언트가 gzip 압축을 지원하면 gzip으로 압축하여 전달하며, 그렇지 않으면 원본 파일을 그대로 전송합니다.
//
// 파라미터:
//   - c: *gin.Context - Gin의 컨텍스트 객체
//   - filepath: string - 클라이언트에게 전달할 JSON 파일 경로
func ApiResponseWithJsonFile(c *gin.Context, filepath string) {
	// 파일 열기
	file, err := os.Open(filepath)
	if err != nil {
		r := CreateReport(&CODE_FILE_OPEN, err)
		ApiResponse(c, r, nil)
		return
	}
	defer file.Close()

	// 파일 정보 가져오기
	fileInfo, err := file.Stat()
	if err != nil {
		r := CreateReport(&CODE_FILE_READ, err)
		ApiResponse(c, r, nil)
		return
	}

	// 클라이언트의 gzip 압축 지원 여부 확인
	acceptEncoding := c.GetHeader("Accept-Encoding")
	supportsGzip := containsGzip(acceptEncoding)

	// 요청 정보 로그 출력
	Info(fmt.Sprintf(`[%s] RES (%s:%s)`, requestid.Get(c), c.Request.Method, c.Request.URL.Path))

	// 클라이언트가 gzip을 지원하면 압축하여 전송
	if supportsGzip {
		// gzip 압축을 위한 io.Pipe 생성
		reader, writer := io.Pipe()
		gzipWriter := gzip.NewWriter(writer)

		// 별도 고루틴에서 gzip 압축 처리
		go func() {
			defer writer.Close()
			defer gzipWriter.Close()
			_, err := io.Copy(gzipWriter, file)
			if err != nil {
				r := CreateReport(&CODE_FILE_READ, err)
				ApiResponse(c, r, nil)
				return
			}
		}()

		// gzip 압축 헤더 설정
		extraHeaders := map[string]string{
			"Content-Encoding": "gzip",
		}

		// 압축된 데이터를 스트리밍하여 클라이언트에게 전송
		c.DataFromReader(http.StatusOK, -1, "application/json", reader, extraHeaders)
	} else {
		// 일반 JSON 데이터를 클라이언트에게 스트리밍 전송
		c.DataFromReader(http.StatusOK, fileInfo.Size(), "application/json", file, nil)
	}
}

// containsGzip 함수는 Accept-Encoding 헤더 값에서 gzip이 포함되어 있는지 확인합니다.
func containsGzip(header string) bool {
	return strings.Contains(header, "gzip")
}
