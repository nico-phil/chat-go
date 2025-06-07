package main

import (
	"time"

	"github.com/google/uuid"
)

type OPT struct {
	Key       string
	CreatedAt time.Time
}

type RetentionMap map[string]OPT

func NewRetentionMap() RetentionMap {
	m := make(RetentionMap)
	return m
}

func (rm RetentionMap) NewOpt() OPT {
	o := OPT{
		Key:       uuid.NewString(),
		CreatedAt: time.Now(),
	}

	rm[o.Key] = o

	return o
}

func (rm RetentionMap) Verify(optKey string) bool {
	if _, ok := rm[optKey]; !ok {
		return false // opt is not valid
	}

	delete(rm, optKey)
	return true
}
