package ingestion

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type TrsSaveFunc func(db *sql.DB, data any, l *slog.Logger) error
type SaveError struct {
	Message  string
	InnerErr error
}

func (e SaveError) Error() string {
	return fmt.Sprintf("%v, err: %+v", e.Message, e.InnerErr)
}

func (e SaveError) Unwrap() error {
	return e.InnerErr
}

type Saver interface {
	Save(data any) error
	Close() error
}

type TursoSaver struct {
	DbConn *sql.DB

	saveFuncs []TrsSaveFunc
	l         *slog.Logger
}

var (
	trsoInstance *TursoSaver
	trsoOnce     sync.Once
)

func GetTursoSaver(opts TSOptions) (*TursoSaver, error) {
	var returnErr error
	trsoOnce.Do(func() {
		opts.Logger.Debug("turso", "message", fmt.Sprintf("opening db at: %v", fmt.Sprintf("%s", opts.Url)))
		url := fmt.Sprintf("%s?authToken=%s", opts.Url, opts.Token)
		db, err := sql.Open("libsql", url)
		if err != nil {
			returnErr = err
			return
		}

		trsoInstance = &TursoSaver{}
		trsoInstance.DbConn = db
		trsoInstance.saveFuncs = append(trsoInstance.saveFuncs, opts.SaveFuncs...)
		trsoInstance.l = opts.Logger
	})

	return trsoInstance, returnErr
}

func (s TursoSaver) Close() error {
	return s.DbConn.Close()
}

func (s TursoSaver) Save(data any) error {
	for _, sf := range s.saveFuncs {
		err := sf(s.DbConn, data, s.l)
		if err == nil {
			// One save per datatype
			return nil
		}

		if errors.As(err, &SaveError{}) {
			s.l.Error("saver", "error", err)
		}
	}

	return errors.New("save func not found for provided data")
}
