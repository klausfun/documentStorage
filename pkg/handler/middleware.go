package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

const userCtx = "userId"

func (h *Handler) userIdentity(c *gin.Context, token string) error {
	userId, err := h.services.Authorization.ParseToken(token)
	if err != nil {
		return err
	}

	c.Set(userCtx, userId)
	return nil
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	return idInt, nil
}
