package bhr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type Field struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type Employee struct {
	ID             string      `json:"id"`
	DisplayName    string      `json:"displayName"`
	FirstName      string      `json:"firstName"`
	LastName       string      `json:"lastName"`
	PreferredName  string      `json:"preferredName"`
	Gender         string      `json:"gender"`
	JobTitle       string      `json:"jobTitle"`
	WorkPhone      string      `json:"workPhone"`
	WorkEmail      string      `json:"workEmail"`
	Department     string      `json:"department"`
	Location       string      `json:"location"`
	Division       string      `json:"division"`
	LinkedIn       interface{} `json:"linkedIn"`
	Supervisor     string      `json:"supervisor"`
	PhotoUploaded  bool        `json:"photoUploaded"`
	PhotoURL       string      `json:"photoUrl"`
	CanUploadPhoto int         `json:"canUploadPhoto"`
	parent         *Employee
	children       []*Employee
}

type Directory struct {
	Fields         []Field    `json:"fields"`
	Employees      []Employee `json:"employees"`
	employeeByName map[string]*Employee
}

type DirectoryCmd struct {
	Client *Client
	Filters
}

type Filters struct {
	Department string
	Title      string
}

type FilterFunc func(e Employee) bool

func (c *Client) GetDirectory(f FilterFunc) (dir *Directory) {
	dir = &Directory{}
	dir.employeeByName = make(map[string]*Employee)

	res, err := c.Request("https://api.bamboohr.com/api/gateway.php/safetyculture/v1/employees/directory")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &dir)
	if err != nil {
		log.Fatal(string(body), err)
	}
	for i, emp := range dir.Employees {
		if !f(emp) {
			continue
		}
		dir.employeeByName[emp.DisplayName] = &dir.Employees[i]
	}
	for i, emp := range dir.Employees {
		if !f(emp) {
			continue
		}
		supervisor := emp.Supervisor
		if supervisor, ok := dir.employeeByName[supervisor]; ok {
			supervisor.children = append(supervisor.children, &dir.Employees[i])
			dir.Employees[i].parent = supervisor
		}
	}
	return
}

func (c *Client) FindEmployeeByName(name string) *Employee {

	nameRegexp := regexp.MustCompile(name)
	nameFilter := func(e Employee) bool {
		return nameRegexp.MatchString(e.DisplayName)
	}

	dir := c.GetDirectory(nameFilter)
	for _, emp := range dir.Employees {
		if !nameFilter(emp) {
			continue
		}
		return &emp
	}
	return nil
}

func Filter(f Filters) FilterFunc {
	departmentRegexp := regexp.MustCompile(f.Department)
	titleRegexp := regexp.MustCompile(f.Title)
	return func(e Employee) bool {
		if f.Department != "" {
			if !departmentRegexp.MatchString(e.Department) {
				return false
			}
			if !titleRegexp.MatchString(e.JobTitle) {
				return false
			}
		}
		return true
	}
}

func (c *DirectoryCmd) Run() error {
	dir := c.Client.GetDirectory(Filter(c.Filters))

	var result strings.Builder
	var currentDepartment string
	for i, emp := range dir.Employees {
		if !Filter(c.Filters)(emp) {
			continue
		}
		if emp.parent == nil {
			currentDepartment = RenderDepartment(&result, 0, currentDepartment, emp.Department)
			Render(&dir.Employees[i], &result, 0, currentDepartment)
		}
	}
	fmt.Println(result.String())
	return nil
}

func Indent(s *strings.Builder, level int) {
	for i := 0; i < level; i++ {
		s.WriteString("  ")
	}
}

func RenderDepartment(s *strings.Builder, level int, currentDepartment string, department string) string {
	if department != currentDepartment {
		s.WriteRune('\n')
		Indent(s, level)
		s.WriteString(fmt.Sprintf("[ %s ]\n", department))
	}
	return department
}

func Render(emp *Employee, s *strings.Builder, level int, currentDepartment string) string {
	currentDepartment = RenderDepartment(s, level, currentDepartment, emp.Department)
	Indent(s, level)
	s.WriteString(fmt.Sprintf("%s (%s)\n", emp.DisplayName, emp.JobTitle))
	for _, emp := range emp.children {
		currentDepartment = Render(emp, s, level+1, currentDepartment)
	}
	return currentDepartment
}
