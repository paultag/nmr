package build

import (
	"pault.ag/go/debian/control"
	"pault.ag/go/debian/dependency"
	"pault.ag/go/nmr/candidate"
	"pault.ag/go/reprepro"
)

type BuildStatus struct {
	Package   reprepro.BuildNeedingPackage
	Buildable bool
}

func ComputeBuildStatus(
	repo reprepro.Repo,
	index candidate.Canidates,
	packages []reprepro.BuildNeedingPackage,
) []BuildStatus {
	ret := []BuildStatus{}

	for _, pkg := range packages {
		dsc, err := control.ParseDscFile(repo.Basedir + "/" + pkg.Location)
		if err != nil {
			panic("OMGWTFALSKJFALSKJ")
			continue
		}

		arch, err := dependency.ParseArch(pkg.Arch)
		if err != nil {
			panic("OMGWTFALSKJFALSKJ")
			continue
		}

		ret = append(ret, BuildStatus{
			Package: pkg,
			Buildable: index.SatisfiesBuildDepends(
				*arch,
				dsc.BuildDepends,
			),
		})
	}

	return ret
}
