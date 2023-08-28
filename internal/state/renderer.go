package state

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
)

type Type string

const (
	DefaultType        Type = ""
	LocalType          Type = "local"
	AwsType            Type = "aws"
	GcpType            Type = "gcp"
	AzureType          Type = "azure"
	TerraformCloudType Type = "terraform_cloud"
)

type Renderer interface {
	Key() string
	Backend() (string, error)
	RemoteState() (string, error)
}

func NewRenderer(typ Type, key string, data map[string]any) (Renderer, error) {
	switch typ {
	case DefaultType:
		//Fallthrough to local
		fallthrough
	case LocalType:
		state := &LocalState{}
		if err := mapstructure.Decode(data, state); err != nil {
			return nil, err
		}
		if err := defaults.Set(state); err != nil {
			return nil, err
		}
		return &LocalRenderer{
			state: state,
			key:   key,
		}, nil
	case AwsType:
		state := &AwsState{}
		if err := mapstructure.Decode(data, state); err != nil {
			return nil, err
		}
		if err := defaults.Set(state); err != nil {
			return nil, err
		}
		return &AwsRenderer{
			key:   key,
			state: state,
		}, nil
	case GcpType:
		state := &GcpState{}
		if err := mapstructure.Decode(data, state); err != nil {
			return nil, err
		}
		if err := defaults.Set(state); err != nil {
			return nil, err
		}
		return &GcpRenderer{
			key:   key,
			state: state,
		}, nil
	case AzureType:
		state := &AzureState{}
		if err := mapstructure.Decode(data, state); err != nil {
			return nil, err
		}
		if err := defaults.Set(state); err != nil {
			return nil, err
		}
		return &AzureRenderer{
			key:   key,
			state: state,
		}, nil
	case TerraformCloudType:
		state := &TerraformCloudState{}
		if err := mapstructure.Decode(data, state); err != nil {
			return nil, err
		}
		if err := defaults.Set(state); err != nil {
			return nil, err
		}
		return &TerraformCloudRenderer{
			key:   key,
			state: state,
		}, nil
	}

	return nil, fmt.Errorf("unknown state type %s", typ)
}
