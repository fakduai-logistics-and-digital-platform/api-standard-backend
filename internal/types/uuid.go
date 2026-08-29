package types

import (
	"database/sql/driver"
	jsonv2 "encoding/json/v2"
	"errors"
	"uuid"

	"go.mongodb.org/mongo-driver/bson"
)

type UUID [16]byte

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

func (u UUID) Value() (driver.Value, error) {
	return uuid.UUID(u).String(), nil
}

func (u *UUID) Scan(src any) error {
	switch v := src.(type) {
	case []byte:
		if len(v) == 16 {
			*u = UUID(v)
			return nil
		}
		return errors.New("invalid UUID bytes length")
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*u = UUID(id)
		return nil
	default:
		return errors.New("invalid UUID source")
	}
}

func (u UUID) MarshalJSON() ([]byte, error) {
	return jsonv2.Marshal(u.String())
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	var str string
	err := jsonv2.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(str)
	if err != nil {
		return err
	}

	*u = UUID(id)
	return nil
}

func (u UUID) MarshalBSON() ([]byte, error) {
	return bson.Marshal(u.String())
}

func (u *UUID) UnmarshalBSON(data []byte) error {
	var str string
	err := bson.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(str)
	if err != nil {
		return err
	}

	*u = UUID(id)
	return nil
}
