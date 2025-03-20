package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Company struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Status      bool    `json:"status"`
}

type CreateCompanyPayload struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func GetCompanies(d *sql.DB) ([]Company, error) {
	var companies []Company

	rows, err := d.Query("SELECT id, name, description, status FROM companies WHERE status = 1;")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var company Company
		err = rows.Scan(&company.ID, &company.Name, &company.Description, &company.Status)
		if err != nil {
			return nil, err
		}

		companies = append(companies, company)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}

func GetCompany(d *sql.DB, ID string) (Company, error) {
	var company Company

	sql := "SELECT id, name, description, status FROM companies WHERE id = ? AND status = 1 LIMIT 1;"
	row := d.QueryRow(sql, ID)
	err := row.Scan(&company.ID, &company.Name, &company.Description, &company.Status)

	if err != nil {
		return Company{}, err
	}

	return company, err
}

func CreateCompany(d *sql.DB, createPayload CreateCompanyPayload) error {
	sql := `
        INSERT INTO companies (id, name, description, status)
        VALUES
            (?, ?, ?, 1);
    `

	_, err := d.Exec(sql, uuid.New().String(), createPayload.Name, createPayload.Description)
	if err != nil {
		return fmt.Errorf("Unable to create new company: %w", err)
	}

	return nil
}

func RemoveCompany(d *sql.DB, companyID string) error {
	sql := `
        DELETE FROM companies WHERE id = ?;
    `

	if _, err := d.Exec(sql, companyID); err != nil {
		return fmt.Errorf("Unable to remove company from id '%s': %w", companyID, err)
	}

	return nil
}
