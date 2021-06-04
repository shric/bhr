package bhr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-sixel"
)

func generateFieldsList() string {
	// https://www.bamboohr.com/api/documentation/employees.php
	// Generated with cat list | paste - - - | awk '{ $1="\"" $1 "\", //"; $2 = $2 ":"; print }'
	fields := []string{
		"address1",                // text: The employee's first address line.
		"address2",                // text: The employee's second address line.
		"age",                     // integer: The employee's age. To change age, update dateOfBirth field.
		"bestEmail",               // email: The employee's work email if set, otherwise their home email.
		"birthday",                // text: The employee's month and day of birth. To change birthday, update dateOfBirth field.
		"city",                    // text: The employee's city.
		"country",                 // country: The employee's country.
		"dateOfBirth",             // date: The date the employee was born.
		"department",              // list: The employee's CURRENT department.
		"division",                // list: The employee's CURRENT division.
		"eeo",                     // list: The employee's EEO job category. These are defined by the U.S. Equal Employment Opportunity Commission.
		"employeeNumber",          // text: Employee number (assigned by your company).
		"employmentHistoryStatus", // list: The employee's CURRENT employment status. Options are customized by account. Read-only starting with version 1.1; update using the employmentStatus table.
		"ethnicity",               // list: The employee's ethnicity.
		"exempt",                  // list: The FLSA Overtime Status (Exempt or Non-exempt).
		"overtimeRate",            // currency: The Overtime Rate
		"firstName",               // text: The employee's first name.
		"flsaCode",                // list: Deprecated please use 'exempt'
		"fullName1",               // text: The employee's first and last name. (e.g., John Doe). Read only.
		"fullName2",               // text: The employee's last and first name. (e.g., Doe, John). Read only.
		"fullName3",               // text: The employee's full name and their preferred name. (e.g., Doe, John Quentin (JDog)). Read only.
		"fullName4",               // text: The employee's full name without their preferred name, last name first. (e.g., Doe, John Quentin). Read only.
		"fullName5",               // text: The employee's full name without their preferred name, first name first. (e.g., John Quentin Doe). Read only.
		"displayName",             // text: The employee's name displayed in a format configured by the user. Read only.
		"gender",                  // gender: The employee's gender (Male or Female).
		"hireDate",                // date: The date the employee was hired.
		"originalHireDate",        // date: The date the employee was originally hired. Available starting with version 1.1.
		"homeEmail",               // email: The employee's home email address.
		"homePhone",               // phone: The employee's home phone number.
		"id",                      // integer: The employee ID automatically assigned by BambooHR. Read only.
		"jobTitle",                // list: The CURRENT value of the employee's job title, updating this field will create a new row in position history.
		"lastChanged",             // timestamp: The date and time that the employee record was last changed.
		"lastName",                // text: The employee's last name.
		"location",                // list: The employee's CURRENT location.
		"maritalStatus",           // list: The employee's marital status (Single, Married, or Domestic Partnership).
		"middleName",              // text: The employee's middle name.
		"mobilePhone",             // phone: The employee's mobile phone number.
		"payChangeReason",         // list: The reason for the employee's last pay rate change.
		"payGroup",                // list: The custom pay group that the employee belongs to.
		"payGroupId",              // integer: The ID value corresponding to the pay group that an employee belongs to.
		"payRate",                 // currency: The employee's CURRENT pay rate (e.g., $8.25).
		"payRateEffectiveDate",    // date: The day the most recent change was made.
		"payType",                 // pay_type: The employee's CURRENT pay type. ie: "hourly","salary","commission","exception hourly","monthly","weekly","piece rate","contract","daily","pro rata".
		"payPer",                  // paid_per: The employee's CURRENT pay per. ie: "Hour", "Day", "Week", "Month", "Quarter", "Year".
		"paidPer",                 // paid_per: The employee's CURRENT pay per. ie: "Hour", "Day", "Week", "Month", "Quarter", "Year".
		"paySchedule",             // list: The employee's CURRENT pay schedule.
		"payScheduleId",           // integer: The ID value corresponding to the pay schedule that an employee belongs to.
		"payFrequency",            // list: The employee's CURRENT pay frequency. ie: "Weekly", "Every other week", "Twice a month", "Monthly", "Quarterly", "Twice a year", or "Yearly"
		"includeInPayroll",        // bool: Should employee be included in payroll (Yes or No)
		"timeTrackingEnabled",     // bool: Should time tracking be enabled for the employee (Yes or No)
		"preferredName",           // text: The employee's preferred name.
		"ssn",                     // ssn: The employee's Social Security number.
		"sin",                     // sin: The employee's Canadian Social Insurance Number.
		"state",                   // state: The employee's state/province.
		"stateCode",               // text: The 2 character abbreviation for the employee's state (US only). Read only.
		"status",                  // status: The employee's employee status (Active or Inactive).
		"supervisor",              // employee: The employee’s CURRENT supervisor. Read only.
		"supervisorId",            // integer: The 'employeeNumber' of the employee's CURRENT supervisor. Read only.
		"supervisorEmail",         // integer: The email address of the employee's CURRENT supervisor. Read only.
		"supervisorEId",           // integer: The ID of the employee's CURRENT supervisor. Read only.
		"terminationDate",         // date: The date the employee was terminated. Read-only starting with version 1.1; update using the employmentStatus table.
		"workEmail",               // email: The employee's work email address.
		"workPhone",               // phone: The employee's work phone number, without extension.
		"workPhonePlusExtension",  // text: The employee's work phone and extension. Read only.
		"workPhoneExtension",      // text: The employee's work phone extension (if any).
		"zipcode",                 // text: The employee's ZIP code.
		"isPhotoUploaded",         // bool: Whether a photo has been uploaded for the employee. Read only.
		"acaStatus",               // text: The employee's ACA (Affordable Care Act) Full-Time status. Options are yes, no, non-assessment
		"acaStatusCategory",       // text: The employee's ACA (Affordable Care Act) status.
		"standardHoursPerWeek",    // integer: The number of hours the employee works in a standard week.
		"bonusDate",               // date: The date of the last bonus.
		"bonusAmount",             // currency: The amount of the most recent bonus.
		"bonusReason",             // list: The reason for the most recent bonus.
		"bonusComment",            // text: Comment about the most recent bonus.
		"commissionDate",          // date: The date of the last commission.
		"commisionDate",           // date: This field name contains a typo, and exists for backwards compatibility.
		"commissionAmount",        // currency: The amount of the most recent commission.
		"commissionComment",       // text: Comment about the most recent commission.
		"employmentStatus",        // status: DEPRECATED. Please use "status" instead. The employee's employee status (Active or Inactive).
		"nickname",                // text: DEPRECATED. Please use "preferredName" instead. The employee's preferred name.
		"payPeriod",               // pay_period: DEPRECATED. Please use paySchedule or payFrequency instead. The employee's CURRENT pay period. ie: "Daily", "Weekly", "Every other week", "Twice a month", "Monthly", "Quarterly", "Twice a year", "Yearly".
		"photoUploaded",           // bool: DEPRECATED. Please use "isPhotoUploaded" instead. The employee has uploaded a photo. Read only.
		"nin",                     // nin: The employee's NIN number
		"nationalId",              // national_id: The employee's National ID number
		"nationality",             // list: The employee's nationality
		"employeeStatusDate",      // date: The effective date on the employment status table
		"employeeTaxType",         // text: The employee's tax type on the employment status table
		"allergies",               // text: The employee's allergies
		"dietaryRestrictions",     // text: The employee's dietary restrictions
		"birthplace",              // text: The employee's birthplace
		"secondaryLanguage",       // text: The employee's secondary language
		"probationEndDate",        // date: The employee's probation end date
		"contractEndDate",         // date: The employee's contract end date
		"citizenship",             // country: The employee's citizenship
		"shirtSize",               // list: The employee's shirt size
		"tShirtSize",              // list: The employee's t-shirt size
		"jacketSize",              // list: The employee's jacket size
		"noticePeriod",            // list: The employee's notice period
		"team",                    // list: The employee's team
	}
	return strings.Join(fields[:], ",")
}

