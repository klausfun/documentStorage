package handler

import (
	"documentStorage/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"regexp"
)

func isValidLogin(login string) bool {
	loginRegex := regexp.MustCompile(`^[a-zA-Z0-9]{8,}$`)
	return loginRegex.MatchString(login)
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	return hasLower && hasUpper && hasDigit && hasSpecial
}

func (h *Handler) signUp(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if !(isValidLogin(input.Login) && isValidPassword(input.Password)) {
		newErrResponse(c, http.StatusBadRequest, "invalid login or password")
		return
	}

	if input.Token != os.Getenv("REGISTRATION_TOKEN") {
		newErrResponse(c, http.StatusForbidden, "no access rights")
		return
	}

	login, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		if errResp, ok := err.(*ErrorResponse); ok {
			newErrResponse(c, errResp.Code, errResp.Text)
			return
		}

		newErrResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, responseModel{
		Response: map[string]interface{}{
			"login": login,
		},
	})
}

type signInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Login, input.Password)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, responseModel{
		Response: map[string]interface{}{
			"token": token,
		},
	})
}

func (h *Handler) signOut(c *gin.Context) {

}
