package api

import (
	"gost/models"
	"io/ioutil"
)

func SingleDataResponse(statusCode int, data models.Modeler) ApiResponse {
	jsonData, err := models.SerializeJson(data)
	if err != nil {
		return InternalServerError(err)
	}

	return ApiResponse{
		StatusCode: statusCode,
		Message:    jsonData,
	}
}

func MultipleDataResponse(statusCode int, data []models.Modeler) ApiResponse {
	jsonData, err := models.SerializeJson(data)
	if err != nil {
		return InternalServerError(err)
	}

	return ApiResponse{
		StatusCode: statusCode,
		Message:    jsonData,
	}
}

func StatusResponse(statusCode int) ApiResponse {
	return ApiResponse{StatusCode: statusCode}
}

func ByteResponse(statusCode int, data []byte) ApiResponse {
	return ApiResponse{
		StatusCode: statusCode,
		Message:    data,
	}
}

func FileResponse(statusCode int, contentType, fullFilePath string) ApiResponse {
	if _, err := ioutil.ReadFile(fullFilePath); err != nil {
		return InternalServerError(err)
	}

	return ApiResponse{
		StatusCode:  statusCode,
		File:        fullFilePath,
		ContentType: contentType,
	}
}
