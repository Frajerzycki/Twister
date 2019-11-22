package test

type TestingWriter struct {
	data []byte
}

func (testingWriter *TestingWriter) Write(p []byte) (int, error) {
	testingWriter.data = p
	return len(p), nil
}
