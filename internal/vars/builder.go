package vars

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
	for k, v := range base {
		copied[k] = v
	}
	return TemplateVars{values: copied}
}

func (v TemplateVars) Map() map[string]string {
	out := make(map[string]string, len(v.values))
	for k, val := range v.values {
		out[k] = val
	}
	return out
}

func (v TemplateVars) Keys() []string {
	out := make([]string, 0, len(v.values))
	for k := range v.values {
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
	for _, k := range RequiredKeys(op) {
		if !v.Has(k) {
			return false
		}
	}
	return true
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
