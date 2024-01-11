package resources

import (
	"fmt"
	"strconv"
	"time"
)

type Employee struct {
	EmployeeID int64     `json:"employee_id"`
	Birthday   time.Time `json:"birthday"`
	Hired_on   time.Time `json:"hired_on"`
	Manager    string    `json:"manager"`
	Position   string    `json:"position"`
	Dept       string    `json:"dept"`
	Office     string    `json:"office"`
	Country    string    `json:"country"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Name       string    `json:"name"`
}

func NewEmployee(data map[string]interface{}) (*Employee, error) {
	e := &Employee{}

	e.EmployeeID, _ = strconv.ParseInt(fmt.Sprint(data["EmployeeID"]), 10, 64)

	e.Name = fmt.Sprint(data["Name"])
	e.Email = fmt.Sprint(data["Email"])
	e.Phone = fmt.Sprint(data["Phone"])

	birthdayStr, _ := data["Birthday"].(string)
	if birthdayStr != "" {
		birthday, err := time.Parse("2006-01-02", birthdayStr)
		if err != nil {
			return nil, err
		}
		e.Birthday = birthday
	}

	hiredOnStr, _ := data["Hired_on"].(string)
	if hiredOnStr != "" {
		hiredOn, err := time.Parse("2006-01-02", hiredOnStr)
		if err != nil {
			return nil, err
		}
		e.Hired_on = hiredOn
	}

	e.Dept = fmt.Sprint(data["Dept"])
	e.Manager = fmt.Sprint(data["Manager"])

	return e, nil
}

// fmt.Println("fmt")
