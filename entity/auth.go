package entity

type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type MetaToken struct {
	AccountID     string `json:"aid"`
	ExpiredAt     int64  `json:"exp"`
	IssuedAt      int64  `json:"iat"`
	Authorization bool   `json:"aut"`
}

type AccessToken struct {
	Claims MetaToken
}

// +++++++++++++ RESPONSE +++++++++++++++

type LoginRes struct {
	Authorization string `json:"authorization"`
}
