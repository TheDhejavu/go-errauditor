package examples

import (
	"errors"
	"fmt"

	"github.com/thedhejavu/go-error-analyzer/pkg/apperrors"
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
	return apperrors.ErrInternalServerError("i am here "), -1
}

func GetDrilldowns() (error, int) {
	return apperrors.ErrInternalServerError("done"), -1
}
