package restc

import (
	"time"

	"github.com/pkg/errors"
)

var JsonTimeFormat = "2006-01-02T15:04:05.000+0800"

type JsonTime struct {
	time.Time
}

func (t *JsonTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var err error
	if t.Time, err = time.Parse(`"`+JsonTimeFormat+`"`, string(data)); err != nil {
		return errors.Wrap(err, "time parse")
	}
	return nil
}
