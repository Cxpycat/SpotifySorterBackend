package mysql

import (
	userModel "SpotifySorter/models"
	"database/sql"
	"errors"
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
			spotify_access_token TEXT,
			country VARCHAR(2),
			id_spotify VARCHAR(255) UNIQUE NOT NULL,
			product VARCHAR(20),
			access_token TEXT);
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

func (s *Storage) SaveUser(email, accessToken, spotifyAccessToken, country, name, idSpotify, product string) (*userModel.User, error) {
	const op = "storage.mysql.SaveUser"

	stmt, err := s.db.Prepare(`
        INSERT INTO users(email, access_token, spotify_access_token, country, name, id_spotify, product)
        VALUES(?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(email, accessToken, spotifyAccessToken, country, name, idSpotify, product)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	user := &userModel.User{
		Email:              email,
		AccessToken:        accessToken,
		SpotifyAccessToken: spotifyAccessToken,
		Country:            country,
		Name:               name,
		IdSpotify:          idSpotify,
		Product:            product,
	}

	return user, nil
}

func (s *Storage) GetUserByEmail(email string) (*userModel.User, error) {
	stmt, err := s.db.Prepare(`SELECT * FROM users WHERE email = ?`)

	if err != nil {
		return nil, err
	}

	var user userModel.User
	err = stmt.QueryRow(email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.SpotifyAccessToken,
		&user.Country,
		&user.IdSpotify,
		&user.Product,
		&user.AccessToken,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}

func (s *Storage) GetUserByAccessToken(accessToken string) (*userModel.User, error) {
	stmt, err := s.db.Prepare(`SELECT * FROM users WHERE access_token = ?`)

	if err != nil {
		return nil, err
	}

	var user userModel.User
	err = stmt.QueryRow(accessToken).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.SpotifyAccessToken,
		&user.Country,
		&user.IdSpotify,
		&user.Product,
		&user.AccessToken,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil

}

func (s *Storage) UpdateUser(email, spotifyAccessToken, country, name, idSpotify, product string) (*userModel.User, error) {
	const op = "storage.mysql.UpdateUser"

	stmt, err := s.db.Prepare(`
        UPDATE users
        SET spotify_access_token = ?, country = ?, name = ?, id_spotify = ?, product = ?
        WHERE email = ?;
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(spotifyAccessToken, country, name, idSpotify, product, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
