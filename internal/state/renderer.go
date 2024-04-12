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
	Identifier() string
	Key() string
	Backend() (string, error)
	RemoteState() (string, error)
}

type BaseRenderer struct {
	identifier string
	key        string
}

func (br *BaseRenderer) Identifier() string {
	return br.identifier
}

func (br *BaseRenderer) Key() string {
	return br.key
}

func NewRenderer(typ Type, identifier, key string, data map[string]any) (Renderer, error) {
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
			BaseRenderer: BaseRenderer{
				identifier: identifier,
				key:        key,
			},
			state: state,
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
			BaseRenderer: BaseRenderer{
				identifier: identifier,
				key:        key,
			},
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
			BaseRenderer: BaseRenderer{
				identifier: identifier,
				key:        key,
			},
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
			BaseRenderer: BaseRenderer{
				identifier: identifier,
				key:        key,
			},
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
			BaseRenderer: BaseRenderer{
				identifier: identifier,
				key:        key,
			},
			state: state,
		}, nil
	}

	return nil, fmt.Errorf("unknown state type %s", typ)
}
