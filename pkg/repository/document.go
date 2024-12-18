package repository

import (
	"context"
	"documentStorage/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

type DocumentPostgres struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func NewDocumentPostgres(db *sqlx.DB, redisDB *redis.Client) *DocumentPostgres {
	return &DocumentPostgres{db: db, redisClient: redisDB}
}

func (r *DocumentPostgres) Create(meta models.GetDocsResp,
	fileData []byte, jsonData string) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var metadataID int

	createMetadataQuery := fmt.Sprintf("INSERT INTO %s (name, mime, file, public)"+
		"VALUES ($1, $2, $3, $4) RETURNING id", metadataTable)
	row := tx.QueryRow(createMetadataQuery, meta.Name, meta.Mime, meta.File, meta.Public)
	if err := row.Scan(&metadataID); err != nil {
		return err
	}

	createFilesQuery := fmt.Sprintf("INSERT INTO %s (metadata_id, file_data)"+
		"VALUES ($1, $2) RETURNING id", filesTable)
	_, err = tx.Exec(createFilesQuery, metadataID, fileData)
	if err != nil {
		return err
	}

	if jsonData == "" {
		jsonData = "{}"
	}
	createJsonDocQuery := fmt.Sprintf("INSERT INTO %s (metadata_id, json_data)"+
		"VALUES ($1, $2) RETURNING id", jsonDocumentTable)
	_, err = tx.Exec(createJsonDocQuery, metadataID, jsonData)
	if err != nil {
		return err
	}

	var uid int
	for _, login := range meta.Grant {
		getUserIdQuery := fmt.Sprintf("SELECT id FROM %s WHERE login=$1", userTable)
		err = r.db.Get(&uid, getUserIdQuery, login)
		if err != nil {
			return err
		}

		createUsersMetadataQuery := fmt.Sprintf("INSERT INTO %s (user_id, metadata_id)"+
			"VALUES ($1, $2) RETURNING id", usersMetadataTable)
		_, err = tx.Exec(createUsersMetadataQuery, uid, metadataID)
		if err != nil {
			return err
		}
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	err = r.redisClient.HSet(context.Background(),
		fmt.Sprintf("document:%d", metadataID),
		"meta", metaJSON,
		"json_data", jsonData).Err()
	if err != nil {
		return err
	}

	ttl := viper.GetInt("ttlDocSec")
	err = r.redisClient.Expire(context.Background(),
		fmt.Sprintf("document:%d", metadataID), time.Duration(ttl)*time.Hour).Err()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *DocumentPostgres) GetList(userId int, docInput models.GetDocsInput) ([]models.GetDocsResp, error) {
	var login string
	if docInput.Login == nil {
		getLoginQuery := fmt.Sprintf("SELECT login FROM %s WHERE id=$1", userTable)
		err := r.db.Get(&login, getLoginQuery, userId)
		if err != nil {
			return nil, err
		}
	} else {
		login = *docInput.Login
	}

	var docs []models.GetDocsResp
	metadataQuery := fmt.Sprintf("SELECT meta.id, meta.name, meta.file,"+
		" meta.public, meta.mime, TO_CHAR(meta.created, 'YYYY-MM-DD HH24:MI:SS') AS created"+
		" FROM %s meta "+
		" INNER JOIN %s us_meta on meta.id = us_meta.metadata_id"+
		" INNER JOIN %s us on us_meta.user_id = us.id"+
		" WHERE us.login = $1 AND meta.%s = $2"+
		" LIMIT $3", metadataTable, usersMetadataTable, userTable, docInput.Key)
	err := r.db.Select(&docs, metadataQuery, login, docInput.Value, docInput.Limit)
	if err != nil {
		return nil, err
	}

	for i, doc := range docs {
		var grant []string
		loginQuery := fmt.Sprintf("SELECT us.login "+
			" FROM %s meta "+
			" INNER JOIN %s us_meta on meta.id = us_meta.metadata_id"+
			" INNER JOIN %s us on us_meta.user_id = us.id"+
			" WHERE meta.id = $1", metadataTable, usersMetadataTable, userTable)
		err = r.db.Select(&grant, loginQuery, doc.Id)
		if err != nil {
			return nil, err
		}

		docs[i].Grant = grant
	}

	return docs, err
}

func (r *DocumentPostgres) GetById(docId int) (models.GetDoc, error) {
	var res models.GetDoc
	notFound := false

	metaJSON, err := r.redisClient.HGet(context.Background(),
		fmt.Sprintf("document:%d", docId), "json_data").Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return res, err
		}

		notFound = true
	}

	if !notFound && metaJSON != "" && metaJSON != "{}" {
		res.IsFile = false
		res.JSON = metaJSON

		return res, nil
	}

	var jsonData string
	getJsonQuery := fmt.Sprintf("SELECT json_data FROM %s WHERE metadata_id=$1", jsonDocumentTable)
	err = r.db.Get(&jsonData, getJsonQuery, docId)
	if err != nil {
		return res, err
	}

	if jsonData != "{}" && jsonData != "" {
		err = r.redisClient.HSet(context.Background(),
			fmt.Sprintf("document:%d", docId), "json_data", jsonData).Err()
		if err != nil {
			return res, err
		}

		res.IsFile = false
		res.JSON = jsonData

		return res, nil
	}

	var file []byte
	getFileQuery := fmt.Sprintf("SELECT file_data FROM %s WHERE metadata_id=$1", filesTable)
	err = r.db.Get(&file, getFileQuery, docId)
	if err != nil {
		return res, err
	}

	var mime string
	getMimeQuery := fmt.Sprintf("SELECT mime FROM %s WHERE id=$1", metadataTable)
	err = r.db.Get(&mime, getMimeQuery, docId)
	if err != nil {
		return res, err
	}

	res.IsFile = true
	res.File = file
	res.MimeType = mime

	return res, nil
}

func (r *DocumentPostgres) Delete(docId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", metadataTable)
	_, err := r.db.Exec(query, docId)

	return err
}
