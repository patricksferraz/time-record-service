package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/patricksferraz/time-record-service/domain/entity"
	_ "gorm.io/driver/sqlite"
)

type Postgres struct {
	Db *gorm.DB
}

func NewPostgres(dsnType, dsn string) (*Postgres, error) {
	postgres := &Postgres{}

	err := postgres.connect(dsnType, dsn)
	if err != nil {
		return nil, err
	}

	return postgres, nil
}

func (p *Postgres) connect(dsnType, dsn string) error {
	db, err := gorm.Open(dsnType, dsn)
	if err != nil {
		return fmt.Errorf("Error connecting to database: %v", err)
	}

	p.Db = db

	return nil
}

func (p *Postgres) Debug(enable bool) {
	p.Db.LogMode(enable)
}

func (p *Postgres) Migrate() {
	p.Db.AutoMigrate(&entity.Employee{}, &entity.TimeRecord{}, &entity.Company{})
}
