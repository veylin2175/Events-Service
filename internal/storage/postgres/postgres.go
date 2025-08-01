package postgres

import (
	"Events-Service/internal/config"
	"Events-Service/internal/models"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func InitDB(cfg *config.Config) (*Storage, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("couldn't connect to the DB: %v", err)
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveEvent(userID int64, dateStr, text string) (int64, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)", userID).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("failed to check user: %v", err)
	}
	if !exists {
		return 0, fmt.Errorf("user with ID %d not found", userID)
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, fmt.Errorf("invalid date format: %v", err)
	}

	var eventID int64
	err = s.db.QueryRow(
		`INSERT INTO event (user_id, date, text) 
         VALUES ($1, $2, $3) RETURNING id`,
		userID, date, text,
	).Scan(&eventID)

	if err != nil {
		return 0, fmt.Errorf("failed to save event: %v", err)
	}

	return eventID, nil
}

func (s *Storage) UpdateEvent(userID, eventID int64, dateStr, text string) error {
	query := "UPDATE event SET"
	args := []interface{}{}
	argPos := 1

	if dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return fmt.Errorf("invalid date format: %v", err)
		}
		query += fmt.Sprintf(" date = $%d,", argPos)
		args = append(args, date)
		argPos++
	}

	if text != "" {
		query += fmt.Sprintf(" text = $%d,", argPos)
		args = append(args, text)
		argPos++
	}

	if len(args) == 0 {
		return nil
	}

	query = strings.TrimSuffix(query, ",")

	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", argPos, argPos+1)
	args = append(args, eventID, userID)

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update event: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("event not found or access denied")
	}

	return nil
}

func (s *Storage) DeleteEvent(userID, eventID int64) error {
	err := s.db.QueryRow(
		"DELETE FROM event WHERE id = $1 AND user_id = $2",
		eventID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete event: %v", err)
	}

	return nil
}

func (s *Storage) GetEventsByDay(userID int64, day string) ([]models.Event, error) {
	date, err := time.Parse("2006-01-02", day)
	rows, err := s.db.Query(
		`SELECT id, date, text FROM event 
         WHERE user_id = $1 AND date = $2 
         ORDER BY date`,
		userID, date,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily events: %v", err)
	}
	defer rows.Close()

	return scanEvents(rows)
}

func (s *Storage) GetEventsByWeek(userID int64, startOfWeek time.Time) ([]models.Event, error) {
	endOfWeek := startOfWeek.AddDate(0, 0, 7)
	rows, err := s.db.Query(
		`SELECT id, date, text FROM event 
         WHERE user_id = $1 AND date >= $2 AND date < $3 
         ORDER BY date`,
		userID,
		startOfWeek.Format("2006-01-02"),
		endOfWeek.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly events: %v", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	return scanEvents(rows)
}

func (s *Storage) GetEventsByMonth(userID int64, year int, month time.Month) ([]models.Event, error) {
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)
	rows, err := s.db.Query(
		`SELECT id, date, text FROM event 
         WHERE user_id = $1 AND date >= $2 AND date < $3 
         ORDER BY date`,
		userID,
		startOfMonth.Format("2006-01-02"),
		endOfMonth.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly events: %v", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	return scanEvents(rows)
}

func (s *Storage) CreateUser() (int64, error) {
	var userID int64
	err := s.db.QueryRow(
		"INSERT INTO users DEFAULT VALUES RETURNING user_id",
	).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("failed to create user: %v", err)
	}
	return userID, nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func scanEvents(rows *sql.Rows) ([]models.Event, error) {
	var events []models.Event
	for rows.Next() {
		var e models.Event
		var eventID int64
		var eventDate time.Time
		if err := rows.Scan(&eventID, &eventDate, &e.Text); err != nil {
			return nil, err
		}
		e.Date = eventDate.Format("2006-01-02")
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
