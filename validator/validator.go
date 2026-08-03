package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/google/uuid"
)

func ValidateRequestCreateEstate(request generated.CreateEstateRequest) error {

	length := request.Length
	width := request.Width

	if err := validParameter("width", width, 1, 5000); err != nil {
		return err
	}

	if err := validParameter("length", length, 1, 5000); err != nil {
		return err
	}
	return nil
}

func ValidateRequestCreateTree(req generated.CreateTreeRequest, estateId string) error {

	//check estateId
	err := isValidId(estateId)
	if err != nil {
		return err
	}

	x := req.X
	y := req.Y
	height := req.Height
	if err := validParameter("x", x, 1, 5000); err != nil {
		return err
	}

	if err := validParameter("y", y, 1, 5000); err != nil {
		return err
	}

	if err := validParameter("height", height, 1, 30); err != nil {
		return err
	}

	return nil
}

func isValidId(estateId string) error {
	isValidUuid := func(s string) bool {
		id, err := uuid.Parse(s)
		version := id.Version()
		return err == nil && version == 4
	}
	if strings.TrimSpace(estateId) == "" || !isValidUuid(estateId) {
		return errors.New("400|must valid estate id")
	}
	return nil
}

func ValidateEstateId(estateId string) error {
	//check estateId
	err := isValidId(estateId)
	if err != nil {
		return err
	}

	return nil
}

func validParameter(parameterName string, value int, minValue int, maxValue int) error {
	if value < minValue || value > maxValue {
		return errors.New(fmt.Sprintf("%s must be between %d and %d", parameterName, minValue, maxValue))
	}
	return nil
}
