package utils

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"api_server/logger"

	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
)

func ReadTabularFile(filePath string) ([][]string, error) {
	switch filepath.Ext(filePath) {
	case EXT_CSV:
		return ReadCsvFile(filePath)
	case EXT_XLS, EXT_XLSX:
		return ReadXlsFile(filePath)
	default:
		return nil, errors.New(filePath + " is not tabular type.")
	}
}

func ReadCsvFile(csvFilePath string) ([][]string, error) {
	if _, err := os.Stat(csvFilePath); os.IsNotExist(err) {
		logger.Error(csvFilePath, " is not exist.")
		return nil, err
	}

	ext := filepath.Ext(csvFilePath)
	if strings.ToLower(ext) != EXT_CSV {
		logger.Error("Invalid file extension.")
		return nil, errors.New("Invalid file extension: " + ext)
	}

	file, err := os.Open(csvFilePath)
	if err != nil {
		logger.Error("Failed to open ", csvFilePath, ": ", err)
		return nil, err
	}
	defer file.Close()

	rdr := csv.NewReader(bufio.NewReader(file))
	rows, err := rdr.ReadAll()
	if err != nil {
		logger.Error("Failed to read ", csvFilePath, ": ", err)
		return nil, err
	}

	rows[0][0] = strings.TrimPrefix(rows[0][0], string('\uFEFF'))

	return rows, nil
}

func ReadXlsFile(xlsFilePath string) ([][]string, error) {
	if _, err := os.Stat(xlsFilePath); os.IsNotExist(err) {
		logger.Error(xlsFilePath, " is not exist.")
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(xlsFilePath))
	if ext != EXT_XLS && ext != EXT_XLSX {
		logger.Error("Invalid file extension.")
		return nil, errors.New("Invalid file extension: " + ext)
	}

	// .xls 파일 처리
	if ext == EXT_XLS {
		// github.com/extrame/xls 라이브러리 사용
		xlsWorkbook, err := xls.Open(xlsFilePath, "utf-8")
		if err != nil {
			logger.Error("Failed to open .xls file: ", err)
			return nil, err
		}

		sheet := xlsWorkbook.GetSheet(0)
		if sheet == nil {
			logger.Error("No sheet found in .xls file")
			return nil, errors.New("no sheet found in .xls file")
		}

		var rows [][]string
		for i := 0; i <= int(sheet.MaxRow); i++ {
			row := sheet.Row(i)
			if row == nil {
				continue
			}

			var rowData []string
			for j := 0; j <= int(row.LastCol()); j++ {
				rowData = append(rowData, row.Col(j))
			}
			rows = append(rows, rowData)
		}

		// BOM 제거 (첫 번째 셀에 있는 경우)
		if len(rows) > 0 && len(rows[0]) > 0 {
			rows[0][0] = strings.TrimPrefix(rows[0][0], string('\uFEFF'))
		}

		return rows, nil
	} else {
		// .xlsx 파일 처리 (기존 코드)
		xlsFile, err := excelize.OpenFile(xlsFilePath)
		if err != nil {
			logger.Error("Failed to open ", xlsFilePath, ": ", err)
			return nil, err
		}
		defer xlsFile.Close()

		sheets := xlsFile.GetSheetList()
		rows, err := xlsFile.GetRows(sheets[0])
		if err != nil {
			logger.Error("Failed to read ", xlsFilePath, ": ", err)
			return nil, err
		}

		if len(rows) > 0 && len(rows[0]) > 0 {
			rows[0][0] = strings.TrimPrefix(rows[0][0], string('\uFEFF'))
		}

		return rows, nil
	}
}

func ReadJsonFile(jsonFilePath string) string {
	file, err := os.Open(jsonFilePath)
	if err != nil {
		logger.Error("File Open Error - ", err.Error())
		return ""
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var data_map map[string]interface{}
	err = decoder.Decode(&data_map)
	if err != nil {
		logger.Error("Failed to Decode Json Reuslt - ", err.Error())
		return ""
	}

	dataBytes, err := json.Marshal(data_map)
	if err != nil {
		logger.Error("Json marshal Error : ", err.Error())
		return ""
	}

	return string(dataBytes)
}

func RemoveDirectory(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		logger.Error("Remove directory error: ", err.Error())
	}
}

func GetFileList(path string) []string {
	files := []string{}

	parent, err := os.Open(path)
	if err != nil {
		logger.Error("Open dataset directory Error : ", err.Error())
		return files
	}
	defer parent.Close()

	childList, err := parent.Readdir(-1)
	if err != nil {
		logger.Error("Open dataset child directories Error : ", err.Error())
		return files
	}

	for _, child := range childList {
		if !child.IsDir() {
			childPath := filepath.Join(path, child.Name())
			files = append(files, childPath)
		}
	}

	return files
}

func MoveFile(src string, dest string) error {
	destFile := filepath.Join(dest, filepath.Base(src))

	err := os.Rename(src, destFile)
	if err != nil {
		return err
	}

	return nil
}

func IsImageFile(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case EXT_PNG:
		return true
	case EXT_JPG:
		return true
	case EXT_JPEG:
		return true
	case EXT_GIF:
		return true
	case EXT_BMP:
		return true
	case EXT_PPM:
		return true
	case EXT_PGM:
		return true
	case EXT_TIF:
		return true
	case EXT_TIFF:
		return true
	case EXT_WEBP:
		return true
	default:
		return false
	}
}

func IsTabularFile(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case EXT_CSV:
		return true
	case EXT_XLS:
		return true
	case EXT_XLSX:
		return true
	default:
		return false
	}
}

func FindImageFile(root_path string) string {
	files, _ := os.ReadDir(root_path)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if IsImageFile(file.Name()) {
			return root_path + "/" + file.Name()
		}
	}

	for _, file := range files {
		if file.IsDir() {
			result := FindImageFile(root_path + "/" + file.Name())
			if result != "" {
				return result
			}
		}
	}

	return ""
}

// ReadDirs reads the directories in the path
func ReadDirs(path string) ([]os.DirEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dirs := []os.DirEntry{}
	allFiles, err := f.ReadDir(-1)
	for _, f := range allFiles {
		if f.IsDir() {
			dirs = append(dirs, f)
		}
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	return dirs, err
}

func ReadFiles(path string, accept []string, filter []string) ([]os.DirEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	files := []os.DirEntry{}
	allFiles, err := f.ReadDir(-1)
	for _, f := range allFiles {
		if f.IsDir() {
			continue
		}
		if f.Name() == "Thumbs.db" || f.Name() == ".DS_Store" {
			continue
		}
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		if len(accept) > 0 {
			for _, filt := range accept {
				if filt == "" {
					continue
				}
				if strings.HasSuffix(strings.ToLower(f.Name()), filt) {
					files = append(files, f)
				}
			}
		} else if len(filter) > 0 {
			hasFilter := false
			for _, filt := range filter {
				if strings.HasSuffix(strings.ToLower(f.Name()), filt) {
					hasFilter = true
					break
				}
			}
			if !hasFilter {
				files = append(files, f)
			}
		} else {
			files = append(files, f)
		}
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	return files, err
}

func RefinePathSeparator(path string) string {
	refinePath := strings.ReplaceAll(path, "/", string(os.PathSeparator))
	refinePath = strings.ReplaceAll(refinePath, "\\", string(os.PathSeparator))
	return refinePath
}
