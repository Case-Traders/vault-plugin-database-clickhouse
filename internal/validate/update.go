package validate

import "fmt"

// UpdateUser checks UpdateUser request fields.
func UpdateUser(username string, hasPassword, hasExpiration bool) error {
	if username == "" {
		return fmt.Errorf("missing username")
	}
	if !hasPassword && !hasExpiration {
		return fmt.Errorf("no changes requested")
	}
	return nil
}

// CreationStatements checks NewUser has at least one creation statement.
func CreationStatements(count int) error {
	if count <= 0 {
		return fmt.Errorf("empty creation statements")
	}
	return nil
}
