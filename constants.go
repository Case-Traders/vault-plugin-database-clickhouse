package clickhouse

// Plugin and Vault database config keys.
const (
	// TypeName is the name passed to vault plugin register.
	TypeName = "clickhouse"
	// ConfigKeyCluster is the optional Vault database config key for ON CLUSTER DDL.
	ConfigKeyCluster = "cluster"
	// ConfigKeyUsernameTemplate is the optional Vault username template key.
	ConfigKeyUsernameTemplate = "username_template"
	// DefaultChangePasswordStatement is used when UpdateUser password statements are empty.
	DefaultChangePasswordStatement = "ALTER USER '{{username}}' IDENTIFIED WITH plaintext_password BY '{{password}}' ON CLUSTER '{{cluster}}';"
	// DefaultChangeExpirationStatement is used when UpdateUser expiration statements are empty.
	DefaultChangeExpirationStatement = "ALTER USER '{{username}}' VALID UNTIL '{{expiration}}' ON CLUSTER '{{cluster}}';"
	// ExpirationFormat is the ClickHouse VALID UNTIL timestamp layout.
	ExpirationFormat = "2006-01-02 15:04:05"
	// DefaultUserNameTemplate generates Vault dynamic usernames when none is configured.
	DefaultUserNameTemplate = `{{ printf "v-%s-%s-%s-%s" (.DisplayName | truncate 8) (.RoleName | truncate 8) (random 20) (unix_time) | truncate 63 }}`
)
