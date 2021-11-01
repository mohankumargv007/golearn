//Import Packages
package mysql

import (
	"database/sql"
	//"alexedwards.net/snippetbox/pkg/models"
	//"errors"
)

//Model Initialization
type EmployeeModel struct {
	DB *sql.DB
}

//Create A Employee Table 
func (m *EmployeeModel) CreateTable() (string, error) {
	//MySql Statement
	stmt := ""
	//w.Write([]byte(stmt))
	print(stmt)
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

func (m *EmployeeModel) showAllEmpList() ([]byte, error) {
	stmt := "SELECT * FROM employees"
	result, err := m.DB.Exec(stmt)
	if err != nil {
		return [] 
	}	
	return result, ""
}