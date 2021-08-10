package exporter

import (
	"errors"

	"github.com/c-4u/time-record-service/domain/entity"
)

type Register string

type Exporter interface {
	Export() []*Register
}

type ExporterType int

const (
	SECULLUM ExporterType = iota
)

func NewExporter(_type ExporterType, timeRecords []*entity.TimeRecord) (Exporter, error) {
	switch _type {
	case SECULLUM:
		exporter := NewSecullumExporter(timeRecords)
		return exporter, nil
	}
	return nil, errors.New("exporter type not implemented")
}

func (t ExporterType) String() string {
	switch t {
	case SECULLUM:
		return "SECULLUM"
	}
	return ""

}
