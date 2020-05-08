package validators

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func TestSchemeValidator(t *testing.T) {
	validator, err := activateValidator()
	if err != nil {
		t.Error(err)
	}
	defer validator.sourceDB.Close()
	defer validator.targetDB.Close()

	if isValid := validator.validateSchema("table1"); !isValid {
		t.Errorf("this table schema is expected to match")
	}
}

func TestMaxKeyValidator(t *testing.T) {
	validator, err := activateValidator()
	if err != nil {
		t.Error(err)
	}
	defer validator.sourceDB.Close()
	defer validator.targetDB.Close()

	if isValid := validator.validateMaxPK("table1"); !isValid {
		t.Errorf("max PK of this table is expected to match")
	}
}

func TestSingleRowValidator(t *testing.T) {
	validator, err := activateValidator()
	if err != nil {
		t.Error(err)
	}
	defer validator.sourceDB.Close()
	defer validator.targetDB.Close()

	if isValid := validator.validateSingleRow("table1"); !isValid {
		t.Errorf("all the rows of this table are expected to match")
	}
}

func activateValidator() (*Validator, error) {
	sourceDB, err := sql.Open("mysql", "root:@tcp(127.0.0.1:29291)/validator")
	if err != nil {
		return nil, err
	}

	targetDB, err := sql.Open("mysql", "root:@tcp(127.0.0.1:29292)/validator")
	if err != nil {
		return nil, err
	}

	log := logrus.New()
	logger := logrus.NewEntry(log)

	return NewValidator(logger, sourceDB, targetDB), nil
}
