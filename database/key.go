package database

import (
	"database/sql"
)

func SaveKey(db *sql.DB, keytype, user, key string) error {
	_, err := db.Exec("INSERT INTO keys (type, user, key) VALUES (?, ?, ?)", keytype, user, key)
	if err != nil {
		return err
	}
	return nil
}

func CheckKey(db *sql.DB, key string) (bool, error) {
    var count int

    err := db.QueryRow("SELECT COUNT(*) FROM keys WHERE key = ?", key).Scan(&count)
    if err != nil {
        return false, err
    }

    return count > 0, nil
}

func KeyType(db *sql.DB, key string) (string, error) {
	var keytype string
    err := db.QueryRow(`SELECT type FROM keys WHERE "key" = ?`, key).Scan(&keytype)
	if err != nil{
		return "", err
	}
	return keytype, nil
}

func KeyOwner(db *sql.DB, key string) (string, error) {
	var keyowner string
    err := db.QueryRow(`SELECT user FROM keys WHERE "key" = ?`, key).Scan(&keyowner)
	if err != nil{
		return "", err
	}
	return keyowner, nil
}