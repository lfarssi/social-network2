package request

import (
	"regexp"
	"social-network/pkg/model/response"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Register struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Birth     string `json:"birth"`
	Nickname  string `json:"nickname"`
	About     string `json:"about"`
	Avatar    string `json:"avatar"`
}

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	nameRegex     = regexp.MustCompile(`^[a-zA-Z]{1,50}$`)
	nicknameRegex = regexp.MustCompile(`^[a-zA-Z0-9.]{3,25}$`)
)

func (r *Register) Validate() (err *response.RegisterError) {
	err = &response.RegisterError{}
	err.Code = 400
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.Firstname = strings.TrimSpace(r.Firstname)
	r.Lastname = strings.TrimSpace(r.Lastname)
	r.Nickname = strings.ToLower(strings.TrimSpace(r.Nickname))
	r.About = strings.TrimSpace(r.About)

	if !emailRegex.MatchString(r.Email) {
		err.Field = "email"
		err.Cause = "invalid email format"
		return
	}

	if !nameRegex.MatchString(r.Firstname) {
		err.Field = "firstName"
		err.Cause = "only letters allowed, 1-50 chars"
		return
	}

	if !nameRegex.MatchString(r.Lastname) {
		err.Field = "lastName"
		err.Cause = "only letters allowed, 1-50 chars"
		return
	}

	if len(r.Nickname) > 0 && !nicknameRegex.MatchString(r.Nickname) {
		err.Field = "nickname"
		err.Cause = "3-25 chars, letters/numbers/dots only"
		return
	}

	if len(r.About) > 1000 {
		err.Field = "about"
		err.Cause = "about section too long (max 1000 characters)"
		return
	}

	if passErr := r.ValidatePassword(); passErr != "" {
		err.Field = "password"
		err.Cause = passErr
		return
	}

	if birthErr := r.ValidateBirth(); birthErr != "" {
		err.Field = "birth"
		err.Cause = birthErr
		return
	}
	if hashErr := r.HashPassword(); hashErr != "" {
		err.Code = 500
		err.Field = "password"
		err.Cause = hashErr
	}
	return nil
}

func (r *Register) ValidateBirth() string {
	now := time.Now()
	newTime, err := time.Parse(time.DateOnly, r.Birth)
	if err != nil {
		return "invalid birth date format, use mm-DD-YYYY"
	}
	if newTime.After(now) {
		return "birth date must be in the past"
	}
	if !(newTime.After(now.Add(-time.Hour*24*365*100)) && newTime.Before(now.Add(-time.Hour*24*365*18))) {
		return "birth date must be Between 18 and 100 years"
	}

	return ""
}

func (r *Register) HashPassword() string {
	passwd, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		//log in error file
		return "Can't hash password"
	}
	r.Password = string(passwd)
	return ""
}

func (r *Register) ValidatePassword() (err string) {
	if len(r.Password) < 6 || len(r.Password) > 100 {
		return "password must be 6-100 characters"
	}
	hasLetter := false
	hasDigit := false
	for _, c := range r.Password {
		switch {
		case c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z':
			hasLetter = true
		case c >= '0' && c <= '9':
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return "must contain both letters and digits"
	}

	return ""
}
