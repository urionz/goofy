package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	Id        uint           `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	CreatedAt FmtTime        `json:"created_at"`
	UpdatedAt FmtTime        `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type JSON []byte

func (j JSON) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("invalid Scan Source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("null point exception")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

func (j JSON) Equals(j1 JSON) bool {
	return bytes.Equal(j, j1)
}

type Strings []string

func (val Strings) Value() (driver.Value, error) {
	data, err := json.Marshal(val)
	return string(data), err
}

func (val *Strings) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &val)
}

type FmtTime struct {
	time.Time
}

func (t FmtTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (t FmtTime) Normalize(layout ...string) string {
	if len(layout) == 0 {
		layout = append(layout, "2006-01-02 15:04:05")
	}
	return fmt.Sprintf("\"%s\"", t.Format(layout[0]))
}

func (t FmtTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *FmtTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = FmtTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
