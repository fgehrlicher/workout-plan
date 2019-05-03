package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func IsGreater(baseVersion string, comparatorVersion string) (bool, error) {
	baseVersionStruct, err := NewVersion(baseVersion)
	if err != nil {
		return false, err
	}

	comparatorVersionStrict, err := NewVersion(comparatorVersion)
	if err != nil {
		return false, err
	}

	MajorBigger := comparatorVersionStrict.Major > baseVersionStruct.Major
	MajorEqual := comparatorVersionStrict.Major == baseVersionStruct.Major

	MinorBigger := comparatorVersionStrict.Minor > baseVersionStruct.Minor
	MinorEqual := comparatorVersionStrict.Minor == baseVersionStruct.Minor

	PatchBigger := comparatorVersionStrict.Patch > baseVersionStruct.Patch
	PatchEqual := comparatorVersionStrict.Patch == baseVersionStruct.Patch

	if MajorBigger || (MajorEqual && MinorBigger) || (MajorEqual && MinorEqual && PatchBigger) {
		return true, nil
	}

	if MajorEqual && MinorEqual && PatchEqual {
		return false, errors.New(
			fmt.Sprintf(
				"can´t compare. versions are identical: %v",
				baseVersion,
			),
		)
	}

	return false, nil
}

func NewVersion(rawVersion string) (Version, error) {
	version := Version{}

	versionSlice := strings.Split(rawVersion, ".")
	if len(versionSlice) != 3 {
		return version, errors.New(
			fmt.Sprintf(
				"invalid version format: `%v` ( should be x.x.x )",
				rawVersion,
			),
		)
	}

	major, err := strconv.Atoi(versionSlice[0])
	if err != nil {
		return version, err
	}
	version.Major = major

	minor, err := strconv.Atoi(versionSlice[0])
	if err != nil {
		return version, err
	}
	version.Minor = minor

	patch, err := strconv.Atoi(versionSlice[0])
	if err != nil {
		return version, err
	}
	version.Patch = patch

	return version, nil
}

type Version struct {
	Major int
	Minor int
	Patch int
}
