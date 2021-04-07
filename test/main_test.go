package test

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/davidwarshaw/golang-employee-hierarchy/internal"
)

func loadFixture(t *testing.T, fixtureFileName string) []internal.Employee {
	employeesJson, err := ioutil.ReadFile("fixtures/" + fixtureFileName)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	var employees []internal.Employee
	if err := json.Unmarshal(employeesJson, &employees); err != nil {
		t.Fatalf("Error: %s", err)
	}
	return employees
}

func TestGood(t *testing.T) {
	employees := loadFixture(t, "good.json")
	employeeHierarchy, err := internal.NewEmployeeHierarchy(employees)
	if err != nil {
		t.Fatalf("Fail: %s", err)
	}

	// we're ok if this just doesn't error
	employeeHierarchy.Print()

	if totalSalary := employeeHierarchy.TotalSalary(); totalSalary != 21 {
		t.Fatalf("Fail: %s", err)
	}
}

func TestTwoCeos(t *testing.T) {
	employees := loadFixture(t, "two-ceos.json")
	expectedError := "found multiple employees without a manager"
	_, err := internal.NewEmployeeHierarchy(employees)
	if err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("Fail: did not find expected error: %s", expectedError)
	}
}

func TestSelfEmployed(t *testing.T) {
	employees := loadFixture(t, "self-employed.json")
	expectedError := "lists themselves as their manager"
	_, err := internal.NewEmployeeHierarchy(employees)
	if err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("Fail: did not find expected error: %s", expectedError)
	}
}

func TestMissingManager(t *testing.T) {
	employees := loadFixture(t, "missing-manager.json")
	expectedError := "not found as employee"
	_, err := internal.NewEmployeeHierarchy(employees)
	if err == nil || !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("Fail: did not find expected error: %s", expectedError)
	}
}
