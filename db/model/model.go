package model

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-module/carbon"
	"github.com/urionz/goofy/contracts"
)

type BaseModel struct {
	Id        uint      `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	CreatedAt FmtTime   `json:"created_at,omitempty"`
	UpdatedAt FmtTime   `json:"updated_at,omitempty"`
	DeletedAt DeletedAt `gorm:"index" json:"-"`
}

var _ contracts.DBConnection = (*BaseModel)(nil)

func (*BaseModel) Connection() string {
	return ""
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

type FmtTime sql.NullTime

func (t FmtTime) MarshalJSON() ([]byte, error) {
	if t.Valid {
		var stamp = fmt.Sprintf("\"%s\"", t.Time.Format("2006-01-02 15:04:05"))
		return []byte(stamp), nil
	}
	return json.Marshal(nil)
}

func (t *FmtTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		t.Valid = false
		return nil
	}
	t.Time = carbon.ParseByFormat("2006-01-02 15:04:05", "2006-01-02 15:04:05").Time
	parsed := carbon.ParseByFormat(string(b), "2006-01-02 15:04:05")
	t.Valid = true
	if parsed.Error == nil {
		t.Time = parsed.Time
	}
	return nil
}

func (t FmtTime) Normalize(layout ...string) string {
	if len(layout) == 0 {
		layout = append(layout, "2006-01-02 15:04:05")
	}
	return fmt.Sprintf("\"%s\"", t.Time.Format(layout[0]))
}

func (t FmtTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *FmtTime) Scan(v interface{}) error {
	return (*sql.NullTime)(t).Scan(v)
	// value, ok := v.(time.Time)
	// if ok {
	// 	*t = FmtTime{Time: value}
	// 	return nil
	// }
	// return fmt.Errorf("can not convert %v to timestamp", v)
}
