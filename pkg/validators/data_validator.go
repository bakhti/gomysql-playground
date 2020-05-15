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
	tblNames, err := showTablesFromSource(v.sourceDB)
	if err != nil {
		v.logger.WithError(err).Error("failed to show tables")
	}

	for _, tbl := range tblNames {
		if !v.validateSchema(tbl) || !v.validateMaxPK(tbl) {
			return false
		}
		if !v.validateSingleRow(tbl) {
			return false
		}
	}
	return true
}

// validateSchema matches the schema of the given table from source DB
// against the one in target DB
// TODO: the database name is hardcoded
func (v *Validator) validateSchema(tbl string) bool {
	sourceTable, err := schema.NewTableFromSqlDB(v.sourceDB, "validator", tbl)
	if err != nil {
		v.logger.WithError(err).Error("couldn't get schema from source DB")
	}
	targetTable, err := schema.NewTableFromSqlDB(v.targetDB, "validator", tbl)
	if err != nil {
		v.logger.WithError(err).Error("couldn't get schema from target DB")
	}
	if !reflect.DeepEqual(sourceTable, targetTable) {
		return false
	}
	return true
}

// validateMaxPK validates if the max PK of the given table from the source DB
// matches the one from target DB
func (v *Validator) validateMaxPK(tbl string) bool {
	sm, err := maxPK(v.sourceDB, tbl)
	if err != nil {
		v.logger.WithError(err).Error("failed to get max PK from source DB")
	}

	tm, err := maxPK(v.targetDB, tbl)
	if err != nil {
		v.logger.WithError(err).Error("failed to get max PK from target DB")
	}

	if sm != tm {
		return false
	}
	return true
}

// validateSingleRow validates if a random row in the given table from the source DB
// matches the same row in the target DB
func (v *Validator) validateSingleRow(tbl string) bool {
	id, scs, err := getCheckSumRand(v.sourceDB, tbl)
	if err != nil {
		v.logger.WithError(err).Error("failed to get checksum of a random row from source DB")
	}
	tcs, err := getCheckSum(v.targetDB, tbl, id)
	if err != nil {
		v.logger.WithError(err).Error("failed to get checksum of a row from target DB")
	}

	if scs != tcs {
		return false
	}
	return true
}

// getCheckSum receives DB handle, table name and row id, returns a checksum of that row
// TODO: should not be bound to id column
func getCheckSum(db *sql.DB, tbl string, id float64) (string, error) {
	// TODO: make it recognise columns
	query := fmt.Sprintf("SELECT MD5(CONCAT(id, IFNULL(data, ''))) FROM %s WHERE id = %f", tbl, id)
	var cs string

	if err := db.QueryRow(query).Scan(&cs); err != nil {
		return "", err
	}
	return cs, nil
}

// getCheckSumRand receives DB handle and table name, returns an ID and a checksum of a random row
func getCheckSumRand(db *sql.DB, tbl string) (float64, string, error) {
	t1 := fmt.Sprintf("SELECT t1.id, MD5(CONCAT(t1.id, IFNULL(t1.data, ''))) FROM %s AS t1", tbl)
	t2 := fmt.Sprintf("JOIN (SELECT CEIL(RAND() * (SELECT MAX(id) FROM %s)) AS id) AS t2", tbl)
	query := fmt.Sprintf("%s %s WHERE t1.id >= t2.id ORDER BY t1.id ASC LIMIT 1", t1, t2)

	var id float64
	var cs string
	if err := db.QueryRow(query).Scan(&id, &cs); err != nil {
		return 0, "", err
	}
	return id, cs, nil
}

// maxPK receives database handle and table name, returns max value of PK
func maxPK(db *sql.DB, tbl string) (float64, error) {
	// TODO: make it recognise what column is PK and act properly when there is no PK
	query := fmt.Sprintf("select max(id) from %s", tbl)
	var m float64
	if err := db.QueryRow(query).Scan(&m); err != nil {
		return 0, err
	}
	return m, nil
}

// showTablesFromSource receives source DB handle, returns list of tables in that DB
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
