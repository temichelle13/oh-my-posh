package segments

import (
	"strings"

	"github.com/jandedobbeleer/oh-my-posh/src/platform"
	"github.com/jandedobbeleer/oh-my-posh/src/properties"
)

type Container struct {
	props properties.Properties
	env   platform.Environment

	Name string
}

func (c *Container) Template() string {
	return " \uf308 {{ .Name }} "
}

func (c *Container) Init(props properties.Properties, env platform.Environment) {
	c.props = props
	c.env = env
}

func (c *Container) Enabled() bool {
	c.Name = c.containerName()

	return len(c.Name) > 0
}

func (c *Container) containerName() string {
	if c.env.HasFiles("/proc/vz") && !c.env.HasFiles("/proc/bc") {
		return "OpenVZ"
	}

	if c.env.HasFiles("/run/host/container-manager") {
		return "OCI"
	}

	// podman and others
	containerEnvPath := "/run/.containerenv"
	if c.env.HasFiles(containerEnvPath) {
		content := c.env.FileContent(containerEnvPath)
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if !strings.HasPrefix(line, `image=\"`) {
				continue
			}

			//TODO: find the correct way to extract the string
			return strings.TrimPrefix(line, `image=\"`)
		}

		return "podman"
	}

	// WSL with systemd will set the contents of this file to "wsl"
	// Avoid showing the container module in that case
	// Honor the contents of this file if "docker" and not running in podman or wsl
	systemdPath := "/run/systemd/container"
	if c.env.HasFiles(systemdPath) {
		content := c.env.FileContent(systemdPath)
		switch strings.TrimSpace(content) {
		case "docker":
			return "Docker"
		case "wsl":
			break
		default:
			return "Systemd"
		}
	}

	if c.env.HasFiles("/.dockerenv") {
		return "Docker"
	}

	return ""
}
