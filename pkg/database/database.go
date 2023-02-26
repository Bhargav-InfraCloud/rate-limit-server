package database

import "github.com/Bhargav-InfraCloud/rate-limit-server/pkg/service/errors"

type DB struct {
	table map[string]uint
}

type RateLimiterDB interface {
	Add(id string, count uint, resetCount bool) *errors.Error
}

func NewDB() RateLimiterDB {
	return &DB{
		table: make(map[string]uint),
	}
}

func (d *DB) Add(id string, count uint, resetCount bool) *errors.Error {
	if count == 0 {
		return errors.RateLimitedError
	}

	if resetCount {
		d.table[id] = count
	}

	if _, ok := d.table[id]; ok {
		if d.table[id] == 0 {
			return errors.RateLimitedError
		}

		d.table[id]--

		return nil
	}

	d.table[id] = count - 1

	return nil
}
