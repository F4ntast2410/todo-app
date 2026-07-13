package entity

type VerificationEmail string

type LoginResult struct {
	SessionToken string
	Requires2FA  bool
}
