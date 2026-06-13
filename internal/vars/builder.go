package vars

import (
	"maps"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
)

// Operation identifies which Vault dbplugin call built a template map.
type Operation int

const (
	// OpNewUser is the NewUser dbplugin path.
	OpNewUser Operation = iota
	// OpUpdatePassword is the UpdateUser password path.
	OpUpdatePassword
	// OpUpdateExpiration is the UpdateUser expiration path.
	OpUpdateExpiration
	// OpDeleteUser is the DeleteUser dbplugin path.
	OpDeleteUser
)

// TemplateVars holds {{placeholder}} substitutions for one statement batch.
type TemplateVars struct {
	values map[string]string
}

// New copies base into a TemplateVars value.
func New(base map[string]string) TemplateVars {
	copied := make(map[string]string, len(base))
	maps.Copy(copied, base)
	return TemplateVars{values: copied}
}

func (v TemplateVars) Map() map[string]string {
	out := make(map[string]string, len(v.values))
	maps.Copy(out, v.values)
	return out
}

func (v TemplateVars) Keys() []string {
	return F.Pipe1(v.values, mapKeys)
}

func mapKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func (v TemplateVars) Has(key string) bool {
	_, ok := v.values[key]
	return ok
}

var requiredByOp = map[Operation][]string{
	OpNewUser:          {"name", "username", "password", "expiration", "cluster"},
	OpUpdatePassword:   {"name", "username", "password", "cluster"},
	OpUpdateExpiration: {"name", "username", "expiration", "cluster"},
	OpDeleteUser:       {"name", "username"},
}

// RequiredKeys lists template placeholders used by default statements for op.
func RequiredKeys(op Operation) []string {
	return requiredByOp[op]
}

// HasRequiredKeys reports whether v contains every key RequiredKeys lists for op.
func HasRequiredKeys(op Operation, v TemplateVars) bool {
	return F.Pipe1(
		RequiredKeys(op),
		A.Reduce(func(acc bool, k string) bool { return acc && v.Has(k) }, true),
	)
}

func ForNewUser(username, password, expiration, cluster string) TemplateVars {
	return New(map[string]string{
		"name":       username,
		"username":   username,
		"password":   password,
		"expiration": expiration,
		"cluster":    cluster,
	})
}

func ForUpdatePassword(username, password, cluster string) TemplateVars {
	return New(map[string]string{
		"name":     username,
		"username": username,
		"password": password,
		"cluster":  cluster,
	})
}

func ForUpdateExpiration(username, expiration, cluster string) TemplateVars {
	return New(map[string]string{
		"name":       username,
		"username":   username,
		"expiration": expiration,
		"cluster":    cluster,
	})
}

func ForDeleteUser(username string) TemplateVars {
	return New(map[string]string{
		"name":     username,
		"username": username,
	})
}
