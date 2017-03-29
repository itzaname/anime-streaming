package routines

// SetupRoutines will setup a routine top be ran
func SetupRoutines() {
	go func() {
		SetupDownloads()
	}()
}