type IndividualEmployee struct {
	ID                      string      `json:"id"`
	Address1                string      `json:"address1"`
	Address2                string      `json:"address2"`
	Age                     string      `json:"age"`
	BestEmail               string      `json:"bestEmail"`
	Birthday                string      `json:"birthday"`
	City                    string      `json:"city"`
	Country                 string      `json:"country"`
	DateOfBirth             string      `json:"dateOfBirth"`
	Department              string      `json:"department"`
	Division                string      `json:"division"`
	EmployeeNumber          string      `json:"employeeNumber"`
	EmploymentHistoryStatus string      `json:"employmentHistoryStatus"`
	FirstName               string      `json:"firstName"`
	FullName1               string      `json:"fullName1"`
	FullName2               string      `json:"fullName2"`
	FullName3               string      `json:"fullName3"`
	FullName4               string      `json:"fullName4"`
	FullName5               string      `json:"fullName5"`
	DisplayName             string      `json:"displayName"`
	Gender                  interface{} `json:"gender"`
	HireDate                string      `json:"hireDate"`
	OriginalHireDate        string      `json:"originalHireDate"`
	HomeEmail               string      `json:"homeEmail"`
	HomePhone               interface{} `json:"homePhone"`
	JobTitle                string      `json:"jobTitle"`
	LastChanged             time.Time   `json:"lastChanged"`
	LastName                string      `json:"lastName"`
	Location                string      `json:"location"`
	MiddleName              interface{} `json:"middleName"`
	MobilePhone             string      `json:"mobilePhone"`
	PayChangeReason         string      `json:"payChangeReason"`
	PayRate                 string      `json:"payRate"`
	PayRateEffectiveDate    string      `json:"payRateEffectiveDate"`
	PayType                 string      `json:"payType"`
	PayPer                  string      `json:"payPer"`
	PaidPer                 string      `json:"paidPer"`
	PaySchedule             string      `json:"paySchedule"`
	PayScheduleID           string      `json:"payScheduleId"`
	PayFrequency            string      `json:"payFrequency"`
	PreferredName           string      `json:"preferredName"`
	State                   string      `json:"state"`
	StateCode               string      `json:"stateCode"`
	Supervisor              string      `json:"supervisor"`
	SupervisorID            string      `json:"supervisorId"`
	SupervisorEmail         string      `json:"supervisorEmail"`
	SupervisorEID           string      `json:"supervisorEId"`
	TerminationDate         string      `json:"terminationDate"`
	WorkEmail               string      `json:"workEmail"`
	WorkPhone               string      `json:"workPhone"`
	Zipcode                 string      `json:"zipcode"`
	IsPhotoUploaded         string      `json:"isPhotoUploaded"`
	AcaStatusCategory       interface{} `json:"acaStatusCategory"`
	Nickname                string      `json:"nickname"`
	PayPeriod               string      `json:"payPeriod"`
	PhotoUploaded           bool        `json:"photoUploaded"`
	PhotoURL                string      `json:"photoUrl"`
	Nationality             interface{} `json:"nationality"`
	EmployeeStatusDate      string      `json:"employeeStatusDate"`
}

