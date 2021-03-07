package web

type loginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registrationModel struct {
	loginCredentials
	ConfirmPassword string `json: "confirmPassword"`
}

type loginResponse struct {
	Jwt string `json:"jwt"`
	Username string `json:"username"`
	Role string `json:"role"`
}

type errorResponse struct {
	Message string	`json:"errorMessage"`
	StatusCode int `json:"statusCode"`
}