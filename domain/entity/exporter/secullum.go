package exporter

import (
	"fmt"
	"time"

	"github.com/c-4u/time-record-service/domain/entity"
)

type SecullumExporter struct {
	Registers []*entity.TimeRecord
}

func NewSecullumExporter(timeRecords []*entity.TimeRecord) *SecullumExporter {
	return &SecullumExporter{Registers: timeRecords}
}

func (e *SecullumExporter) Export() []*Register {
	var registers []*Register
	var i int = 1
	for _, tr := range e.Registers {
		var r Register

		loc := time.FixedZone("", tr.TzOffset)
		time := tr.Time.In(loc)
		// (10) i = fake id; (8) date; (4) clock; (12) pis
		r = Register(fmt.Sprintf(
			"%010d%02d%02d%d%02d%02d%012s",
			i, time.Day(), time.Month(), time.Year(), time.Hour(), time.Minute(), tr.Employee.Pis))
		registers = append(registers, &r)
		i++
	}
	return registers
}
