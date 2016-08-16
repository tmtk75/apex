// Package infra proxies Terraform commands.
package infra

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/infra"
	"github.com/apex/apex/stats"
)

// example output.
const example = `  View change plan
  $ apex infra plan

  Apply changes
  $ apex infra apply`

// Command config.
var Command = &cobra.Command{
	Use:     "infra",
	Short:   "Infrastructure management",
	Example: example,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)
}

// Run command.
func run(c *cobra.Command, args []string) error {
	stats.Track("Infra", nil)

	err := root.Project.LoadFunctions()

	// Hack to prevent initial `apex infra apply` from failing,
	// as we load functions to expose their ARNs.
	if err != nil {
		if !strings.Contains(err.Error(), "Role: zero value") {
			return err
		}
	}

	p := &infra.Proxy{
		ProjectName: root.Project.Name,
		Functions:   root.Project.Functions,
		Region:      *root.Session.Config.Region,
		Environment: root.Project.Environment,
		Role:        root.Project.Role,
	}

	return p.Run(args...)
}
