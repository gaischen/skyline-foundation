package protocol

import "math"

type VersionNumber uint32

// The version numbers, making grepping easier
const (
	VersionTLS      VersionNumber = 0x51474fff
	VersionWhatever VersionNumber = 1 // for when the version doesn't matter
	VersionUnknown  VersionNumber = math.MaxUint32
	VersionDraft29  VersionNumber = 0xff00001d
	VersionDraft32  VersionNumber = 0xff000020
)

var SupportedVersions = []VersionNumber{VersionTLS}


// IsValidVersion says if the version is known to quic-go
func IsValidVersion(v VersionNumber) bool {
	return v == VersionTLS || IsSupportedVersion(SupportedVersions, v)
}

// IsSupportedVersion returns true if the server supports this version
func IsSupportedVersion(supported []VersionNumber, v VersionNumber) bool {
	for _, t := range supported {
		if t == v {
			return true
		}
	}
	return false
}