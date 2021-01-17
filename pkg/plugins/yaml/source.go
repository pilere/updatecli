package yaml

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Source return the latest version
func (y *Yaml) Source() (string, error) {
	// By default workingDir is set to local directory
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	workingDir := filepath.Dir(pwd)

	if y.Value != "" {
		fmt.Println("WARNING: Key 'Value' is not used by source YAML")
	}

	if dir, base, err := isFileExist(y.File); err == nil && y.Path == "" {
		// if no scm configuration has been provided and neither file path then we try to guess the file directory.
		// if file name contains a path then we use it otherwise we fallback to the current path
		y.Path = dir
		y.File = base
	} else if _, _, err := isFileExist(y.File); err != nil && y.Path == "" {

		y.Path = workingDir

	} else if y.Path != "" && !isDirectory(y.Path) {

		fmt.Printf("Directory '%s' is not valid so fallback to '%s'", y.Path, workingDir)
		y.Path = workingDir

	} else {
		return "", fmt.Errorf("Something weird happened while trying to set working directory")
	}

	data, err := y.ReadFile()
	if err != nil {
		return "", err
	}

	var out yaml.Node

	err = yaml.Unmarshal(data, &out)

	if err != nil {
		return "", fmt.Errorf("cannot unmarshal data: %v", err)
	}

	valueFound, value, _ := replace(&out, strings.Split(y.Key, "."), y.Value, 1)

	if !valueFound {
		fmt.Printf("\u2717 cannot find key '%s' from file '%s'\n",
			y.Key,
			filepath.Join(y.Path, y.File))
		return "", nil

	}

	fmt.Printf("\u2714 Key '%s', from file '%v', is correctly set to %s'\n",
		y.Key,
		filepath.Join(y.Path, y.File),
		y.Value)
	return value, nil

}
