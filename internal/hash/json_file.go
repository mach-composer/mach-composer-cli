package hash

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"sync"
)

var mutex = &sync.RWMutex{}

type Hashes map[string]string

type JsonFileHandler struct {
	file string
}

func NewJsonFileHandler(file string) Handler {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Info().Msgf("Creating new hash file %s", file)
		err = os.MkdirAll(path.Dir(file), 0777)
		if err != nil {
			log.Panic().Err(err).Msgf("Failed to create directory %s", path.Dir(file))
		}

		err = os.WriteFile(file, []byte("{}"), 0777)
		if err != nil {
			log.Panic().Err(err).Msgf("Failed to create file %s", file)
		}
	}

	return &JsonFileHandler{
		file: file,
	}
}

func (h *JsonFileHandler) getHashes() (*Hashes, error) {
	var hashes Hashes

	f, err := os.OpenFile(h.file, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c, err := os.ReadFile(f.Name())
	if err != nil {
		return nil, err
	}

	if len(c) == 0 {
		return &Hashes{}, nil
	}

	err = json.Unmarshal(c, &hashes)
	if err != nil {
		return nil, err
	}

	return &hashes, nil
}

func (h *JsonFileHandler) Fetch(_ context.Context, n graph.Node) (string, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	hashes, err := h.getHashes()
	if err != nil {
		return "", err
	}

	switch n.Type() {
	case graph.ProjectType:
		return "", nil
	case graph.SiteType:
		s := n.(*graph.Site)
		graph.SortSiteComponentNodes(s.NestedNodes)

		var componentHashes []string
		for _, component := range s.NestedNodes {
			h := (*hashes)[component.Identifier()]
			componentHashes = append(componentHashes, h)
		}

		return utils.ComputeHash(componentHashes)
	case graph.SiteComponentType:
		return (*hashes)[n.Identifier()], nil
	default:
		return "", fmt.Errorf("unknown node type %T", n)
	}
}

func (h *JsonFileHandler) Store(_ context.Context, n graph.Node) error {
	mutex.Lock()
	defer mutex.Unlock()

	hashes, err := h.getHashes()
	if err != nil {
		return err
	}

	switch n.Type() {
	case graph.ProjectType:
		return nil
	case graph.SiteType:
		for _, nn := range n.(*graph.Site).NestedNodes {
			(*hashes)[nn.Identifier()], err = nn.Hash()
			if err != nil {
				return err
			}
		}
	case graph.SiteComponentType:
		(*hashes)[n.Identifier()], err = n.Hash()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown node type %T", n)
	}

	c, err := json.Marshal(hashes)
	if err != nil {
		return err
	}

	return os.WriteFile(h.file, c, 0777)
}
