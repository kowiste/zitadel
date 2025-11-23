package domain

import "time"

type Organization struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	PrimaryDomain string    `json:"primaryDomain"`
	State         string    `json:"state"`
	CreationDate  time.Time `json:"creationDate"`
	ChangedDate   time.Time `json:"changedDate"`
}

type User struct {
	ID               string    `json:"id"`
	UserName         string    `json:"userName"`
	State            string    `json:"state"`
	PreferredLoginName string  `json:"preferredLoginName"`
	Email            string    `json:"email,omitempty"`
	FirstName        string    `json:"firstName,omitempty"`
	LastName         string    `json:"lastName,omitempty"`
	CreationDate     time.Time `json:"creationDate"`
	ChangedDate      time.Time `json:"changedDate"`
}

type OrganizationList struct {
	Organizations []Organization `json:"result"`
	TotalResult   string         `json:"totalResult,omitempty"`
}

type UserList struct {
	Users       []User `json:"result"`
	TotalResult string `json:"totalResult,omitempty"`
}
