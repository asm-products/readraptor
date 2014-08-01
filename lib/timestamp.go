package readraptor

import (
	"database/sql/driver"
	"strconv"
	"time"
)

type Timestamp struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *Timestamp) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt Timestamp) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte(`""`), nil
	}
	return []byte(strconv.FormatInt(t.Time.Local().Unix(), 10)), nil
}
