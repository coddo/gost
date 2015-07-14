package api

import (
	"encoding/json"
	"go-server-template/models"
	"io/ioutil"
)

func SingleDataResponse(statusCode int, data models.Serializable) ApiResponse {
	jsonData, err := data.SerializeJson()
	if err != nil {
		return InternalServerError(err)
	}

	return ApiResponse{
		StatusCode: statusCode,
		Message:    jsonData,
	}
}

func MultipleDataResponse(statusCode int, data map[int]models.Serializable) ApiResponse {
	serializableData := make([]models.Serializable, len(data))

	for i := 0; i < len(data); i++ {
		serializableData[i] = data[i]
	}

	jsonData, err := json.MarshalIndent(serializableData, models.JsonPrefix, models.JsonIndent)
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
