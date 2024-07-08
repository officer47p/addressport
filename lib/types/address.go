package types

import (
	"fmt"
	"regexp"
)

const (
	minAddressLen = 10
	maxAddressLen = 100
	minNetworkLen = 2
	maxNetworkLen = 20
	minReasonLen  = 5
	maxReaonLen   = 500
	minEmailLen   = 5
	maxEmailLen   = 100
)

func (params UpdateReportParams) Validate() []error {
	errors := []error{}
	if len(params.Reason) < minReasonLen {
		errors = append(errors, fmt.Errorf("reason length should be at least %d characters", minReasonLen))
	}
	if !isEmailValid(params.ContactEmail) {
		errors = append(errors, fmt.Errorf("contact email is not valid"))
	}
	return errors
}

type CreateReportParams struct {
	Address      string `json:"address"`
	Network      string `json:"network"`
	Reason       string `json:"reason"`
	ContactEmail string `json:"contactEmail"`
}

type UpdateReportParams struct {
	ContactEmail string `bson:"contactEmail,omitempty" json:"contactEmail"`
	Reason       string `bson:"reason,omitempty" json:"reason"`
}

func (params CreateReportParams) Validate() []error {
	errors := []error{}
	if len(params.Address) < minAddressLen {
		errors = append(errors, fmt.Errorf("address length should be at least %d characters", minAddressLen))
	}
	if len(params.Network) < minNetworkLen {
		errors = append(errors, fmt.Errorf("network length should be at least %d characters", minNetworkLen))
	}
	if len(params.Reason) < minReasonLen {
		errors = append(errors, fmt.Errorf("reason length should be at least %d characters", minReasonLen))
	}
	if !isEmailValid(params.ContactEmail) {
		errors = append(errors, fmt.Errorf("contact email is not valid"))
	}

	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type Report struct {
	ID           string `bson:"_id,omitempty" json:"id,omitempty"`
	Address      string `bson:"address" json:"address"`
	Network      string `bson:"network" json:"network"`
	Reason       string `bson:"reason" json:"reason"`
	ContactEmail string `bson:"contactEmail" json:"contactEmail"`
}

func NewReportFromParams(params CreateReportParams) (*Report, error) {
	// encPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	// if err != nil {
	// 	return nil, err
	// }
	return &Report{
		Address:      params.Address,
		Network:      params.Network,
		Reason:       params.Reason,
		ContactEmail: params.ContactEmail,
	}, nil
}