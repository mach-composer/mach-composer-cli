package cmd

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Basic imports

func cleanWorkingDir(workdir string) {
	err := os.RemoveAll(path.Join(workdir, "deployments"))
	if err != nil {
		panic(err)
	}
	err = os.RemoveAll(path.Join(workdir, "states"))
	if err != nil {
		panic(err)
	}
	err = os.RemoveAll(path.Join(workdir, "hashes.json"))
	if err != nil {
		panic(err)
	}
}

func CompareDirectories(dir1, dir2 string) error {
	// Read the file info of both directories
	files1, err := os.ReadDir(dir1)
	if err != nil {
		return err
	}

	files2, err := os.ReadDir(dir2)
	if err != nil {
		return err
	}

	// Iterate over files in both directories
	for _, file1 := range files1 {
		file1Path := filepath.Join(dir1, file1.Name())

		// Check if the file exists in the second directory
		found := false
		for _, file2 := range files2 {
			if file2.Name() == file1.Name() {
				file2Path := filepath.Join(dir2, file2.Name())

				// Compare file contents if both are regular files
				if file1.Type().IsRegular() && file2.Type().IsRegular() {
					if err := compareFiles(file1Path, file2Path); err != nil {
						return err
					}
				}

				// If both are directories, recursively compare them
				if file1.Type().IsDir() && file2.Type().IsDir() {
					if err := CompareDirectories(file1Path, file2Path); err != nil {
						return err
					}
				}

				found = true
				break
			}
		}

		// If file1 is not found in dir2
		if !found {
			fmt.Printf("%s exists in %s but not in %s\n", file1.Name(), dir1, dir2)
		}
	}

	// Check for files in dir2 that are not in dir1
	for _, file2 := range files2 {
		found := false
		for _, file1 := range files1 {
			if file1.Name() == file2.Name() {
				found = true
				break
			}
		}

		// If file2 is not found in dir1
		if !found {
			fmt.Printf("%s exists in %s but not in %s\n", file2.Name(), dir2, dir1)
		}
	}

	return nil
}

func compareFiles(file1, file2 string) error {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	tplCtx := struct {
		PWD string
	}{
		PWD: path,
	}

	content1, err := os.ReadFile(file1)
	if err != nil {
		return err
	}
	parsedContent1, err := utils.RenderGoTemplate(string(content1), tplCtx)
	if err != nil {
		return err
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		return err
	}
	parsedContent2, err := utils.RenderGoTemplate(string(content2), tplCtx)
	if err != nil {
		return err
	}

	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(
		strings.TrimSpace(string(parsedContent1)),
		strings.TrimSpace(string(parsedContent2)),
		true,
	)

	if len(diffs) != 1 {
		return fmt.Errorf("file %s and %s differ:\n\n %s:\n\n%s\n\n%s\n\n%s", file1, file2, file1, parsedContent1, file2, parsedContent2)
	}

	return nil
}