type EmployeeCmd struct {
	Client *Client
	Image  bool
	EmployeeFilters
}

type EmployeeFilters struct {
	Name string
	ID   int
}

func (c *Client) GetEmployee(id int) *IndividualEmployee {
	dir := &IndividualEmployee{}
	res, err := c.Request(fmt.Sprintf("https://api.bamboohr.com/api/gateway.php/safetyculture/v1/employees/%d?fields=%s",
		id, generateFieldsList()))
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
	return dir
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func IRender(e *IndividualEmployee, s *strings.Builder) {
	pad := 15
	s.WriteString(fmt.Sprintf("%-*s%s\n", pad, "ID: ", e.ID))

	if e.DisplayName != "" {
		s.WriteString(fmt.Sprintf("%-*s%s", pad, "Name: ", e.DisplayName))
		if e.JobTitle != "" {
			s.WriteString(fmt.Sprintf(" (%s)\n", e.JobTitle))
		}
	}

	if e.WorkEmail != "" {
		s.WriteString(fmt.Sprintf("%-*s%s\n", pad, "Email: ", e.WorkEmail))
	}

	if e.WorkPhone != "" {
		s.WriteString(fmt.Sprintf("%-*s%s\n", pad, "Phone: ", e.WorkPhone))
	}
	/*
		if e.Address1 != "" {
			s.WriteString(fmt.Sprintf("%-*s%s", pad, "Address: ", e.Address1))
			if e.Address2 != "" {
				s.WriteString(fmt.Sprintf(", %s", e.Address2))
			}
			if e.City != "" {
				s.WriteString(fmt.Sprintf(", %s", e.City))
			}
			if e.State != "" {
				s.WriteString(fmt.Sprintf(", %s", e.State))
			}
			if e.Zipcode != "" {
				s.WriteString(fmt.Sprintf(", %s", e.Zipcode))
			}
			if e.Country != "" {
				s.WriteString(fmt.Sprintf(", %s", e.Country))
			}
			s.WriteRune('\n')
		}
	*/

	if e.Department != "" {
		s.WriteString(fmt.Sprintf("%-*s%s\n", pad, "Department: ", e.Department))
	}

	if e.Supervisor != "" {
		s.WriteString(fmt.Sprintf("%-*s%s\n", pad, "Supervisor: ", e.Supervisor))
	}

	if e.HireDate != "" {
		s.WriteString(fmt.Sprintf("%-*s%s\n", pad, "Hire date: ", e.HireDate))
	}

	if e.Location != "" {
		s.WriteString(fmt.Sprintf("%-*s%s\n", pad, "Location: ", e.Location))
	}
}

func (c *EmployeeCmd) Run() error {
	var employee *IndividualEmployee
	if c.ID != -1 {
		employee = c.Client.GetEmployee(c.ID)
	} else if c.Name != "" {
		e := c.Client.FindEmployeeByName(c.Name)
		id, err := strconv.Atoi(e.ID)
		if err != nil {
			return err
		}
		employee = c.Client.GetEmployee(id)
	} else {
		employee = c.Client.GetEmployee(0)

	}
	var result strings.Builder

	IRender(employee, &result)
	if c.Image && employee.PhotoUploaded && employee.PhotoURL != "" {
		err := c.Client.ImageShow(employee.PhotoURL, &result)
		if err != nil {
			return err
		}
	}
	fmt.Println()
	fmt.Println(result.String())
	return nil
}

func (c *Client) ImageShow(url string, s *strings.Builder) error {
	res, err := c.Request(strings.Replace(url, "-1.jpg", "-2.jpg", -1))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	if err != nil {
		return err
	}
	buf := bufio.NewWriter(os.Stdout)
	defer buf.Flush()

	enc := sixel.NewEncoder(buf)
	enc.Dither = true
	return enc.Encode(img)
}
