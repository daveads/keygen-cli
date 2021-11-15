package cmd

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	"github.com/spf13/cobra"
)

var (
	genkeyOpts = &CommandOptions{}
	genkeyCmd  = &cobra.Command{
		Use:   "genkey",
		Short: "generate an ed25519 key pair for code signing",
		Args:  cobra.NoArgs,
		RunE:  genkeyRun,

		// Encountering an error should not display usage
		SilenceUsage: true,
	}
)

func init() {
	genkeyCmd.Flags().StringVar(&genkeyOpts.signingKey, "out", "keygen.key", "output the private publishing key to specified file")
	genkeyCmd.Flags().StringVar(&genkeyOpts.verifyKey, "pubout", "keygen.pub", "output the public upgrade key to specified file")

	rootCmd.AddCommand(genkeyCmd)
}

func genkeyRun(cmd *cobra.Command, args []string) error {
	signingKeyPath, err := homedir.Expand(genkeyOpts.signingKey)
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, genkeyOpts.signingKey, err)
	}

	verifyKeyPath, err := homedir.Expand(genkeyOpts.verifyKey)
	if err != nil {
		return fmt.Errorf(`path "%s" is not expandable (%s)`, genkeyOpts.verifyKey, err)
	}

	if _, err := os.Stat(signingKeyPath); err == nil {
		return fmt.Errorf(`signing key file "%s" already exists`, signingKeyPath)
	}

	if _, err := os.Stat(verifyKeyPath); err == nil {
		return fmt.Errorf(`verify key file "%s" already exists`, verifyKeyPath)
	}

	verifyKey, signingKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}

	if err := writeSigningKeyFile(signingKeyPath, signingKey); err != nil {
		return err
	}

	if err := writeVerifyKeyFile(verifyKeyPath, verifyKey); err != nil {
		return err
	}

	if abs, err := filepath.Abs(signingKeyPath); err == nil {
		signingKeyPath = abs
	}

	if abs, err := filepath.Abs(verifyKeyPath); err == nil {
		verifyKeyPath = abs
	}

	fmt.Printf(`Private publishing key: %s
Public upgrade key: %s

Warning: never share your publishing key, it's a secret!
`,
		signingKeyPath,
		verifyKeyPath,
	)

	return nil
}

func writeSigningKeyFile(signingKeyPath string, signingKey ed25519.PrivateKey) error {
	file, err := os.OpenFile(signingKeyPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := hex.EncodeToString(signingKey)
	_, err = file.WriteString(enc)
	if err != nil {
		return err
	}

	return nil
}

func writeVerifyKeyFile(verifyKeyPath string, verifyKey ed25519.PublicKey) error {
	file, err := os.OpenFile(verifyKeyPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := hex.EncodeToString(verifyKey)
	_, err = file.WriteString(enc)
	if err != nil {
		return err
	}

	return nil
}
