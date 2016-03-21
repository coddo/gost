package api

import (
	"gost/models"
	"io/ioutil"
)

// SingleDataResponse creates a Response from the api, containing a single entity encoded as JSON
func SingleDataResponse(statusCode int, data interface{}) Response {
	jsonData, err := models.SerializeJSON(data)
	if err != nil {
		return InternalServerError(err)
	}

	return Response{
		StatusCode: statusCode,
		Content:    jsonData,
	}
}

// MultipleDataResponse creates a Response from the api, containing an array of entities encoded as JSON
func MultipleDataResponse(statusCode int, data interface{}) Response {
	jsonData, err := models.SerializeJSON(data)
	if err != nil {
		return InternalServerError(err)
	}

	return Response{
		StatusCode: statusCode,
		Content:    jsonData,
	}
}

// StatusResponse creates a Response from the api, containing just a status code
func StatusResponse(statusCode int) Response {
	return Response{StatusCode: statusCode}
}

// ByteResponse creates a Response from the api, containing a status code and a message in the form of a byte array
func ByteResponse(statusCode int, data []byte) Response {
	return Response{
		StatusCode: statusCode,
		Content:    data,
	}
}

// FileResponse creates a Response from the api, containing a file path (download, load or stream)
// and the content type of the file that is returned
func FileResponse(statusCode int, contentType, fullFilePath string) Response {
	if _, err := ioutil.ReadFile(fullFilePath); err != nil {
		return InternalServerError(err)
	}

	return Response{
		StatusCode:  statusCode,
		File:        fullFilePath,
		ContentType: contentType,
	}
}
