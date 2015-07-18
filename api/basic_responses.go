package api

import (
	"go-server-template/models"
	"io/ioutil"
)

func SingleDataResponse(statusCode int, data models.Expander) ApiResponse {
	jsonData, err := models.SerializeJson(data)
	if err != nil {
		return InternalServerError(err)
	}

	return ApiResponse{
		StatusCode: statusCode,
		Message:    jsonData,
	}
}

func MultipleDataResponse(statusCode int, data map[int]models.Expander) ApiResponse {
	serializableData := make([]models.Expander, len(data))

	for i := 0; i < len(data); i++ {
		serializableData[i] = data[i]
	}

	jsonData, err := models.SerializeJson(serializableData)
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
