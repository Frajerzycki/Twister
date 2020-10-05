package test

// TestingWriter is writer used for tests.
type TestingWriter struct {
	data []byte
}

func (testingWriter *TestingWriter) Write(p []byte) (int, error) {
	testingWriter.data = append(testingWriter.data, p...)
	return len(p), nil
}
