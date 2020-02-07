package pdf

type Version struct {
	major, minor int
}

func (v Version) Set(major, minor int) (err error) {
	v.major = major
	v.minor = minor
	if !validateVersion(major, minor) {
		err = ErrInvalidVersion
	}
	return
}

func (v Version) Get() (int, int) {
	return v.major, v.minor
}

func validateVersion(major, minor int) bool {
	if major == 1 && 0 < minor {
		return minor <= 7
	}
	return false
}
