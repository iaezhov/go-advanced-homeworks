package auth

type LoginPayload struct {
	Phone string `json:"phone" validate:"required"`
}
type LoginResponse struct {
	SessionId string `json:"sessionId"`
}

type VerifyPayload struct {
	SessionId string `json:"sessionId" validate:"required"`
	Code      int    `json:"code" validate:"required"`
}
type VerifyResponse struct {
	Token string `json:"token"`
}
