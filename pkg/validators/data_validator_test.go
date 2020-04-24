package validators

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func TestMaxKeyValidator(t *testing.T) {
	sourceDB, err := sql.Open("mysql", "root:@tcp(127.0.0.1:29291)/validator")
	if err != nil {
		t.Error(err)
	}
	defer sourceDB.Close()

	targetDB, err := sql.Open("mysql", "root:@tcp(127.0.0.1:29292)/validator")
	if err != nil {
		t.Error(err)
	}
	defer targetDB.Close()

	log := logrus.New()
	logger := logrus.NewEntry(log)

	validator := NewValidator(logger, sourceDB, targetDB)
	isValid := validator.validateMaxPK()
	if isValid != true {
		t.Errorf("max PK of one of the tables doesn't match to the one in target")
	}
}

func TestSchemeValidator(t *testing.T) {
	sourceDB, err := sql.Open("mysql", "root:@tcp(127.0.0.1:29291)/validator")
	if err != nil {
		t.Error(err)
	}
	defer sourceDB.Close()

	targetDB, err := sql.Open("mysql", "root:@tcp(127.0.0.1:29292)/validator")
	if err != nil {
		t.Error(err)
	}
	defer targetDB.Close()

	log := logrus.New()
	logger := logrus.NewEntry(log)

	validator := NewValidator(logger, sourceDB, targetDB)
	isValid := validator.validateSchema()
	if isValid != true {
		t.Errorf("one of the tables schema in source doesn't match to the one in target")
	}
}
