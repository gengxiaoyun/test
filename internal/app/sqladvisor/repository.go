package sqladvisor

import (
	"github.com/romberli/das/global"
	"github.com/romberli/das/internal/dependency/sqladvisor"
	"github.com/romberli/go-util/middleware"
	"github.com/romberli/log"
)

var _ sqladvisor.Repository = (*Repository)(nil)

type Repository struct {
	Database middleware.Pool
}

// NewRepository returns *Repository with given middleware.Pool
func NewRepository(db middleware.Pool) *Repository {
	return &Repository{Database: db}
}

// NewRepository returns *Repository with global mysql pool
func NewRepositoryWithGlobal() *Repository {
	return NewRepository(global.DASMySQLPool)
}

// Execute executes given command and placeholders on the middleware
func (r *Repository) Execute(command string, args ...interface{}) (middleware.Result, error) {
	conn, err := r.Database.Get()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf("healthcheck DASRepo.Execute(): close database connection failed.\n%s", err.Error())
		}
	}()

	return conn.Execute(command, args...)
}

// Transaction returns a middleware.Transaction that could execute multiple commands as a transaction
func (r *Repository) Transaction() (middleware.Transaction, error) {
	return r.Database.Transaction()
}

// Save saves sql tuning advice into the middleware
func (r *Repository) Save(dbID int, sqlText, advice, message string) error {
	sql := `insert into t_sa_operation_info(db_id, sql_text, advice, message) values(?, ?, ?, ?);`
	_, err := r.Execute(sql, dbID, sqlText, advice, message)

	return err
}
