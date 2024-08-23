package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "Ok"
	StatusError = "Error"
)

func Ok() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

/*func ValidationError(errs validator.ValidationErrors) Response {
	var errsMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errsMsgs = append(errsMsgs, fmt.Sprintf("Field %s is required", err.Field()))
		case "alpha":
			errsMsgs = append(errsMsgs, fmt.Sprintf("Field %s must contain only letters", err.Field()))
		case "alphanum":
			errsMsgs = append(errsMsgs, fmt.Sprintf("Field %s must contain only letters and numbers", err.Field()))
		default:
			errsMsgs = append(errsMsgs, fmt.Sprintf("Field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errsMsgs, ", "),
	}
}*/
