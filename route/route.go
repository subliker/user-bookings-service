package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	dv "github.com/subliker/backendproj/datavalidator"
	"github.com/subliker/backendproj/db"
	"github.com/subliker/backendproj/model"

	"github.com/gin-gonic/gin"
)

var DataBase db.DataBase
var DB = DataBase.Init()

// AddNewUser godoc
//
//	@Summary		Add new user data in db
//	@Description	Prepairing user data for new user in db
//	@Tags			user
//
//	@Param   username   formData   string     true        "username (3 <= length <= 20, exclude=\"\\\/")"
//	@Param   password   formData   string     true        "password (6 <= length <= 20, exclude=\"\\\/")"
//
//	@Success		200				{object}	model.User
//	@Failure		400				{object}	dv.ResError
//	@Failure		500				{object}	dv.ResError
//	@Router			/user [post]
func AddNewUser(c *gin.Context) {
	var user model.User
	fmt.Println(c.PostForm("username"))
	err := dv.ValidateUsername(c.PostForm("username"))
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}
	user.Username = c.PostForm("username")

	err = dv.ValidatePassword(c.PostForm("password"))
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}
	passwordHashed, errh := dv.HashPassword(c.PostForm("password"))
	if errh != nil {
		dv.ResMessage(c, http.StatusInternalServerError, dv.ErrToString(errh))
		return
	}
	user.Password = passwordHashed

	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")
	user.Created_at = ts
	user.Updated_at = ts

	user_id, httpCodeA, errA := DataBase.AddNewUser(user)
	if errA != nil {
		dv.ResMessage(c, int(httpCodeA), dv.ErrToString(errA))
		return
	}

	user, httpCodeG, errG := DataBase.GetUserDataByID(user_id)
	if errG != nil {
		dv.ResMessage(c, int(httpCodeG), dv.ErrToString(errG))
		return
	}

	userData, e := json.Marshal(user)
	if e != nil {
		dv.ResMessage(c, http.StatusInternalServerError, dv.ErrToString(e))
		return
	}

	c.Data(http.StatusOK, "application/json", userData)
}

// GetUserDataById godoc
//
//	@Summary		Return user data (json) by id
//	@Description	If user isn't found, it returns blank json
//	@Tags			user
//	@Produce		json
//	@Param id path int required "id to find user"
//	@Success		200				{object}	model.User
//	@Failure		400				{object}	dv.ResError
//	@Failure		500				{object}	dv.ResError
//	@Router			/user/{id} [get]
func GetUserDataById(c *gin.Context) {
	id := c.Param("id")
	if c.Param("id") == "" {
		dv.ResMessage(c, http.StatusBadRequest, "id isn`t set")
		return
	}
	idI, err := strconv.Atoi(id)
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}

	user, httpCodeG, errG := DataBase.GetUserDataByID(idI)
	if err != nil {
		dv.ResMessage(c, int(httpCodeG), dv.ErrToString(errG))
		return
	}

	if user == (model.User{}) {
		c.Data(http.StatusOK, "application/json", []byte("{}"))
		return
	}
	userData, e := json.Marshal(user)
	if e != nil {
		fmt.Println(e)
	}
	c.Data(http.StatusOK, "application/json", userData)
}

// DeleteUserDataById godoc
//
//	@Summary		Delete user data (user and bookings) by id
//	@Description	Delete user and bookings (if exists)
//	@Tags			user
//	@Produce		json
//	@Param id path int required "id to find user"
//	@Success		200				{object}	dv.ResMesOK
//	@Failure		400				{object}	dv.ResError
//	@Failure		500				{object}	dv.ResError
//	@Router			/user/{id} [delete]
func DeleteUserDataByID(c *gin.Context) {
	id := c.Param("id")
	if c.Param("id") == "" {
		dv.ResMessage(c, http.StatusBadRequest, "id isn`t set")
		return
	}
	idI, err := strconv.Atoi(id)
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}

	httpCodeD, errD := DataBase.DeleteUserByID(idI)
	if errD != nil {
		dv.ResMessage(c, int(httpCodeD), dv.ErrToString(errD))
		return
	}

	dv.ResMessage(c, http.StatusOK, "user was successfully deleted")
}

