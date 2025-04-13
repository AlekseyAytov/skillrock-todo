package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AlekseyAytov/skillrock-todo/internal/models/task"
	"github.com/AlekseyAytov/skillrock-todo/internal/store"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	db *sql.DB
}

// NewDBStorage creates DBStorage and checks connection
func NewDBStorage(dbDSN string) (*DBStorage, error) {
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		return nil, err
	}
	result := &DBStorage{db: db}

	err = result.checkDB()
	if err != nil {
		return nil, err
	}

	err = result.checkOrCreateTable()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *DBStorage) checkDB() error {
	err := d.db.PingContext(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

func (d *DBStorage) checkOrCreateTable() error {
	_, err := d.db.ExecContext(context.TODO(),
		`CREATE TABLE IF NOT EXISTS public.tasks (
			"id" SERIAL PRIMARY KEY,
			"title" TEXT NOT NULL,
			"description" TEXT,
			"status" TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
			"created_at" TIMESTAMP DEFAULT now(),
			"updated_at" TIMESTAMP DEFAULT now()
    	);`)
	return err
}

// Add task to database
func (d *DBStorage) Add(t task.Task) error {
	_, err := d.db.ExecContext(context.TODO(),
		`INSERT INTO public.tasks (
			"title",
			"description",
			"status",
			"created_at",
			"updated_at"
		) VALUES ($1, $2, $3, $4, $5);`,
		t.Title, t.Description, t.Status, t.CreatedAt, t.UpdatedAt)
	return err
}

func (d *DBStorage) FindBy(id string) (task.Task, error) {
	row := d.db.QueryRowContext(context.TODO(),
		`SELECT
			"id",
			"title",
			"description",
			"status",
			"created_at",
			"updated_at"
		FROM public.tasks WHERE "id" = $1;`, id)

	var t task.Task
	// If no row matches the query, Scan returns [ErrNoRows]
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, store.ErrTaskNotFound
		}
		return t, err
	}
	return t, nil
}

func (d *DBStorage) GetAll() ([]task.Task, error) {
	rows, err := d.db.QueryContext(context.TODO(),
		`SELECT
			"id",
			"title",
			"description",
			"status",
			"created_at",
			"updated_at"
		FROM public.tasks;`)
	defer rows.Close()

	tasks := make([]task.Task, 0)
	for rows.Next() {
		var t task.Task
		err = rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (d *DBStorage) Update(t task.Task) error {
	// запрос не обновляет поле created_at
	_, err := d.db.ExecContext(context.TODO(),
		`UPDATE public.tasks SET
			"title"       = $1,
			"description" = $2,
			"status"      = $3,
			"updated_at"  = $4
		WHERE "id"        = $5;`,
		t.Title, t.Description, t.Status, t.UpdatedAt, t.ID)
	return err
}

func (d *DBStorage) Delete(t task.Task) error {
	_, err := d.db.ExecContext(context.TODO(),
		`DELETE FROM public.tasks WHERE "id" = $1;`, t.ID)
	return err
}
