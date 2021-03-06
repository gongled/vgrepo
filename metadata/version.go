package metadata

import (
	"sort"

	"pkg.re/essentialkaos/ek.v9/sortutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// VMetadataVersionsList is a list of VMetadataVersion structs
type VMetadataVersionsList []*VMetadataVersion

// VMetadataVersion provides struct for version entity
type VMetadataVersion struct {
	Version   string                 `json:"version"`   // version of the image
	Providers VMetadataProvidersList `json:"providers"` // list of available providers
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Len implements interface method for Sort
func (s VMetadataVersionsList) Len() int {
	return len(s)
}

// Swap implements interface method for Sort
func (s VMetadataVersionsList) Swap(i, j int) {
	*s[i], *s[j] = *s[j], *s[i]
}

// Less implements interface method for Sort
func (s VMetadataVersionsList) Less(i, j int) bool {
	return sortutil.VersionCompare(s[i].Version, s[j].Version)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// equalVersions returns true if versions are equal
func equalVersions(first *VMetadataVersion, second *VMetadataVersion) bool {
	return first.Version == second.Version
}

// notEqualVersions returns true if versions are not equal
func notEqualVersions(first *VMetadataVersion, second *VMetadataVersion) bool {
	return !equalVersions(first, second)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// AnyVersions returns true if version if present on the list
func (m *VMetadata) AnyVersions(version *VMetadataVersion, f func(*VMetadataVersion, *VMetadataVersion) bool) bool {
	for _, v := range m.Versions {
		if f(v, version) {
			return true
		}
	}
	return false
}

// isVersionExist returns true if version is already exist in the metadata
// TODO: improve comparisons: version -> checksum
func (m *VMetadata) isVersionExist(version *VMetadataVersion) bool {
	return m.AnyVersions(version, equalVersions)
}

// CountVersions returns number of available versions
func (m *VMetadata) CountVersions() int {
	return len(m.Versions)
}

// OldestVersion returns the oldest version from the list
func (m *VMetadata) OldestVersion() string {
	if !m.IsEmptyMeta() {
		return m.Versions[0].Version
	}

	return ""
}

// LatestVersion returns the latest version from the list
func (m *VMetadata) LatestVersion() *VMetadataVersion {
	if !m.IsEmptyMeta() {
		return m.Versions[m.CountVersions()-1]
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// SortVersions sorts list of versions in the metadata
func (m *VMetadata) SortVersions() {
	sort.Sort(VMetadataVersionsList(m.Versions))
}

// AddVersion adds version to the metadata list
func (m *VMetadata) AddVersion(version *VMetadataVersion) error {
	if m.isVersionExist(version) {
		versionMatch := m.FindVersion(version.Version)
		for _, p := range version.Providers {
			return versionMatch.AddProvider(p)
		}
	} else {
		m.Versions = append(m.Versions, version)
	}

	m.SortVersions()

	return nil
}

// FilterVersion filters list of versions in the metadata by given function
func (m *VMetadata) FilterVersion(version *VMetadataVersion, f func(*VMetadataVersion, *VMetadataVersion) bool) {
	versionsList := make(VMetadataVersionsList, 0)

	for _, v := range m.Versions {
		if f(v, version) {
			versionsList = append(versionsList, v)
		}
	}

	m.Versions = versionsList
}

// FindVersion returns version by string id
func (m *VMetadata) FindVersion(version string) *VMetadataVersion {
	for _, v := range m.Versions {
		if v.Version == version {
			return v
		}
	}
	return nil
}

// RemoveVersion removes version from the list or do nothing
func (m *VMetadata) RemoveVersion(version *VMetadataVersion) {
	m.FilterVersion(version, notEqualVersions)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewMetadataVersion returns new VMetadataVersion struct
func NewMetadataVersion(version string, providers VMetadataProvidersList) *VMetadataVersion {
	m := &VMetadataVersion{
		Version:   version,
		Providers: providers,
	}

	return m
}