// UpdateUserDataById godoc
//
// @Summary	Update user data by id
// @Description (option) update username(to unique), password
// @Tags user
// @Produce json
// @Param   id   path   int     true        "user id"
// @Param   username   formData   string     false        "username (3 <= length <= 20, exclude=\"\\\/")"
// @Param   password   formData   string     false        "password (6 <= length <= 20, exclude=\"\\\/")"
// @Success		200				{object}	model.User
// @Failure		400				{object}	dv.ResError
// @Failure		500				{object}	dv.ResError
// @Router /user/{id} [put]
func UpdateUserDataById(c *gin.Context) {
	id := c.Param("id")
	if c.Param("id") == "" {
		dv.ResMessage(c, http.StatusBadRequest, "id isn`t set")
		return
	}
	idI, err := strconv.Atoi(id)
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}

	var user model.User
	user, httpCodeG, errG := DataBase.GetUserDataByID(idI)
	if errG != nil {
		dv.ResMessage(c, int(httpCodeG), dv.ErrToString(err))
		return
	}

	if c.PostForm("username") != "" {
		usernameExists, httpCode, err := DataBase.CheckUsernameExists(c.PostForm("username"))
		if err != nil {
			dv.ResMessage(c, int(httpCode), dv.ErrToString(err))
			return
		}
		if usernameExists {
			dv.ResMessage(c, http.StatusBadRequest, "username already exists")
			return
		}
		err = dv.ValidateUsername(c.PostForm("username"))
		if err != nil {
			dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
			return
		}
		user.Username = c.PostForm("username")
	}

	if c.PostForm("password") != "" {
		err := dv.ValidatePassword(c.PostForm("password"))
		if err != nil {
			dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
			return
		}
		passwordHashed, errh := dv.HashPassword(c.PostForm("password"))
		if errh != nil {
			dv.ResMessage(c, http.StatusInternalServerError, dv.ErrToString(errh))
			return
		}
		user.Password = passwordHashed
	}

	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")
	user.Updated_at = ts
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}

	user, httpCodeU, errU := DataBase.UpdateUserData(user)
	if err != nil {
		dv.ResMessage(c, int(httpCodeU), dv.ErrToString(errU))
		return
	}

	if user == (model.User{}) {
		dv.ResMessage(c, http.StatusBadRequest, "User wasn't found")
		return
	}
	userData, e := json.Marshal(user)
	if e != nil {
		dv.ResMessage(c, http.StatusInternalServerError, dv.ErrToString(e))
		return
	}

	c.Data(http.StatusOK, "application/json", userData)
}

// AddNewBooking godoc
//
//	@Summary		Add new booking data in db
//	@Description	Prepairing booking data for new booking (linked to user) in db
//	@Tags			booking
//
//	@Param   user_id   formData   int     true        "user_id (user is exists)"
//	@Param   start_time   formData   string     true        "start_time (YYYY-MM-DD HH:MM:SS)"
//	@Param   end_time   formData   string     true        "end_time (YYYY-MM-DD HH:MM:SS)"
//	@Param   comment   formData   string     false        "comment (5 <= length <= 120, exclude=\"\\\/")"
//
//	@Success		200				{object}	model.Booking
//	@Failure		400				{object}	dv.ResError
//	@Failure		500				{object}	dv.ResError
//	@Router			/booking [post]
func AddNewBooking(c *gin.Context) {
	var booking model.Booking
	user_id := c.PostForm("user_id")
	if user_id == "" {
		dv.ResMessage(c, http.StatusBadRequest, "user_id isn`t set")
		return
	}
	user_idI, err := strconv.Atoi(user_id)
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}

	user, httpCodeG, errG := DataBase.GetUserDataByID(user_idI)
	if errG != nil {
		dv.ResMessage(c, int(httpCodeG), dv.ErrToString(errG))
		return
	}
	if user == (model.User{}) {
		dv.ResMessage(c, http.StatusBadRequest, "user with this user_id doesn't exist")
		return
	}
	booking.User_id = user_idI

	booking.Start_time = c.PostForm("start_time")
	booking.End_time = c.PostForm("end_time")
	err = dv.CheckCorrectTimeDuration(booking.Start_time, booking.End_time)
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}

	comment := c.PostForm("comment")
	booking.Comment = comment
	err = dv.ValidateComment(comment)
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}

	booking_id, httpCodeA, errA := DataBase.AddNewBooking(booking)
	if errA != nil {
		dv.ResMessage(c, int(httpCodeA), dv.ErrToString(errA))
		return
	}

	booking, httpCodeGN, errGN := DataBase.GetBookingDataByID(booking_id)
	if errGN != nil {
		dv.ResMessage(c, int(httpCodeGN), dv.ErrToString(errGN))
		return
	}

	bookingData, e := json.Marshal(booking)
	if e != nil {
		dv.ResMessage(c, http.StatusInternalServerError, dv.ErrToString(e))
		return
	}

	c.Data(http.StatusOK, "application/json", bookingData)
}

