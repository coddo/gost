package api

import (
	"gost/util/jsonutil"
	"io/ioutil"
)

// JSONResponse creates a Response from the api, containing a single entity encoded as JSON
func JSONResponse(statusCode int, data interface{}) Response {
	jsonData, err := jsonutil.SerializeJSON(data)
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

// PlainTextResponse creates a Response from the api, containing a status code and a text message
func PlainTextResponse(statusCode int, text string) Response {
	return DataResponse(statusCode, []byte(text), ContentTextPlain)
}

// TextResponse creates a Response from the api, containing a status code, a text message and a custom content-type
func TextResponse(statusCode int, text, contentType string) Response {
	return DataResponse(statusCode, []byte(text), contentType)
}

// DataResponse creates a Response from the api, containing a status code,
// a custom content-type and a message in the form of a byte array
func DataResponse(statusCode int, data []byte, contetType string) Response {
	return Response{
		StatusCode:  statusCode,
		Content:     data,
		ContentType: contetType,
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
