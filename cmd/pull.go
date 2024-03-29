package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/deitch/eci-distribution/pkg/registry"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	pullDir string
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull an ECI from a registry to a local directory",
	Long:  `pull an Edge Container Image (ECI) from an OCI compliant registry`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		// must be exactly one arg, the URL to the manifest
		if len(args) != 1 {
			log.Fatal("must be exactly one arg, the name of the image to download")
		}
		image := args[0]
		desc, err := registry.Pull(image, pullDir, verbose, os.Stdout)
		if err != nil {
			log.Fatalf("error pulling from registry: %v", err)
		}
		fmt.Printf("Pulled image %s with digest %s to directory %s\n", image, string(desc.Digest), pullDir)

	},
}

func pullInit() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pullCmd.Flags().StringVar(&pullDir, "dir", cwd, "directory where to install the ECI, optional")
	pullCmd.Flags().BoolVar(&debug, "debug", false, "debug output")
	pullCmd.Flags().BoolVar(&verbose, "verbose", false, "verbose output")
}
