package exporter

import (
	"fmt"
	"time"

	"github.com/c-4u/time-record-service/domain/entity"
)

type SecullumExporter struct {
	EmployeeRegisters []*entity.Employee
}

func NewSecullumExporter(employees []*entity.Employee) (*SecullumExporter, error) {
	return &SecullumExporter{
		EmployeeRegisters: employees,
	}, nil
}

func (e *SecullumExporter) Export() []*Register {
	var registers []*Register
	var i int = 1
	for _, employee := range e.EmployeeRegisters {
		for _, tr := range employee.TimeRecords {
			var r Register

			loc := time.FixedZone("", tr.TzOffset)
			time := tr.Time.In(loc)
			// (10) i = fake id; (8) date; (4) clock; (12) pis
			r = Register(fmt.Sprintf(
				"%010d%02d%02d%d%02d%02d%012s",
				i, time.Day(), time.Month(), time.Year(), time.Hour(), time.Minute(), employee.Pis))
			registers = append(registers, &r)
			i++
		}
	}
	return registers
}
