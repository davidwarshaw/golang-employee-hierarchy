package internal

import (
	"errors"
	"fmt"
	"sort"
)

// Employee holds all the data for a single employee in the hierarchy,
// including their id, name, salary, and the id of their manager.
type Employee struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Salary    uint   `json:"salary"`
	ManagerId uint   `json:"manager_id"`
	checked   bool
}

// EmployeeHierarchy is the point of entry to manipulate the employees in
// the context of the hierarchy. It holds the single employee at the top
// of the hierarchy.
// EmployeeHierarchy.Head holds the employee at the top of the hierarchy
// EmployeeHierarchy.Managers holds a map from every manager id to a sorted
// array of that manager's reports
type EmployeeHierarchy struct {
	Head     *Employee
	Managers map[uint][]Employee
}

func depthFirstSum(employeeHierarchy *EmployeeHierarchy, employee *Employee) uint {
	// recurse through the employee's reports and total their salaries
	reportsSalaries := uint(0)
	for _, report := range employeeHierarchy.Managers[employee.Id] {
		reportsSalaries += depthFirstSum(employeeHierarchy, &report)
	}
	// The total salary at every point in the hierarchy is that manager's
	// salary, plus the salaries of their reports
	return reportsSalaries + employee.Salary
}

func depthFirstOutput(employeeHierarchy *EmployeeHierarchy, employee *Employee, indent string, isLast bool) {

	// output the ASCII org chart lines and then the employee name
	fmt.Print(indent)
	if isLast {
		fmt.Print("\u2514\u2500")
		indent += "  "
	} else {
		fmt.Print("\u251C\u2500")
		indent += "\u2502 "
	}
	fmt.Println(employee.Name)

	// recurse through the employee's reports to output them
	reports := employeeHierarchy.Managers[employee.Id]
	for reportNum, report := range reports {
		depthFirstOutput(employeeHierarchy, &report, indent, reportNum+1 == len(reports))
	}
}

// NewEmployeeHierarchy creates and returns an EmployeeHierarchy from an array of EmployeeRecord
func NewEmployeeHierarchy(employees []Employee) (*EmployeeHierarchy, error) {
	// declare and initialize the hierarchy
	employeeHierarchy := &EmployeeHierarchy{}
	employeeHierarchy.Managers = map[uint][]Employee{}
	// For every manager we find in the array of employees, add that manager to the map of managers
	for i, employee := range employees {
		// if the employee is their own manager, we have bad data
		if employee.Id == employee.ManagerId {
			return nil, fmt.Errorf("employee id: %d lists themselves as their manager", employee.Id)
		}
		// If the employee's manager has id 0, the employee is the head of the hierarchy
		if employee.ManagerId == 0 {
			// If we already found a head, then we have bad data
			if employeeHierarchy.Head != nil {
				return nil, errors.New("found multiple employees without a manager. Only 1 expected. ")
			}
			employeeHierarchy.Head = &employees[i]
			continue
		}
		// Add the employee as a report of their manager
		employeeHierarchy.Managers[employee.ManagerId] = append(employeeHierarchy.Managers[employee.ManagerId], employee)
	}

	// If there is any manager who is not also a report, then we have bad data
	for managerId := range employeeHierarchy.Managers {
		// Exclude the Head from this check
		if managerId == employeeHierarchy.Head.Id {
			continue
		}
		found := false
		for _, reports := range employeeHierarchy.Managers {
			for _, report := range reports {
				if managerId == report.Id {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("manager id: %d not found as employee", managerId)
		}
	}

	for _, reports := range employeeHierarchy.Managers {
		sort.Slice(reports, func(i, j int) bool {
			return reports[i].Name < reports[j].Name
		})
	}

	return employeeHierarchy, nil
}

// Print outputs the employee hierarchy to stdout formatted with indentation
func (employeeHierarchy *EmployeeHierarchy) Print() {
	depthFirstOutput(employeeHierarchy, employeeHierarchy.Head, "", true)
}

// TotalSalary returns the total salary of the employee hierachy in dollars as a uint
func (employeeHierarchy *EmployeeHierarchy) TotalSalary() uint {
	return depthFirstSum(employeeHierarchy, employeeHierarchy.Head)
}
