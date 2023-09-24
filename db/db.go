package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/subliker/backendproj/model"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    start_time TIMESTAMP NOT NULL,
	end_time TIMESTAMP NOT NULL,
	comment TEXT
)`

type httpCode int

type BookingsData struct {
	Count int             `json:"count"`
	Rows  []model.Booking `json:"rows"`
}

type DataBase struct {
	base *sqlx.DB
}

func (c *DataBase) Init() *sqlx.DB {
	e := godotenv.Load()
	if e != nil {
		fmt.Println(e)
	}

	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	if db_host == "" {
		db_host = "localhost"
	}
	if db_port == "" {
		db_port = "5432"
	}
	if db_user == "" {
		db_user = "postgres"
	}
	if db_name == "" {
		db_name = "postgres"
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_password, db_name)
	fmt.Println(connStr)
	c.base = sqlx.MustConnect("postgres", connStr)
	c.base.MustExec(schema)
	c.base.MustExec("SET timezone = 'Europe/Moscow'")
	return c.base
}

func (c *DataBase) AddNewUser(user model.User) (int, httpCode, error) {
	tx := c.base.MustBegin()
	var user_id int
	err := tx.QueryRow(`INSERT INTO users (username, password, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`, user.Username, user.Password, user.Created_at, user.Updated_at).Scan(&user_id)
	if err != nil {
		return -1, http.StatusInternalServerError, err
	}
	tx.Commit()
	return user_id, 200, nil
}

func (c *DataBase) AddNewBooking(booking model.Booking) (int, httpCode, error) {
	tx := c.base.MustBegin()
	var booking_id int
	err := tx.QueryRow(`INSERT INTO bookings (user_id, start_time, end_time, comment) VALUES ($1, $2, $3, $4) RETURNING id`, booking.User_id, booking.Start_time, booking.End_time, booking.Comment).Scan(&booking_id)
	if err != nil {
		return -1, http.StatusInternalServerError, err
	}
	tx.Commit()
	return booking_id, 200, nil
}

func (c *DataBase) GetUserDataByID(id int) (model.User, httpCode, error) {
	tx := c.base.MustBegin()
	var user model.User
	err := tx.QueryRowx(`SELECT * FROM users WHERE id=$1`, strconv.Itoa(id)).StructScan(&user)
	if err == sql.ErrNoRows {
		return model.User{}, 200, nil
	} else if err != nil {
		return model.User{}, http.StatusInternalServerError, err
	} else {
		return user, 200, nil
	}
}

func (c *DataBase) GetBookingDataByID(id int) (model.Booking, httpCode, error) {
	tx := c.base.MustBegin()
	var booking model.Booking
	err := tx.QueryRowx(`SELECT * FROM bookings WHERE id=$1`, id).StructScan(&booking)
	if err == sql.ErrNoRows {
		return model.Booking{}, 200, nil
	} else if err != nil {
		return model.Booking{}, http.StatusInternalServerError, err
	} else {
		return booking, 200, nil
	}
}

func (c *DataBase) GetBookings(limit, page, offset string) (BookingsData, httpCode, error) {
	tx := c.base.MustBegin()
	bookingsData := BookingsData{}
	var count int
	err := tx.QueryRow(`SELECT COUNT(*) as count FROM bookings`).Scan(&count)
	if err != nil {
		return BookingsData{}, http.StatusInternalServerError, err
	}
	bookingsData.Count = count
	bookings := make([]model.Booking, 0)
	var strQuery string
	if limit != "" && offset != "" {
		limitI, errL := strconv.Atoi(limit)
		if errL != nil {
			return BookingsData{}, http.StatusBadRequest, errL
		}

		offsetI, errO := strconv.Atoi(offset)
		if errO != nil {
			return BookingsData{}, http.StatusBadRequest, errO
		}

		strQuery = fmt.Sprintf(` SELECT * FROM bookings ORDER BY id LIMIT %d OFFSET %d`, limitI, offsetI)
	} else if limit != "" && page != "" {
		limitI, errL := strconv.Atoi(limit)
		if errL != nil {
			return BookingsData{}, http.StatusBadRequest, errL
		}

		pageI, errP := strconv.Atoi(page)
		if errP != nil {
			return BookingsData{}, http.StatusBadRequest, errP
		}

		strQuery = fmt.Sprintf(` SELECT * FROM bookings ORDER BY id LIMIT %d OFFSET %d`, limitI, limitI*(pageI-1))
	} else if limit != "" {
		limitI, errL := strconv.Atoi(limit)
		if errL != nil {
			return BookingsData{}, http.StatusBadRequest, errL
		}

		strQuery = fmt.Sprintf(` SELECT * FROM bookings ORDER BY id LIMIT %d`, limitI)
	} else {
		strQuery = " SELECT * FROM bookings"
	}
	rows, err := tx.Queryx(strQuery)
	if err != nil {
		return BookingsData{}, http.StatusInternalServerError, err
	}
	for rows.Next() {
		var booking model.Booking
		err := rows.StructScan(&booking)
		if err != nil {
			return BookingsData{}, http.StatusInternalServerError, err
		}
		bookings = append(bookings, booking)
	}

	bookingsData.Rows = bookings
	return bookingsData, 200, nil
}

func (c *DataBase) DeleteUserByID(id int) (httpCode, error) {
	tx := c.base.MustBegin()

	isExists, httpCodeE, errE := c.CheckUserExists(id)
	if errE != nil {
		return httpCodeE, errE
	}
	if !isExists {
		return http.StatusBadRequest, errors.New("user with this id doesn't exis")
	}

	_, err := tx.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return http.StatusBadRequest, err
	}

	_, err = tx.Exec("DELETE FROM bookings WHERE user_id =$1", id)
	if err != nil && err != sql.ErrNoRows {
		return http.StatusInternalServerError, err
	}
	tx.Commit()
	return 200, nil
}

func (c *DataBase) DeleteBookingByID(id int) (httpCode, error) {
	tx := c.base.MustBegin()
	_, err := tx.Exec("DELETE FROM bookings WHERE id =$1", id)
	if err != nil {
		return http.StatusBadRequest, err
	}
	tx.Commit()
	return 200, nil
}

func (c *DataBase) CheckUsernameExists(username string) (bool, httpCode, error) {
	tx := c.base.MustBegin()
	var user model.User
	err := tx.QueryRowx("SELECT * FROM users WHERE username=$1", username).StructScan(&user)
	tx.Commit()
	if err == sql.ErrNoRows {
		return false, 200, nil
	} else if err != nil {
		return false, http.StatusInternalServerError, err
	} else {
		return true, 200, nil
	}
}

func (c *DataBase) CheckUserExists(id int) (bool, httpCode, error) {
	tx := c.base.MustBegin()
	var user model.User
	err := tx.QueryRowx("SELECT * FROM users WHERE id=$1", id).StructScan(&user)
	tx.Commit()
	if err == sql.ErrNoRows {
		return false, 200, nil
	} else if err != nil {
		return false, http.StatusInternalServerError, err
	} else {
		return true, 200, nil
	}
}

func (c *DataBase) UpdateUserData(user model.User) (model.User, httpCode, error) {
	tx := c.base.MustBegin()
	_, err := tx.Exec(`UPDATE users SET username=$1, password=$2, updated_at=$3 WHERE id=$4`, user.Username, user.Password, user.Updated_at, user.Id)
	if err != nil {
		return model.User{}, http.StatusInternalServerError, err
	}

	err = tx.QueryRowx("SELECT * FROM users WHERE id=$1", user.Id).StructScan(&user)
	if err != nil {
		return model.User{}, http.StatusInternalServerError, err
	}
	tx.Commit()
	return user, 200, err
}

func (c *DataBase) UpdateBookingData(booking model.Booking) (model.Booking, httpCode, error) {
	tx := c.base.MustBegin()
	_, err := tx.Exec(`UPDATE bookings SET start_time=$1, end_time=$2, comment=$3 WHERE id=$4`, booking.Start_time, booking.End_time, booking.Comment, booking.Id)
	if err != nil {
		return model.Booking{}, http.StatusInternalServerError, err
	}

	err = tx.QueryRowx("SELECT * FROM bookings WHERE id=$1", booking.Id).StructScan(&booking)
	if err != nil {
		return model.Booking{}, http.StatusInternalServerError, err
	}
	tx.Commit()
	return booking, 200, err
}
