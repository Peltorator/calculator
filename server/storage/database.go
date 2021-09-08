package storage

import (
	"database/sql"
	"fmt"
)

type DataBaseHistoryStorage struct {
	conn *sql.DB
}

func New(conn *sql.DB) *DataBaseHistoryStorage {
	return &DataBaseHistoryStorage{conn: conn}
}

const queryStoreCalculation = `
	INSERT INTO history(
		expression,
		result 
	) VALUES ($1, $2)
`

func (p *DataBaseHistoryStorage) StoreCalculation(c Calculation) error {
	fmt.Println("kek1 ", c.Expression, " ", c.Result)
	row := p.conn.QueryRow(queryStoreCalculation, 1, 2)
	fmt.Println("kek2")
	err := row.Scan()
	if err != nil {
		return err
	}
	return nil
}

const queryGetHistory = `
	SELECT (expression, result) FROM history
`

func (p *DataBaseHistoryStorage) GetHistory() ([]Calculation, error) {
	c := []Calculation{}

	rows, err := p.conn.Query(queryGetHistory)
	if err != nil {
		return c, err
	}

	for rows.Next() {
		r := new(Calculation)
		err = rows.Scan(&r.Expression, &r.Result)
		if err != nil {
			return c, err
		}
		c = append(c, *r)
	}
	return c, nil
}


