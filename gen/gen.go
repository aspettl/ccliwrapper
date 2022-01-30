package gen

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/aspettl/ccliwrapper/cfg"
	"github.com/aspettl/ccliwrapper/tpl"
)

type ToolParams struct {
	Engine       string
	Name         string
	ImageName    string
	ImageTag     cfg.ImageTagConfig
	WorkDir      string
	HomeDir      string
	Command      cfg.CommandConfig
	Mounts       []cfg.MountConfig
	Env          []cfg.EnvConfig
	CustomScript string
}

// Generate writes a shell script based on the template and the tool config
func Generate(outputDir, templateFile string, toolParams ToolParams) error {
	t, err := loadTemplate(templateFile)
	if err != nil {
		return err
	}

	tmpFileName := path.Join(outputDir, toolParams.Name+".tmp")
	fileName := path.Join(outputDir, toolParams.Name)

	if err := os.Remove(tmpFileName); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	f, err := os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, toolParams)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return os.Rename(tmpFileName, fileName)
}

// GenerateAlias produces a symlink to the specified tool with the given alias name
func GenerateAlias(outputDir, toolName, aliasName string) error {
	tmpFileName := path.Join(outputDir, toolName+".tmp")
	aliasFileName := path.Join(outputDir, aliasName)

	if err := os.Remove(tmpFileName); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	err := os.Symlink(toolName, tmpFileName)
	if err != nil {
		return err
	}

	return os.Rename(tmpFileName, aliasFileName)
}

func loadTemplate(templateFile string) (*template.Template, error) {
	if templateFile != "" {
		return template.ParseFiles(templateFile)
	} else {
		return template.New("root").Parse(tpl.WrapperScript)
	}
}
