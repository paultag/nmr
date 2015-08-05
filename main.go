package main

import (
	"fmt"
	"os"

	// "pault.ag/go/nmr/helpers"
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

	dsc := "/home/tag/tmp/repo/pool/main/f/fbautostart/fbautostart_2.718281828-3.dsc"

	cmd, err := build.BuildCommand(dsc)
	if err != nil {
		panic(err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	ftbfs := err == nil

	// if ftbfs {
	// 	helpers.LogChangesFromDsc()
	// } else {
	// 	helpers.AppendLogToChanges()
	// }
	//
	// changes.Move(incoming)
}
