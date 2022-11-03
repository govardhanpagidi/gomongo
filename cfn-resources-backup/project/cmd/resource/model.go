// Code generated by 'cfn generate', changes will be undone by the next invocation. DO NOT EDIT.
// Updates to this type are made my editing the schema file and executing the 'generate' command.
package resource

// Model is autogenerated from the json schema
type Model struct {
	Name                      *string           `json:",omitempty"`
	OrgId                     *string           `json:",omitempty"`
	ProjectOwnerId            *string           `json:",omitempty"`
	WithDefaultAlertsSettings *bool             `json:",omitempty"`
	Id                        *string           `json:",omitempty"`
	Created                   *string           `json:",omitempty"`
	ClusterCount              *int              `json:",omitempty"`
	ProjectSettings           *ProjectSettings  `json:",omitempty"`
	ApiKeys                   *ApiKeyDefinition `json:",omitempty"`
	ProjectTeams              []ProjectTeam     `json:",omitempty"`
	ProjectApiKeys            []ProjectApiKey   `json:",omitempty"`
}

// ProjectSettings is autogenerated from the json schema
type ProjectSettings struct {
	IsCollectDatabaseSpecificsStatisticsEnabled *bool `json:",omitempty"`
	IsDataExplorerEnabled                       *bool `json:",omitempty"`
	IsPerformanceAdvisorEnabled                 *bool `json:",omitempty"`
	IsRealtimePerformancePanelEnabled           *bool `json:",omitempty"`
	IsSchemaAdvisorEnabled                      *bool `json:",omitempty"`
}

// ApiKeyDefinition is autogenerated from the json schema
type ApiKeyDefinition struct {
	PublicKey  *string `json:",omitempty"`
	PrivateKey *string `json:",omitempty"`
}

// ProjectTeam is autogenerated from the json schema
type ProjectTeam struct {
	TeamId    *string  `json:",omitempty"`
	RoleNames []string `json:",omitempty"`
}

// ProjectApiKey is autogenerated from the json schema
type ProjectApiKey struct {
	Key       *string  `json:",omitempty"`
	RoleNames []string `json:",omitempty"`
}