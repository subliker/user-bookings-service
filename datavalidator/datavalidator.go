package datavalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// swagger:model
type ResError struct {
	//exclude = \"\\\/
	Message string `json:"message" example:"banned symbols in username"`
}

// swagger:model
type ResMesOK struct {
	//exclude = \"\\\/
	Message string `json:"message" example:"... successfully ..." `
}

func ResMessage(c *gin.Context, httpStatus int, message string) {
	var err ResError
	err.Message = message
	errorData, e := json.Marshal(err)
	if e != nil {
		fmt.Println(e)
		return
	}
	c.Data(httpStatus, "application/json", errorData)
}

func ErrToString(err error) string {
	replacer := strings.NewReplacer(`"`, `\"`, `\`, `\\`, `/`, `\/`)
	return replacer.Replace(fmt.Sprint(err))
}

func CheckCorrectTimeDuration(t1, t2 string) error {
	replacer := strings.NewReplacer("T", " ", "Z", "")
	t1 = replacer.Replace(t1)
	t2 = replacer.Replace(t2)
	start_timeP, errS := time.Parse("2006-01-02 15:04:05", t1)
	if errS != nil {
		return errors.New("incorrect start_time")
	}

	end_timeP, errE := time.Parse("2006-01-02 15:04:05", t2)
	if errE != nil {
		return errors.New("incorrect end_time")
	}

	if !end_timeP.After(start_timeP) {
		err := errors.New("incorrect time duration")
		return err
	}

	return nil
}

func ValidateUsername(username string) error {
	if strings.Contains(username, `"`) || strings.Contains(username, `\`) || strings.Contains(username, `/`) {
		return errors.New("banned symbols in username")
	}
	if len(username) > 20 || len(username) < 3 {
		return errors.New("incorrect username length (3 <= length <= 20)")
	}
	fmt.Println(username)
	return nil
}

func ValidatePassword(password string) error {
	if strings.Contains(password, `"`) || strings.Contains(password, `\`) || strings.Contains(password, `/`) {
		return errors.New("banned symbols in password")
	}
	if len(password) > 20 || len(password) < 6 {
		return errors.New("incorrect password length (6 <= length <= 20)")
	}
	return nil
}

func ValidateComment(comment string) error {
	if strings.Contains(comment, `"`) || strings.Contains(comment, `\`) || strings.Contains(comment, `/`) {
		return errors.New("banned symbols in comments")
	}
	if (len(comment) > 120 || len(comment) < 5) && (len(comment) != 0) {
		return errors.New("incorrect comment length (5 <= length <= 120)")
	}
	return nil
}
