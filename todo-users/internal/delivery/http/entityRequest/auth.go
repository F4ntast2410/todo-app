package dto

type RegisterRequest struct {
	Email string `json:"email"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Responce2FA struct {
	Status bool `json:"status"`
}
type Request2FA struct {
	InputCode string `json:"input_code"`
	Email     string `json:"email"`
	NewEmail  string `json:"new_email,omitempty"`
}
type Request2FARegister struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	InputCode string `json:"input_code"`
}

type RequestSetEmail struct {
	Email string `json:"email"`
}
