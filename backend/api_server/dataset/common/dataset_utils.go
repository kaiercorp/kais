package common

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"

	repo "api_server/dataset/repository"
	"api_server/logger"
)

// Child directory들을 가져옴
func GetChildDirDataset(path string) ([]repo.DatasetDTO, error) {
	childDirs := []repo.DatasetDTO{}
	parent, err := os.Open(path)

	if err != nil {
		logger.Error("Open dataset directory Error : ", err.Error())
		return nil, err
	}

	childList, err := parent.Readdir(-1)
	if err != nil {
		logger.Error("Open dataset child directories Error : ", err.Error())
		return nil, err
	}

	for _, child := range childList {
		if child.IsDir() {
			childPath := filepath.Join(path, child.Name())
			childDirs = append(childDirs, repo.DatasetDTO{Name: child.Name(), Path: childPath})
		}
	}

	return childDirs, nil
}

// Child directory들을 가져옴
func GetChildDirNames(path string) ([]string, error) {
	childDirs := []string{}
	parent, err := os.Open(path)

	if err != nil {
		logger.Error("Open directory Error : ", err.Error())
		return nil, err
	}

	childList, err := parent.Readdir(-1)
	if err != nil {
		logger.Error("Open child directories Error : ", err.Error())
		return nil, err
	}

	for _, child := range childList {
		if child.IsDir() {
			childDirs = append(childDirs, child.Name())
		}
	}

	return childDirs, nil
}

func GetLabelNames(path string) ([]string, error) {
	file, err := os.Open(filepath.Join(path, "label.txt"))
	if err != nil {
		logger.Error("Labeling File is Not Open - ", path)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	labelMap := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()

		labelInfo := strings.Split(line, " ")[1:]

		for _, label := range labelInfo {
			if label == "" {
				continue
			}
			labelMap[label] = true
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Labeling File Scanning Error - ", path)
		return nil, err
	}

	var labels []string
	for key := range labelMap {
		labels = append(labels, key)
	}

	sort.Strings(labels)

	return labels, nil
}

func DeleteDuplicate(arr []string) []string {
	result := []string{}
	temp := make(map[string]struct{})

	for _, val := range arr {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
			result = append(result, val)
		}
	}
	return result
}
