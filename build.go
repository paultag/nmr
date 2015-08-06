package main

import (
	"fmt"
	"os"
	"os/exec"

	"pault.ag/go/nmr/repo"
	"pault.ag/go/reprepro"
	"pault.ag/go/sbuild"
)

func GetBuildNeeding(repoRoot, suite, arch string) []repo.BuildStatus {
	repreproRepo := reprepro.NewRepo(repoRoot)

	config, err := repo.LoadConfig(repreproRepo.Basedir)
	i, err := config.LoadIndex(suite)

	if err != nil {
		return []repo.BuildStatus{}
	}

	needsBuild, err := repreproRepo.BuildNeeding(suite, arch)
	if err != nil {
		return []repo.BuildStatus{}
	}

	return repo.ComputeBuildStatus(*repreproRepo, *i, needsBuild)
}

func IsArchAllArch(repoRoot, arch string) bool {
	repreproRepo := reprepro.NewRepo(repoRoot)
	config, err := repo.LoadConfig(repreproRepo.Basedir)
	if err != nil {
		return false
	}
	return arch == config.Global.ArchIndepBuildArch
}

func GetIncoming(repoRoot, dist string) (string, error) {
	repreproRepo := reprepro.NewRepo(repoRoot)
	config, err := repo.LoadConfig(repreproRepo.Basedir)
	if err != nil {
		return "", err
	}

	distConfig, err := config.GetDistConfig(dist)
	if err != nil {
		return "", err
	}
	return distConfig.Incoming, nil
}

func SbuildCommand(
	chroot,
	suite,
	arch,
	dscFile,
	repoRoot string,
	verbose bool,
) (*exec.Cmd, error) {
	repreproRepo := reprepro.NewRepo(repoRoot)
	config, err := repo.LoadConfig(repreproRepo.Basedir)
	if err != nil {
		return nil, err
	}

	suiteConfig, err := repo.LoadDistributions(repreproRepo.Basedir)
	if err != nil {
		return nil, err
	}

	distConfig, err := config.GetDistConfig(suite)
	if err != nil {
		return nil, err
	}

	suiteDistConfig, err := suiteConfig.GetDistConfig(suite)
	if err != nil {
		return nil, err
	}

	build := sbuild.NewSbuild(suite, suite)

	if verbose {
		build.Verbose()
	}

	if arch == config.Global.ArchIndepBuildArch {
		build.AddFlag("--arch-all")
	} else {
		build.AddFlag("--no-arch-all")
	}

	// build.AddArgument("build-dep-resolver", "aptitude")
	build.AddArgument("chroot-setup-commands",
		fmt.Sprintf("apt-key add /schroot/%s.asc", suiteDistConfig.SignWith))

	build.AddArgument("extra-repository",
		fmt.Sprintf("deb %s %s main",
			config.Global.PublicArchiveRoot,
			distConfig.Upstream.Dist))

	cmd, err := build.BuildCommand(dscFile)

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd, err
}
