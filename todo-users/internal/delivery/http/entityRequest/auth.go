package dto

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` // сырой пароль, только для входа
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"` // сырой пароль, только для входа
}
type Responce2FA struct {
	Status bool `json:"status"`
}
type Request2FA struct {
	InputCode string `json:"input_code"`
	Email     string `json:"email"`
}
