package repo

import (
	"pault.ag/go/debian/control"
	"pault.ag/go/debian/dependency"
	"pault.ag/go/reprepro"
	"pault.ag/go/resolver"
)

type BuildStatus struct {
	Package   reprepro.BuildNeedingPackage
	Buildable bool
	Why       string
}

func ComputeBuildStatus(
	repo reprepro.Repo,
	index resolver.Canidates,
	packages []reprepro.BuildNeedingPackage,
) []BuildStatus {
	ret := []BuildStatus{}

	for _, pkg := range packages {
		dsc, err := control.ParseDscFile(repo.Basedir + "/" + pkg.Location)
		if err != nil {
			continue
		}

		arch, err := dependency.ParseArch(pkg.Arch)
		if err != nil {
			/// XXX: ERROR OUT
			continue
		}

		buildable, why := index.ExplainSatisfiesBuildDepends(*arch, dsc.BuildDepends)

		ret = append(ret, BuildStatus{
			Package:   pkg,
			Buildable: buildable,
			Why:       why,
		})
	}

	return ret
}
