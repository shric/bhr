package bhr

type PercentCmd struct {
	Client *Client
	EmployeeFilters
}

func (c *PercentCmd) Run() error {
	return nil
}
