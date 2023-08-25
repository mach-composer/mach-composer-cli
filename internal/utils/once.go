package utils

import "sync"

// OnceMap allows for creating a map of sync.Once elements. This allows for running multiple different processes only once.
//
// A use case is for example a list of downloads, where we only want to run every unique download only once.
//
// var downloadFiles = utils.OnceMap{}
//
//	for _, file := files {
//		downloadFiles.Get(file).Do(func() {
//			downloadFile(ctx, file)
//		})
//	}
type OnceMap[K comparable] map[K]*sync.Once

// Get will search for an existing key. If it exists the sync.Once will be returned. If not a new sync.
// Once will be created and set on the map
func (m *OnceMap[any]) Get(key any) *sync.Once {
	im := *m

	if _, ok := im[key]; !ok {
		im[key] = &sync.Once{}
	}

	return im[key]
}
