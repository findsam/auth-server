package error

var errors = map[int]string{
	401: "Unauthorized request",
	400: "Client error occured",
	500: "Server error occured",
	204: "No content found",
}

func Message(code int) string {
	if msg, exists := errors[code]; exists {
		return msg
	}
	return "An error has occured, try again later."
}
