package database

import "database/sql"

func Ban(db *sql.DB, user string) (bool, error) {
	res, err := db.Exec("INSERT OR IGNORE INTO blacklist (user) VALUES (?)", user)
    if err != nil {
        return false, err
    }
    n, err := res.RowsAffected()
    if err != nil {
        return false, err
    }
    return n == 1, nil
}

func CheckUser(db *sql.DB, user string) (bool, error) {
    var count int

    err := db.QueryRow("SELECT COUNT(*) FROM blacklist WHERE user = ?", user).Scan(&count)
    if err != nil {
        return false, err
    }

    return count > 0, nil
}