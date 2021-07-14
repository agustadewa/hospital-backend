package entity

type TAppointment struct {
	AppointmentID     string   `bson:"appointment_id" json:"appointment_id"`
	DoctorID          string   `bson:"doctor_id" json:"doctor_id"`
	PatientAccountIDs []string `bson:"patient_account_id" json:"patient_account_i_ds"`
	Description       string   `bson:"description" json:"description"`
	MaxAppointment    uint8    `bson:"max_appointment" json:"max_appointment"`
}

type TUpdateAppointment struct {
	DoctorID          *string   `bson:"doctor_id,omitempty" json:"doctor_id,omitempty"`
	PatientAccountIDs *[]string `bson:"patient_account_id,omitempty" json:"patient_account_i_ds,omitempty"`
	Description       *string   `bson:"description,omitempty" json:"description,omitempty"`
	MaxAppointment    *uint8    `bson:"max_appointment,omitempty" json:"max_appointment,omitempty"`
}
