package entity

type TDoctor struct {
	DoctorID  string `bson:"doctor_id" json:"doctor_id"`
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
}

type TCreateDoctorReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type TGetManyDoctorReq struct {
	Keyword string `json:"keyword"`
}

// ++++++++++++ RESPONSE ++++++++++++

type TCreateDoctorRes struct {
	DoctorID string `json:"docker_id"`
}
