package cliconfig_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/jxsl13/cli-config-boilerplate/cliconfig"
	"github.com/spf13/cobra"
)

func TestCliParse(t *testing.T) {

	port := "9090"

	_ = os.Setenv("SVC_APP_PORT", port)
	buf := bytes.NewBuffer(nil)

	cli := NewTestCmd()
	cli.SetOutput(buf)
	err := cli.Execute()
	if err != nil {
		t.Error(err)
	}

	output := buf.String()
	if output != port+"\n" {
		t.Errorf("expected output %s, got %s", port, output)
	}
}

func NewTestCmd() *cobra.Command {
	cli := CLI{
		Config: Config{
			AppPort: 8080,
		},
	}
	cmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.PreRunE = cli.PreRunE(cmd)
	cmd.RunE = cli.RunE

	return cmd
}

type Config struct {
	AppPort int `koanf:"app.port"`
}

func (c *Config) Validate() error {
	if c.AppPort < 0 || c.AppPort > 65535 {
		return fmt.Errorf("invalid port number")
	}
	return nil
}

type CLI struct {
	Config Config
}

func (c *CLI) PreRunE(cmd *cobra.Command) func(cmd *cobra.Command, args []string) error {
	parse := cliconfig.RegisterFlags(&c.Config, false, cmd, cliconfig.WithEnvPrefix("SVC_"))
	return func(cmd *cobra.Command, args []string) error {
		return parse()
	}
}

func (c *CLI) RunE(cmd *cobra.Command, args []string) error {
	fmt.Fprintln(cmd.OutOrStdout(), c.Config.AppPort)
	return nil
}
