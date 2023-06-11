package companiespb 

import (
	"golang.org/x/exp/constraints"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func verifyString(value *string, maxLength int) error {
	if value != nil {
		if len(*value) > int(maxLength) {
			return status.Errorf(
				codes.InvalidArgument,
				"Value should be shorter than %d",
				maxLength,
			)
		}
	}
	return nil
}

func verifyInteger[T constraints.Integer](value *T, low T, high T) error {
	if value != nil {
		if *value > high || *value <= low {
			return status.Errorf(
				codes.InvalidArgument,
				"Value should be smaller than %d and greater than %d",
				high,
				low,
			)
		}
	}
	return nil
}
