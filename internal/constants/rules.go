package constants

import "events-api/pkg/utils"

var EnvValidationRules = []utils.ValidationRule{
	// Server validation
	{
		Variable: "PORT",
		Default:  "3005",
		Rule:     utils.IsValidPort,
		Message:  "server port is required and must be a valid port number",
	},
	{
		Variable: "GO_ENV",
		Default:  "development",
		Rule:     func(v string) bool { return v == "development" || v == "production" },
		Message:  "GO_ENV must be either 'development' or 'production'",
	},

	// Database validation
	{
		Variable: "DB_URI",
		Rule:     func(v string) bool { return v != "" },
		Message:  "database uri is required",
	},
	{
		Variable: "DB_NAME",
		Default:  "events",
		Rule:     func(v string) bool { return v != "" },
		Message:  "database name is required",
	},
}
