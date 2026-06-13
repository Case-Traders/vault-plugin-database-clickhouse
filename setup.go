package clickhouse

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/helper/template"
)

func setupUsernameTemplate(config map[string]interface{}) (template.StringTemplate, error) {
	var empty template.StringTemplate
	raw, err := strutil.GetString(config, ConfigKeyUsernameTemplate)
	if err != nil {
		return empty, fmt.Errorf("retrieve %s: %w", ConfigKeyUsernameTemplate, err)
	}
	if raw == "" {
		raw = DefaultUserNameTemplate
	}
	up, err := template.NewTemplate(template.Template(raw))
	if err != nil {
		return empty, fmt.Errorf("initialize username template: %w", err)
	}
	if _, err := up.Generate(dbplugin.UsernameMetadata{}); err != nil {
		return empty, fmt.Errorf("invalid username template: %w", err)
	}
	return up, nil
}

func clusterFromConfig(config map[string]interface{}) (string, error) {
	return strutil.GetString(config, ConfigKeyCluster)
}

func formatExpiration(t time.Time) string {
	return t.Format(ExpirationFormat)
}
