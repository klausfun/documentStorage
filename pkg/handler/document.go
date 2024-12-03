package handler

import (
	"documentStorage/models"
	"documentStorage/pkg"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

const maxSize = 2 * 1024 * 1024

func (h *Handler) createDocument(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(maxSize); err != nil {
		newErrResponse(c, http.StatusBadRequest, "file too large or an incorrect MultipartForm")
		return
	}

	metaStr := c.Request.FormValue("meta")
	if metaStr == "" {
		newErrResponse(c, http.StatusBadRequest, "meta is required")
		return
	}

	meta := models.Document{}
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		newErrResponse(c, http.StatusBadRequest, "invalid meta format")
		return
	}

	token := meta.Token
	if err := h.userIdentity(c, token); err != nil {
		var errResp *pkg.ErrorResponse
		if errors.As(err, &errResp) {
			newErrResponse(c, errResp.Code, errResp.Text)
			return
		}

		newErrResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var fileData []byte
	jsonData := make(map[string]interface{})

	if meta.File {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			newErrResponse(c, http.StatusBadRequest, "file is required")
			return
		}
		defer file.Close()

		if header.Header.Get("Content-Type") != meta.Mime {
			newErrResponse(c, http.StatusBadRequest, "invalid file type")
			return
		}

		fileData, err = io.ReadAll(file)
		if err != nil {
			newErrResponse(c, http.StatusInternalServerError, "failed to read file")
			return
		}
	}

	jsonStr := c.Request.FormValue("json")
	if jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
			newErrResponse(c, http.StatusBadRequest, "invalid json format")
			return
		}
	}

	dos := models.GetDocsResp{
		Id:     meta.Id,
		Name:   meta.Name,
		Mime:   meta.Mime,
		File:   meta.File,
		Public: meta.Public,
		Grant:  meta.Grant,
	}

	err := h.services.Document.Create(dos, fileData, jsonStr)
	if err != nil {
		var errResp *pkg.ErrorResponse
		if errors.As(err, &errResp) {
			newErrResponse(c, errResp.Code, errResp.Text)
			return
		}

		newErrResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dataModel{
		Data: map[string]interface{}{
			"json": jsonData,
			"file": meta.Name,
		},
	})
}

func (h *Handler) getListOfDocs(c *gin.Context) {
	var input models.GetDocsInput
	if err := c.BindJSON(&input); err != nil {
		newErrResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token := input.Token
	if err := h.userIdentity(c, token); err != nil {
		var errResp *pkg.ErrorResponse
		if errors.As(err, &errResp) {
			newErrResponse(c, errResp.Code, errResp.Text)
			return
		}

		newErrResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, "error getting the user ID")
		return
	}

	docs, err := h.services.Document.GetListOfDocs(userId, input)
	if err != nil {
		var errResp *pkg.ErrorResponse
		if errors.As(err, &errResp) {
			newErrResponse(c, errResp.Code, errResp.Text)
			return
		}

		newErrResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if c.Request.Method == http.MethodHead {
		if len(docs) == 0 {
			c.Status(http.StatusNotFound)
			return
		}

		c.Status(http.StatusOK)
		return
	}

	c.JSON(http.StatusOK, dataModel{
		Data: map[string]interface{}{
			"docs": docs,
		},
	})
}

func (h *Handler) getDocumentById(c *gin.Context) {

}

func (h *Handler) deleteDocument(c *gin.Context) {

}
