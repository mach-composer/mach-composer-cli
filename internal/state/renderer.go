package state

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	"strings"
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
	// Identifier returns the full identifier for the renderer. This can be used to fetch a renderer for a node
	Identifier() string
	// StateKey returns the terraform state key for the renderer
	StateKey() string
	// Backend returns the terraform backend configuration for the renderer
	Backend() (string, error)
	// RemoteState returns the terraform remote state data configuration for the renderer
	RemoteState() (string, error)
}

type BaseRenderer struct {
	identifier string
	stateKey   string
}

func (br *BaseRenderer) Identifier() string {
	return br.identifier
}

func (br *BaseRenderer) StateKey() string {
	return br.stateKey
}

func NewRenderer(typ Type, identifier string, data map[string]any) (Renderer, error) {
	//We only use the last part of the identifier as the state key.
	keyParts := strings.Split(identifier, "/")
	if len(keyParts) < 1 {
		return nil, fmt.Errorf("invalid identifier %s", identifier)
	}
	key := keyParts[len(keyParts)-1]

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
				stateKey:   key,
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
				stateKey:   key,
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
				stateKey:   key,
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
				stateKey:   key,
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
				stateKey:   key,
			},
			state: state,
		}, nil
	}

	return nil, fmt.Errorf("unknown state type %s", typ)
}
