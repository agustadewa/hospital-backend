package entity

type TAccountRole string

const (
	PATIENT       TAccountRole = "PATIENT"
	ADMINISTRATOR TAccountRole = "ADMINISTRATOR"
)

type TAccount struct {
	AccountID string       `bson:"account_id" json:"account_id"`
	Role      TAccountRole `bson:"role" json:"role"`
	FirstName string       `bson:"first_name" json:"first_name"`
	LastName  string       `bson:"last_name" json:"last_name"`
	Age       uint8        `bson:"age" json:"age"`
	Email     string       `bson:"email" json:"email"`
	Username  string       `bson:"username" json:"username"`
	Password  string       `bson:"password" json:"password"`
}

// type TUpdateAccount struct {
// 	FirstName *string       `bson:"first_name,omitempty" json:"first_name,omitempty"`
// 	LastName  *string       `bson:"last_name,omitempty" json:"last_name,omitempty"`
// 	Age       *uint8        `bson:"age,omitempty" json:"age,omitempty"`
// 	Password  *string       `bson:"password,omitempty" json:"password,omitempty"`
// }

type TCreateAccountReq struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Age       uint8  `json:"age" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// ++++++++++++ RESPONSE +++++++++++++

type TCreateAccountRes struct {
	AccountID string `json:"account_id"`
}
