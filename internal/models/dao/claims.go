package dao

type UserClaims struct {
	UserID      string   `json:"user_id"`
	Permissions []string `json:"permissions"`
	Roles       []string `json:"roles"`
}
