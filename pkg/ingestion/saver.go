package ingestion

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

type TrsSaveFunc func(db *sql.DB, data any) error

type Saver interface {
	Save(data any) error
	Close() error
}

type TursoSaver struct {
	DbConn *sql.DB

	saveFuncs []TrsSaveFunc
}

var (
	trsoInstance *TursoSaver
	trsoOnce     sync.Once
)

func GetTursoSaver(opts TSOptions) (*TursoSaver, error) {
	var returnErr error
	trsoOnce.Do(func() {
		url := fmt.Sprintf("libsql://%s?authToken=%s", opts.Url, opts.Token)
		db, err := sql.Open("libsql", url)
		if err != nil {
			returnErr = err
			return
		}

		trsoInstance.DbConn = db
		trsoInstance.saveFuncs = append(trsoInstance.saveFuncs, opts.SaveFuncs...)
	})

	return trsoInstance, returnErr
}

func (s TursoSaver) Close() error {
	return s.DbConn.Close()
}

func (s TursoSaver) Save(data any) error {
	for _, sf := range s.saveFuncs {
		err := sf(s.DbConn, data)
		if err == nil {
			// One save per datatype
			return nil
		}
	}

	return errors.New("save func not found for provided data")
}
