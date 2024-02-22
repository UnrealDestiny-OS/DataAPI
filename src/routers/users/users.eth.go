package users

import (
	"strconv"
)

func GenerateSignValidationMessage(wallet string, chain int) string {
	return "Login attempt from " + wallet + " on unrealdestiny.com using the chain " + strconv.Itoa(chain)
}
