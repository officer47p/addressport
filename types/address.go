package types

import (
	"fmt"
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

// type UpdateAddressParams struct {
// 	FirstName string `bson:"firstName,omitempty" json:"firstName,omitempty"`
// 	LastName  string `bson:"lastName,omitempty" json:"lastName,omitempty"`
// }

// func (params UpdateAddressParams) Validate() []error {
// 	errors := []error{}
// 	if params.FirstName != "" && len(params.FirstName) < minNameLen {
// 		errors = append(errors, fmt.Errorf("firstName length should be at least %d characters", minNameLen))
// 	}
// 	if params.LastName != "" && len(params.LastName) < minNameLen {
// 		errors = append(errors, fmt.Errorf("lastName length should be at least %d characters", minNameLen))
// 	}
// 	return errors
// }

type CreateAddressParams struct {
	Address      string `json:"address"`
	Network      string `json:"network"`
	Reason       string `json:"reason"`
	ContactEmail string `json:"contactEmail"`
}

func (params CreateAddressParams) Validate() []error {
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
	// if !isEmailValid(params.Email) {
	// 	errors = append(errors, fmt.Errorf("email is not valid"))
	// }

	return errors
}

// func isEmailValid(e string) bool {
// 	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
// 	return emailRegex.MatchString(e)
// }

type Address struct {
	ID           string `bson:"_id,omitempty" json:"id,omitempty"`
	Address      string `bson:"address" json:"address"`
	Network      string `bson:"network" json:"network"`
	Reason       string `bson:"reason" json:"reason"`
	ContactEmail string `bson:"contactEmail" json:"contactEmail"`
}

func NewAddressFromParams(params CreateAddressParams) (*Address, error) {
	// encPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	// if err != nil {
	// 	return nil, err
	// }
	return &Address{
		Address:      params.Address,
		Network:      params.Network,
		Reason:       params.Reason,
		ContactEmail: params.ContactEmail,
	}, nil
}
