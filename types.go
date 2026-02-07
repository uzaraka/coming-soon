package main

import "time"

type EmailRequest struct {
	Email string `json:"email" format:"email" required:"true"`
}
type User struct {
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}
type EmailList struct {
	Emails []User `json:"emails"`
}
