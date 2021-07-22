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

func NewExporter(exporter ExporterType, employees []*entity.Employee) (Exporter, error) {
	switch exporter {
	case SECULLUM:
		exporter, err := NewSecullumExporter(employees)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	}
	return nil, errors.New("exporter type not implemented")
}

func (e ExporterType) String() string {
	switch e {
	case SECULLUM:
		return "SECULLUM"
	}
	return ""

}
