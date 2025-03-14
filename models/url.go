package models

type URLUsageFiltered struct {
	Url      string
	Content  string
	Filepath string
	// @TODO review the fields below
	Request     string
	Description string
}
