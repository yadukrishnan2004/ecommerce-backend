package response


type ApiResponse struct{
	StatusCode  int    			`json:"status"` 
	Message		string 			`json:"message"`
	Data		interface{}		`json:"data,omitemty"` 
	Error		interface{}		`json:"error,omitemty"`
}