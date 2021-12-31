package commands

import (
	"strings"
)

var (
	SkeletonVersion = "1.1.20"
	CLIVersion      = "1.1.21"
)

var (
	VersionBig   = 1
	VersionSmall = 2
	VersionEqual = 0
)

func VersionCompare(verA, verB string) int {
	verStrArrA := splitStrByNet(verA)
	verStrArrB := splitStrByNet(verB)
	lenStrA := len(verStrArrA)
	lenStrB := len(verStrArrB)
	if lenStrA != lenStrB {
		panic("version length is inconsistent")
	}
	return compareArrStrVers(verStrArrA, verStrArrB)
}

func compareArrStrVers(verA, verB []string) int {
	for index, _ := range verA {
		littleResult := compareLittleVer(verA[index], verB[index])
		if littleResult != VersionEqual {
			return littleResult
		}
	}
	return VersionEqual
}

func compareLittleVer(verA, verB string) int {
	bytesA := []byte(verA)
	bytesB := []byte(verB)
	lenA := len(bytesA)
	lenB := len(bytesB)
	if lenA > lenB {
		return VersionBig
	}
	if lenA < lenB {
		return VersionSmall
	}
	return compareByBytes(bytesA, bytesB)
}

func compareByBytes(verA, verB []byte) int {
	for index, _ := range verA {
		if verA[index] > verB[index] {
			return VersionBig
		}
		if verA[index] < verB[index] {
			return VersionSmall
		}

	}
	return VersionEqual
}

func splitStrByNet(strV string) []string {
	return strings.Split(strV, ".")
}
