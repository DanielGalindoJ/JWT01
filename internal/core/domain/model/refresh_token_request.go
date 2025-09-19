package model

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
