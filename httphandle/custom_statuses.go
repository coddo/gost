package httphandle

const (
	StatusTooManyRequests = 429
)

var statusText = map[int]string{
	StatusTooManyRequests: "Too many requests",
}

func StatusText(statusCode int) string {
	return statusText[statusCode]
}
