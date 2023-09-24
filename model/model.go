package model

// AddNewUser provides data to create a User.
//
// swagger:model
type User struct {
	Id int `json:"id" db:"id" example:"906"`
	//min_length = 3
	//min_length = 20
	//exclude = \"\\\/
	Username string `json:"username" db:"username" example:"Andrew"`
	//min_length = 6
	//min_length = 20
	//exclude = \"\\\/
	Password string `json:"password" db:"password" example:"$2a$14$kv/sGmTWIlNYocbZqd88GuRsrOtKrs9bBFMM7N7HRNZ.qPxF.b.GG"`
	//YYYY-MM-DD HH:MM:SS
	Created_at string `json:"created_at" db:"created_at" example:"2023-09-24T17:13:42Z"`
	//YYYY-MM-DD HH:MM:SS
	Updated_at string `json:"updated_at" db:"updated_at" example:"2023-09-27T11:10:23Z"`
}

// AddNewBooking provides data to create a Booking.
//
// swagger:model
type Booking struct {
	Id      int `json:"id" db:"id" example:"1021"`
	User_id int `json:"user_id" db:"user_id" example:"906"`
	//YYYY-MM-DD HH:MM:SS
	Start_time string `json:"start_time" db:"start_time" example:"2023-10-01T12:00:00Z"`
	//YYYY-MM-DD HH:MM:SS
	End_time string `json:"end_time" db:"end_time" example:"2023-10-01T14:30:00Z"`
	//min_length = 5
	//min_length = 120
	//exclude = \"\\\/
	Comment string `json:"comment" db:"comment" example:"I may be a little late"`
}
