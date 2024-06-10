package record

// Config holds the configuration for recording
type Config struct {
	CaptureMethod string

	// Record config
	RecordFPS        int
	RecordResolution Resolution
}

type Resolution struct {
	X, Y int
}
