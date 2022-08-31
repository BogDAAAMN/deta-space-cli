package manifest

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/deta/pc-cli/pkg/util/fs"
	"github.com/deta/pc-cli/shared"
	"gopkg.in/yaml.v3"
)

var (
	// ManifestName manifest file name
	ManifestName = "deta.yml"
)

func IsManifestPresent(sourceDir string) (bool, error) {
	var exists bool
	var err error
	for _, name := range getSupportedManifestNames() {
		exists, err = fs.FileExists(sourceDir, name)
		if err != nil {
			return false, err
		}

		if exists {
			return true, nil
		}
	}
	return false, nil
}

func Open(sourceDir string) (*Manifest, error) {
	var exists bool
	var err error
	for _, name := range getSupportedManifestNames() {
		exists, err = fs.FileExists(sourceDir, name)
		if err != nil {
			return nil, err
		}

		if exists {
			ManifestName = name
			break
		}
	}

	if !exists {
		return nil, ErrManifestNotFound
	}

	// read raw contents from manifest file
	c, err := ioutil.ReadFile(filepath.Join(sourceDir, ManifestName))
	if err != nil {
		return nil, fmt.Errorf("failed to read contents of manifest file: %w", err)
	}

	// parse raw manifest file content
	m := Manifest{}
	err = yaml.Unmarshal([]byte(c), &m)
	if err != nil {
		return nil, fmt.Errorf("failed to do parse manifest file, please check for correct syntax: %w", err)
	}

	return &m, nil
}

// OpenRaw returns the raw manifest file content from sourceDir if it exists
func OpenRaw(sourceDir string) ([]byte, error) {
	var exists bool
	var err error
	for _, name := range getSupportedManifestNames() {
		exists, err = fs.FileExists(sourceDir, name)
		if err != nil {
			return nil, err
		}

		if exists {
			ManifestName = name
			break
		}
	}

	if !exists {
		return nil, ErrManifestNotFound
	}

	// read raw contents from manifest file
	c, err := ioutil.ReadFile(filepath.Join(sourceDir, ManifestName))
	if err != nil {
		return nil, fmt.Errorf("failed to read contents of manifest file: %w", err)
	}

	return c, nil
}

func (m *Manifest) Save(sourceDir string) error {
	// marshall manifest object
	c, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshall manifest object: %w", err)
	}

	// write manifest object to file
	err = ioutil.WriteFile(filepath.Join(sourceDir, ManifestName), c, 0644)
	if err != nil {
		return fmt.Errorf("failed to write manifest object: %w", err)
	}

	return nil
}

func (m *Manifest) AddMicro(newMicro *shared.Micro) error {
	// mark new micro as primary if it is the only one
	if len(m.Micros) == 0 {
		newMicro.Primary = true
	}

	for _, micro := range m.Micros {
		if micro.Name == newMicro.Name {
			return fmt.Errorf("a micro with the same name already exists in \"deta.yml\"")
		}
		if micro.Src == newMicro.Src {
			return fmt.Errorf("another micro already exists at the same location %s in the manifest", newMicro.Src)
		}
	}
	m.Micros = append(m.Micros, newMicro)

	return nil
}

func CreateBlankManifest(sourceDir string) (*Manifest, error) {
	manifest := &Manifest{
		V: 0,
	}

	err := manifest.Save(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create a blank manifest in %s, %w", sourceDir, err)
	}

	return manifest, nil
}
