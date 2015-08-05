package main

import (
	"fmt"
	"os"

	"pault.ag/go/debian/control"
	"pault.ag/go/nmr/helpers"
	"pault.ag/go/nmr/repo"
	"pault.ag/go/reprepro"
	"pault.ag/go/sbuild"
)

func main() {
	fmt.Printf(`

		░█▀█░█▀█░█░█░█▀▀░█▀▄░░░█▄█░█▀█░█░░░█▀▀░░░█▀▄░█▀█░▀█▀
		░█░█░█▀█░█▀▄░█▀▀░█░█░░░█░█░█░█░█░░░█▀▀░░░█▀▄░█▀█░░█░
		░▀░▀░▀░▀░▀░▀░▀▀▀░▀▀░░░░▀░▀░▀▀▀░▀▀▀░▀▀▀░░░▀░▀░▀░▀░░▀░


`)
	arch := "amd64"
	dscFile := "/home/tag/tmp/repo/pool/main/f/fbautostart/fbautostart_2.718281828-3.dsc"
	dsc, err := control.ParseDscFile(dscFile)
	if err != nil {
		panic(err)
	}

	repreproRepo := reprepro.NewRepo(os.Args[1])
	config, err := repo.LoadConfig(repreproRepo.Basedir)
	if err != nil {
		panic(err)
	}

	suiteConfig, err := repo.LoadDistributions(repreproRepo.Basedir)
	if err != nil {
		panic(err)
	}

	suite := os.Args[2]

	distConfig, err := config.GetDistConfig(suite)
	if err != nil {
		panic(err)
	}

	suiteDistConfig, err := suiteConfig.GetDistConfig(suite)
	if err != nil {
		panic(err)
	}

	build := sbuild.NewSbuild(suite, suite)
	build.Verbose()
	build.AddArgument("build-dep-resolver", "aptitude")
	build.AddArgument("chroot-setup-commands",
		fmt.Sprintf("apt-key add /schroot/%s.asc", suiteDistConfig.SignWith))

	build.AddArgument("extra-repository",
		fmt.Sprintf("deb %s %s main",
			config.Global.PublicArchiveRoot,
			distConfig.Upstream.Dist,
		),
	)

	cmd, err := build.BuildCommand(dscFile)
	if err != nil {
		panic(err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	ftbfs := err != nil

	// control file paths have no epoch.
	//
	// ... because that makes sense.
	var logVersion string
	if dsc.Version.IsNative() {
		logVersion = fmt.Sprintf("%d", dsc.Version.Version)
	} else {
		logVersion = fmt.Sprintf("%s-%s", dsc.Version.Version, dsc.Version.Revision)
	}

	logPath := fmt.Sprintf(
		"%s_%s_%s.build",
		dsc.Source,
		logVersion,
		arch,
	)

	if ftbfs {
		changes, err := helpers.LogChangesFromDsc(logPath, dscFile, suite, arch)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", changes)
	} // else {
	// 	helpers.AppendLogToChanges()
	// }
	//
	// changes.Move(incoming)
}
