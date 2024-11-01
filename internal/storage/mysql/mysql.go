package mysql

import (
	UserModel "SpotifySorter/models"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func Init(cfg Config) (*Storage, error) {
	const op = "storage.mysql.New"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		   CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			access_token TEXT,
			country VARCHAR(50),
			href VARCHAR(255),
			id_spotify VARCHAR(255),
			product VARCHAR(100),
			uri VARCHAR(255)
		);
`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(email, accessToken, country, name, href, idSpotify, product, uri string) (*UserModel.User, error) {
	const op = "storage.mysql.SaveUser"

	stmt, err := s.db.Prepare(`
        INSERT INTO users(email, access_token, country, name, href, id_spotify, product, uri)
        VALUES(?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(email, accessToken, country, name, href, idSpotify, product, uri)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	user := &UserModel.User{
		Id:          id,
		Email:       email,
		AccessToken: accessToken,
		Country:     country,
		Name:        name,
		Href:        href,
		IdSpotify:   idSpotify,
		Product:     product,
		Uri:         uri,
	}

	return user, nil
}

func (s *Storage) GetUserByEmail(email string) (*UserModel.User, error) {
	const op = "storage.mysql.GetUserByEmail"

	stmt, err := s.db.Prepare(`
        SELECT id, name, email, access_token, country, href, id_spotify, product, uri 
        FROM users
        WHERE email = ?;
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var user UserModel.User

	row := stmt.QueryRow(email)
	err = row.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.AccessToken,
		&user.Country,
		&user.Href,
		&user.IdSpotify,
		&user.Product,
		&user.Uri,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) UpdateUser(email, accessToken, country, name, href, product, uri string) (*UserModel.User, error) {
	const op = "storage.mysql.UpdateUser"

	stmt, err := s.db.Prepare(`
        UPDATE users
        SET access_token = ?, country = ?, name = ?, href = ?, product = ?, uri = ?
        WHERE email = ?;
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(accessToken, country, name, href, product, uri, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
