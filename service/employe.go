package service

import (
	"encoding/json"
	"fmt"
	"importxcel/db"
	"importxcel/resources"
	"strconv"
	"time"

	// "github.com/360EntSecGroup-Skylar/excelize"

	// "demo/utils"

	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gorilla/mux"
)

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	employeeID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	fmt.Println(employeeID)
	// getting from redis
	// valueStr, err := db.RdsClientLocal.HGet("employees", fmt.Sprintf("%v", employeeID)).Result()
	// if err != nil {	if err != nil {
	// 	return "", err
	// }
	res, err := db.GetEmployee(employeeID)
	if err != nil {
		http.Error(w, "employee not found", http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Error converting employee data to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getIntValue(val string) int {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	return intVal
}

const MaxMemory int64 = 32 << 20

func Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxMemory)
	fmt.Println("input", r)
	// // Set the response header to indicate that the response contains a file
	// w.Header().Set("Content-Disposition", "attachment; filename=data.csv")

	// // Set the response header to indicate the content type of the response
	// w.Header().Set("Content-Type", "text/csv")

	// Parse the request body as multipart/form-data
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Println("Error parsing request:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	xlFile, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	sheetName := xlFile.GetSheetName(0)
	rows, err := xlFile.GetRows(sheetName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	employees := []*resources.Employee{}
	row := 1
	for {
		if len(rows[row]) < 11 {
			fmt.Println("One or more required fields are missing", rows[row])
			break
		}
		id, err := strconv.ParseInt(rows[row][0], 10, 64)
		if err != nil {
			fmt.Println("Error parsing EmployeeID:", err)
			return
		}
		// birthdayStr :=.Format("2006-01-02")
		// hiredStr := rows[row][2].Format("2006-01-02")
		t, err := time.Parse("01-02-06", rows[row][1])
		if err != nil {
			fmt.Println(":", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h, err := time.Parse("01-02-06", rows[row][2])
		if err != nil {
			fmt.Println(":", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		e := &resources.Employee{
			EmployeeID: id,
			Birthday:   t,
			Hired_on:   h,
			Manager:    rows[row][3],
			Position:   rows[row][4],
			Dept:       rows[row][5],
			Office:     rows[row][6],
			Country:    rows[row][7],
			Email:      rows[row][8],
			Phone:      rows[row][9],
			Name:       rows[row][10]}

		// Create a new employee
		newEmployee := &resources.Employee{EmployeeID: e.EmployeeID, Birthday: e.Birthday, Hired_on: e.Hired_on, Manager: e.Manager, Position: e.Position, Dept: e.Dept, Office: e.Office, Country: e.Country, Email: e.Email, Phone: e.Phone, Name: e.Name}
		// Convert the employee to a slice
		// empJSON, err := json.Marshal(newEmployee)
		// err = db.RdsClientLocal.HSet(r.Context(), strconv.Itoa(int(newEmployee.EmployeeID)), string(empJSON)).Err()
		// if err != nil {
		// 	fmt.Println("Error adding employee to the cache:", err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		employees = append(employees, newEmployee)
		row++
		if row == len(rows) {
			break
		}
	}
	// Add the employee to the database
	err = db.AddEmployees(employees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully uploaded %d files\n", len(r.MultipartForm.File))
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	employeeID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var existingEmployee resources.Employee
	existingEmployee, err = db.GetEmployee(employeeID)
	if err != nil {
		http.Error(w, "employee not found", http.StatusInternalServerError)
		return
	}

	var updatedEmployee resources.Employee
	if err := json.NewDecoder(r.Body).Decode(&updatedEmployee); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	existingEmployee.Birthday = updatedEmployee.Birthday
	existingEmployee.Hired_on = updatedEmployee.Hired_on
	existingEmployee.Manager = updatedEmployee.Manager
	existingEmployee.Position = updatedEmployee.Position
	existingEmployee.Office = updatedEmployee.Office
	existingEmployee.Country = updatedEmployee.Country
	existingEmployee.Email = updatedEmployee.Email
	existingEmployee.Phone = updatedEmployee.Phone
	existingEmployee.Name = updatedEmployee.Name

	err = db.SaveEmployee(existingEmployee)
	if err != nil {

		// return err
	}
	// Update in Redis
	// empJSON, err := json.Marshal(existingEmployee)
	// if err != nil {
	// 	log.Println("Error marshaling JSON:", err)
	// } else {
	// 	rdb.HSet("employees", strconv.Itoa(existingEmployee.EmployeeID), string(empJSON)).Err()
	// }

	fmt.Fprintf(w, "Successfully Updated")

}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	employeeID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}
	err = db.DeleteEmployee(employeeID)
	if err != nil {

	}
	// Delete from Redis
	// db.RdsClientLocal.HDel(r.Context(), strconv.Itoa(employeeID)).Err()
	fmt.Fprintf(w, "Successfully Deleted")

	// respondJSON(w, http.StatusNoContent, nil)
}
