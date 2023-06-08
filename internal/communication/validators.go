package communication

import (
	"golang.org/x/exp/constraints"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func verifyString(value *string, maxLength int) (*string, error) {
	if value != nil {
		if len(*value) > int(maxLength) {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"Value should be shorter than %d",
				maxLength,
			)
		}
	}
	return value, nil
}

func verifyInteger[T constraints.Integer](value *T, low T, high T) (*T, error) {
	if value != nil {
		if *value > high || *value <= low {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"Value should be smaller than %d and greater than %d",
				high,
				low,
			)
		}
	}
	return value, nil
}
