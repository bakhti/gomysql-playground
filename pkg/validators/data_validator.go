package validators

import (
	"github.com/sirupsen/logrus"
)

type Validator struct {
	logger *logrus.Entry
}

func NewValidator(logger *logrus.Entry) *Validator {
	return &Validator{
		logger: logger,
	}
}

func (v *Validator) Run() error {
	//The business logic can be implemented here and then extracted to another struct of interface

	//Uncomment the lines below if you want to return error. Just for demonstration
	//err := errors.New("Error running the validator")
	//v.logger.WithError(err).Error(err.Error())
	//return err

	v.logger.WithField("from", "validator").Info("Hello I am the validator!")

	return nil
}
