package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func WrapperTimeStamp(t *timestamppb.Timestamp) *time.Time {
	if t != nil {
		tmp := t.AsTime()
		return &tmp
	}
	return nil
}
