package project

import (
	"errors"
	"fmt"

	"github.com/thedhejavu/errauditor/examples/project/pkg/apperrors"
)

type usecase struct{}

func GetAddressByUser() error {
	address := ""
	err := errors.New("doe")
	if address == "" {
		return fmt.Errorf("unable to update appraisal by user: %w", err)
	}
	return apperrors.ErrInternalServerError("done")
}

func GetDrilldown() (error, int) {
	err := errors.New("doe")
	return fmt.Errorf("unable to update appraisal by user: %w", err), -1
}

func GetDrilldowns() (error, int) {

	return apperrors.ErrInternalServerError("done"), -1
}

func (u usecase) GetDrixxlldowns() (error, int) {
	return apperrors.ErrInternalServerError("done"), -1
}
