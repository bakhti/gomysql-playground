package validators

import (
	"database/sql"
	"fmt"
	"github.com/siddontang/go-mysql/schema"
	"github.com/sirupsen/logrus"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

type Validator struct {
	logger   *logrus.Entry
	sourceDB *sql.DB
	targetDB *sql.DB
}

func NewValidator(logger *logrus.Entry, sourceDB, targetDB *sql.DB) *Validator {
	return &Validator{
		logger:   logger,
		sourceDB: sourceDB,
		targetDB: targetDB,
	}
}

func (v *Validator) Validate() bool {
	return v.validateSchema()
}

func (v *Validator) validateSchema() bool {
	tblNames, err := showTablesFromSource(v.sourceDB)
	if err != nil {
		v.logger.WithError(err).Error("failed to show tables")
	}

	for _, tbl := range tblNames {
		sourceTable, err := schema.NewTableFromSqlDB(v.sourceDB, "validator", tbl)
		if err != nil {
			v.logger.WithError(err).Error("couldn't read from source DB")
		}
		targetTable, err := schema.NewTableFromSqlDB(v.targetDB, "validator", tbl)
		if err != nil {
			v.logger.WithError(err).Error("couldn't read from target DB")
		}
		if !reflect.DeepEqual(sourceTable, targetTable) {
			return false
		}
	}
	return true
}

func (v *Validator) validateMaxPK() bool {
	tblNames, err := showTablesFromSource(v.sourceDB)
	if err != nil {
		v.logger.WithError(err).Error("failed to show tables")
	}

	for _, tbl := range tblNames {
		sm, err := maxPK(v.sourceDB, tbl)
		if err != nil {
			v.logger.WithError(err).Error("failed to query the source DB")
		}

		tm, err := maxPK(v.targetDB, tbl)
		if err != nil {
			v.logger.WithError(err).Error("failed to query the target DB")
		}

		if sm != tm {
			return false
		}
	}
	return true
}

func maxPK(d *sql.DB, t string) (float64, error) {
	rows, err := d.Query(fmt.Sprintf("select max(id) from %s", t))
	if err != nil {
		return 0, err
	} else {
		var m float64
		for rows.Next() {
			rows.Scan(&m)
		}
		return m, nil
	}
}

func showTablesFromSource(s *sql.DB) ([]string, error) {
	rows, err := s.Query(fmt.Sprint("show tables"))
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	tables := make([]string, 0)
	for rows.Next() {
		var tbl string
		err = rows.Scan(&tbl)
		if err != nil {
			return tables, err
		}
		tables = append(tables, tbl)
	}

	return tables, nil
}

func (v *Validator) Run() error {
	//The business logic can be implemented here and then extracted to another struct of interface

	//Uncomment the lines below if you want to return error. Just for demonstration
	//err := errors.New("Error running the validator")
	//v.logger.WithError(err).Error(err.Error())
	//return err
	var db *sql.DB
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:29291)/validator")
	if err != nil {
		v.logger.Fatalf("err: %v", err)
	}
	defer db.Close()

	// Get the table description

	v.logger.WithField("from", "validator").Info("Hello I am the validator!")

	return nil
}
