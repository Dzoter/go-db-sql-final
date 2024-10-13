package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, fmt.Errorf("parcel with number %d not found", number)
		}
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// выполняем запрос к БД для получения строк по client
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel

	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("number", parcel.Number),
		sql.Named("status", status))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}
	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("parcel with number %d is not registered status", number)
	}
	// менять адрес можно только если значение статуса registered
	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
		sql.Named("address", address),
		sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}
	if parcel.Status == ParcelStatusRegistered {
		_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", parcel.Number))
		if err != nil {
			return err
		}
	}

	return nil
}
