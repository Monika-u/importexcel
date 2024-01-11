package db

import (
	"fmt"
	"importxcel/resources"
	"log"
)

// 	"golang.org/x/crypto/bcrypt"

// func (d *taskRepo) ImportTasks(ctx context.Context, tasks []*models.Task) error {
// 	span, _ := apm.StartSpan(ctx, "ImportTasks", "db")
// 	defer span.End()
// 	for _, v := range tasks {
// 		db := getTransactionFromContext(ctx, d.db)
// 		res := db.WithContext(ctx).
// 			Model(&models.Task{}).
// 			Create(v)
// 		if res.Error != nil {
// 			return res.Error
// 		}
// 	}
// 	return nil
// }

func GetEmployee(employeeID int) (resources.Employee, error) {
	// var employee resources.Employee

	db := DbClient
	var result resources.Employee
	query := "select * from employee where employeeID=?"

	err := db.QueryRow(query, employeeID).Scan(&result.EmployeeID, &result.Birthday, &result.Hired_on, &result.Manager, &result.Position, &result.Dept, result.Office,
		&result.Country, &result.Email, &result.Phone, &result.Name)
	if err != nil {
		return result, err
	}
	fmt.Println(result, "result")
	return result, nil
}

func AddEmployees(employees []*resources.Employee) error {
	db := DbClient
	query := "insert into employee (EmployeeID, birthday, hired_on, manager, position, dept, office, country, email, phone, name) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	// var rows *sql.Rows
	// var err error
	for _, emp := range employees {
		values := []interface{}{
			emp.EmployeeID,
			emp.Birthday,
			emp.Hired_on,
			emp.Manager,
			emp.Position,
			emp.Dept,
			emp.Office,
			emp.Country,
			emp.Email,
			emp.Phone,
			emp.Name,
		}

		rows, err := db.Query(query, values...)
		// handle error

		if err != nil {
			// Handle the error here.
			log.Printf("Error executing query: %v", err)
		} else {
			for rows.Next() {
				// Retrieve the values of the current row here.
			}
			rows.Close()
		}
	}

	// for _, employee := range employees {
	// 	err := db.QueryRow(query, employee.EmployeeID, employee.Birthday, employee.Hired_on, employee.Manager, employee.Position, employee.Dept, employee.Office,
	// 		employee.Country, employee.Email, employee.Phone, employee.Name)
	// 	if err != nil {
	// 		// Handle the error here.
	// 		// if err == sql.ErrNoRows {
	// 		// 	// No rows were returned by the query.
	// 		// 	log.Println("No rows were returned bxy the query.")
	// 		// } else {
	// 		// There was an error executing the query.
	// 		fmt.Println("Error executing query: %v", err)

	// 	}
	// }

	return nil
}

func SaveEmployee(employee resources.Employee) error {
	stmt, err := DbClient.Prepare("UPDATE employees SET name=?, email=?, phone=?, position=?, office=?, country=?, hired_on=?, manager=?, birthday=? WHERE employeeID=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(employee.Name, employee.Email, employee.Phone, employee.Position, employee.Office, employee.Country, employee.Hired_on, employee.Manager, employee.Birthday, employee.EmployeeID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEmployee(employeeID int) error {
	stmt, err := DbClient.Prepare("DELETE FROM employees WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(employeeID)
	if err != nil {
		return err
	}

	return nil
}
