//Import Packages
package mysql

import (
	"database/sql"

	"alexedwards.net/snippetbox/pkg/models"

	//"errors"
	"net/http"
	//"encoding/json"
)

//Model Initialization
type EmployeeModel struct {
	DB *sql.DB
}

var employees models.Employee

//Create A Employee Table
func (m *EmployeeModel) CreateTable() (string, error) {
	//MySql Statement
	stmt := ""
	result, err := m.DB.Exec(stmt)
	if err != nil {
		return "c", err
	}
	print(result)
	return "c", nil
}

//Inserting Data Into Employee Model
func (m *EmployeeModel) Insert(empID, empName, Role string) (int, error) {
	//MySql Statement
	stmt := `INSERT INTO employees (emp_id, emp_name, role, created, updated) VALUES(?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := m.DB.Exec(stmt, empID, empName, Role)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	print(id)
	return 0, nil
}

//Updating Data Into Employee Model
func (m *EmployeeModel) Update(empID, empName, Role string) (int, error) {
	//MySql Statement
	stmt := `UPDATE employees SET emp_name = ?, role = ?, updated = UTC_TIMESTAMP() WHERE emp_id = ?`
	_, err := m.DB.Exec(stmt, empName, Role, empID)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (m *EmployeeModel) Show(w http.ResponseWriter, r *http.Request) ([]models.Employee, error) {
	stmt := "SELECT * FROM employees"
	result, err := m.DB.Query(stmt)
	print(123)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	var rows []models.Employee

	for result.Next() {
		var row models.Employee
		if err := result.Scan(&row.ID, &row.EmpID, &row.Role, &row.Created, &row.Updated); err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

//Get Latest Data
func (m *EmployeeModel) Latest() ([]*models.Employee, error) {
	stmt := `SELECT id, emp_id, emp_name, role FROM employees ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := []*models.Employee{}

	for rows.Next() {
		s := &models.Employee{}
		err = rows.Scan(&s.ID, &s.EmpID, &s.EmpName, &s.Role)
		if err != nil {
			return nil, err
		}
		employees = append(employees, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}