// GetBookingDataById godoc
//
//	@Summary		Return booking data (json) by id
//	@Description	If booking isn't found, it returns blank json
//	@Tags			booking
//	@Produce		json
//	@Param id path int required "id to find booking"
//	@Success		200				{object}	model.Booking
//	@Failure		400				{object}	dv.ResError
//	@Failure		500				{object}	dv.ResError
//	@Router			/booking/{id} [get]
func GetBookingDataById(c *gin.Context) {
	var booking model.Booking
	id := c.Param("id")
	if c.Param("id") == "" {
		dv.ResMessage(c, http.StatusBadRequest, "id isn`t set")
		return
	}
	idI, err := strconv.Atoi(id)
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}

	booking, httpCodeG, errG := DataBase.GetBookingDataByID(idI)
	if errG != nil {
		dv.ResMessage(c, int(httpCodeG), dv.ErrToString(errG))
		return
	}

	if booking == (model.Booking{}) {
		c.Data(http.StatusOK, "application/json", []byte("{}"))
		return
	}
	bookingData, e := json.Marshal(booking)
	if e != nil {
		fmt.Println(e)
		return
	}
	c.Data(http.StatusOK, "application/json", bookingData)
}

// DeleteBookingById godoc
//
//	@Summary		Delete booking data by id
//
//	@Tags			booking
//	@Produce		json
//	@Param id path int required "id to find booking"
//	@Success		200				{object}	dv.ResMesOK
//	@Failure		400				{object}	dv.ResError
//	@Failure		500				{object}	dv.ResError
//	@Router			/booking/{id} [delete]
func DeleteBookingByID(c *gin.Context) {
	id := c.Param("id")
	if c.Param("id") == "" {
		dv.ResMessage(c, http.StatusBadRequest, "id isn`t set")
		return
	}
	idI, err := strconv.Atoi(id)
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}

	httpCodeD, errD := DataBase.DeleteBookingByID(idI)
	if errD != nil {
		dv.ResMessage(c, int(httpCodeD), dv.ErrToString(errD))
		return
	}

	dv.ResMessage(c, http.StatusOK, "booking was successfully deleted")
}

// GetBookings godoc
//
//	@Summary		Return all bookings
//	@Description	(optional) set limit or limit with page or limit with offset
//	@Tags			booking
//	@Produce		json
//	@Param        limit    query     int  false  "limit"
//	@Param        page    query     int  false  "page"
//	@Param        offset    query     int  false  "offset"
//	@Success		200				{object}	db.BookingsData
//	@Failure		400				{object}	dv.ResError
//	@Failure		500				{object}	dv.ResError
//	@Router			/booking [get]
func GetBookings(c *gin.Context) {
	bookings, httpCode, err := DataBase.GetBookings(c.Query("limit"), c.Query("page"), c.Query("offset"))
	if err != nil {
		dv.ResMessage(c, int(httpCode), dv.ErrToString(err))
		return
	}

	jsonData, e := json.Marshal(bookings)
	if e != nil {
		dv.ResMessage(c, http.StatusInternalServerError, dv.ErrToString(e))
		return
	}

	c.Data(int(httpCode), "appplication/json", jsonData)
}

// UpdateBookingDataById godoc
//
// @Summary	Update booking data by id
// @Description (option) update start_time, end_time, comment
// @Tags booking
// @Produce json
// @Param   id   path   int     true        "booking id"
// @Param   start_time   formData   string     false        "start_time (YYYY-MM-DD HH:MM:SS)"
// @Param   end_time   formData   string     false        "end_time (YYYY-MM-DD HH:MM:SS)"
// @Param   comment   formData   string     false        "comment (5 <= length <= 120, exclude=\"\\\/")"
// @Success		200				{object}	model.Booking
// @Failure		400				{object}	dv.ResError
// @Failure		500				{object}	dv.ResError
// @Router /booking/{id} [put]
func UpdateBookingDataById(c *gin.Context) {
	id := c.Param("id")
	if c.Param("id") == "" {
		dv.ResMessage(c, http.StatusBadRequest, "id isn`t set")
		return
	}
	idI, err := strconv.Atoi(id)
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}

	var booking model.Booking
	booking, httpCodeG, errG := DataBase.GetBookingDataByID(idI)
	if errG != nil {
		dv.ResMessage(c, int(httpCodeG), dv.ErrToString(err))
		return
	}

	if c.PostForm("start_time") != "" {
		booking.Start_time = c.PostForm("start_time")
	}
	if c.PostForm("end_time") != "" {
		booking.End_time = c.PostForm("end_time")
	}

	err = dv.CheckCorrectTimeDuration(booking.Start_time, booking.End_time)
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}

	err = dv.ValidateComment(c.PostForm("comment"))
	if err != nil {
		dv.ResMessage(c, http.StatusBadRequest, dv.ErrToString(err))
		return
	}
	booking.Comment = c.PostForm("comment")

	booking, httpCodeU, errU := DataBase.UpdateBookingData(booking)
	if err != nil {
		dv.ResMessage(c, int(httpCodeU), dv.ErrToString(errU))
		return
	}

	if booking == (model.Booking{}) {
		dv.ResMessage(c, http.StatusBadRequest, "Booking wasn't found")
		return
	}
	bookingData, e := json.Marshal(booking)
	if e != nil {
		dv.ResMessage(c, http.StatusInternalServerError, dv.ErrToString(e))
		return
	}

	c.Data(http.StatusOK, "application/json", bookingData)
}
