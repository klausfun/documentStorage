package repository

import (
	"context"
	"documentStorage/models"
	"encoding/json"
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

func (r *DocumentPostgres) Create(userId int, meta models.Document,
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

	fmt.Println("fileData: ", fileData)

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
