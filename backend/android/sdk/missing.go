package sdk

// Missing is the error a partial install produces. It is an error and also a
// report: it names each thing that is wrong and, where one exists, the
// sdkmanager line that puts it right.
type Missing struct {
	Root     string
	Problems []Problem
}
