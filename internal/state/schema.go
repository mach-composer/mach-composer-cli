package state

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed schemas/*
var schemas embed.FS

func GetSchema(key Type) (*map[string]any, error) {
	s := map[string]any{}

	switch key {
	case DefaultType:
		fallthrough
	case AwsType:
		loadSchemaNode("schemas/aws.schema.json", &s)
	case AzureType:
		loadSchemaNode("schemas/azure.schema.json", &s)
	case GcpType:
		loadSchemaNode("schemas/gcp.schema.json", &s)
	case HttpType:
		loadSchemaNode("schemas/http.schema.json", &s)
	case LocalType:
		loadSchemaNode("schemas/local.schema.json", &s)
	case TerraformCloudType:
		loadSchemaNode("schemas/terraform_cloud.schema.json", &s)
	default:
		return nil, fmt.Errorf("unknown schema %s", key)
	}

	return &s, nil
}

func loadSchemaNode(filename string, dst any) {
	body, err := schemas.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, dst); err != nil {
		panic(err)
	}
}
