package terraform

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	workingPath string
}

type version struct {
	TerraformVersion string `json:"terraform_version"`
}

func NewCmd(workingPath string) *Cmd {
	return &Cmd{
		workingPath: workingPath,
	}
}

func (c *Cmd) GetVersion() (string, error) {
	output, err := c.runTerraformCommand("version", "-json")
	if err != nil {
		return "", err
	}

	var version version
	if err := json.Unmarshal(output, &version); err != nil {
		return "", err
	}

	return version.TerraformVersion, nil
}

func (c *Cmd) GeneratePlan(planPath string) error {
	if _, err := c.runTerraformCommand("plan", "-out", planPath); err != nil {
		return err
	} else {
		return nil
	}
}

func (c *Cmd) ConvertPlanToJson(planPath string, jsonPath string) error {
	output, err := c.runTerraformCommand("show", "-json", planPath)
	if err != nil {
		return err
	}

	// keep relative pathing so I don't lose my mind
	if strings.HasPrefix(jsonPath, "./") || strings.HasPrefix(jsonPath, "../") {
		jsonPath = c.workingPath + "/" + jsonPath
	}

	if err := os.WriteFile(jsonPath, output, 0644); err != nil {
		return err
	}

	return nil
}

func (c *Cmd) DoImport(resourceName string, resourceId string) error {
	if _, err := c.runTerraformCommand("import", resourceName, resourceId); err != nil {
		return err
	}

	return nil
}

func (c *Cmd) DoStateRm(resourceName string) error {
	if _, err := c.runTerraformCommand("state", "rm", resourceName); err != nil {
		return err
	}

	return nil
}

func (c *Cmd) runTerraformCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = c.workingPath

	if output, err := cmd.Output(); err != nil {
		return nil, err
	} else {
		return output, nil
	}
}
