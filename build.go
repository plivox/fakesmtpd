//go:build tools

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	s "github.com/plivox/gulmy/shell"
	v "github.com/plivox/gulmy/version"
	"github.com/spf13/cobra"
)

var (
	name       = "fakesmtpd"
	sources    = []string{"cmd/main.go"}
	version    string
	buildDir   string
	releaseDir string

	binaryFlagBuildDir   string
	binaryFlagArch       string
	binaryFlagOS         string
	binaryFlagCC         string
	binaryFlagCGOEnabled bool
	binaryFlagLDFlags    []string

	binaryCmd = &cobra.Command{
		Use:   "binary",
		Short: "Build binary",
		Run: func(cmd *cobra.Command, args []string) {
			binaryCmdRun()
		},
	}

	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clean build directory",
		Run: func(cmd *cobra.Command, args []string) {
			cleanCmdRun()
		},
	}

	rootCmd = &cobra.Command{Use: "build"}
)

func init() {
	s.MakeStyle()
	cobra.OnInitialize(initConfig)

	binaryCmd.Flags().StringVar(&binaryFlagArch, "arch", "", "Architecture (GOARCH)")
	binaryCmd.Flags().StringVar(&binaryFlagOS, "os", "", "Platform (GOOS)")
	binaryCmd.Flags().StringVar(&binaryFlagCC, "cc", "", "Flag CC")
	binaryCmd.Flags().BoolVar(&binaryFlagCGOEnabled, "cgo", false, "Enable cgo")

	rootCmd.PersistentFlags().StringVarP(&buildDir, "build-dir", "b", "build", "Build directory")
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(binaryCmd)
}

func initConfig() {
	version = v.FromFile()

	fmt.Println(version)

	releaseDir = s.Join(buildDir, "release")
	binaryFlagLDFlags = []string{
		fmt.Sprintf("-X %s/internal/cmd.Version=%s", name, version),
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func cleanCmdRun() {
	s.Remove(buildDir)
}

func binaryCmdRun() {
	var kernels []string

	if binaryFlagOS != "" {
		kernels = []string{binaryFlagOS}
	} else {
		kernels = []string{s.Windows, s.Darwin, s.Linux}
	}

	s.Mkdir(buildDir)

	if binaryFlagCGOEnabled {
		os.Setenv("CGO_ENABLED", "1")
	}

	for _, kernel := range kernels {
		os.Setenv("GOOS", kernel)

		var architectures []string

		if binaryFlagArch != "" {
			architectures = []string{binaryFlagArch}
		} else {
			switch kernel {
			case s.Windows:
				architectures = []string{"amd64"}
			case s.Darwin:
				architectures = []string{"amd64", "arm64"}
			case s.Linux:
				architectures = []string{"amd64", "386", "arm64"}
			}
		}

		for _, arch := range architectures {
			os.Setenv("GOARCH", arch)

			target := fmt.Sprintf("%s/%s-%s-%s-%s", releaseDir, name, version, kernel, arch)
			if kernel == s.Windows {
				target += ".exe"
			}

			args := []string{"build"}
			if len(binaryFlagLDFlags) != 0 {
				args = append(args, "-ldflags", fmt.Sprintf("%s", strings.Join(binaryFlagLDFlags, " ")))
			}

			args = append(args, "-o", target, strings.Join(sources, " "))
			s.Cmd("go", args...).Run()
		}
	}
}
