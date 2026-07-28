package entity

type VerificationEmail string

type AuthResult struct {
	SessionToken string
	Requires2FA  bool
}

type EmailResult struct {
	Requires2FA bool
}
