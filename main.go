package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davidwarshaw/golang-employee-hierarchy/internal"
	"github.com/dustin/go-humanize"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error: %v\n", "employee records filename expected")
		os.Exit(1)
	}
	filename := os.Args[1]

	employeesJson, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	var employees []internal.Employee
	if err := json.Unmarshal(employeesJson, &employees); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	employeeHierarchy, err := internal.NewEmployeeHierarchy(employees)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Display the hierarchy
	fmt.Println("Employee hierarchy:")
	employeeHierarchy.Print()

	// Calculate and display the total salaries of the hierarchy
	totalSalaries := employeeHierarchy.TotalSalary()
	fmt.Printf("Total salaries: $%s.\n", humanize.Comma(int64(totalSalaries)))
}
