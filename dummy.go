package buildinfo

const (
	DummyVersion = "0.0.0"
	DummyDate    = "1997-08-29 13:37:00"
	DummyBranch  = "HEAD"
	DummyCommit  = "abcdef"
)

func Dummy(bld *BuildInfo) {
	if bld.Version == "" {
		bld.Version = DummyVersion
	}
	if bld.Date == "" {
		bld.Date = DummyDate
	}
	if bld.Branch == "" {
		bld.Branch = DummyBranch
	}
	if bld.Commit == "" {
		bld.Commit = DummyCommit
	}
}
