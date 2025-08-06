package pkg

// Bytes returns the given number of bytes.
func Bytes(n int64) int64 {
	return n
}

// KB returns the number of bytes for the given kilobytes.
func KB(n int64) int64 {
	return n * 1024
}

// MB returns the number of bytes for the given megabytes.
func MB(n int64) int64 {
	return n * 1024 * 1024
}

// GB returns the number of bytes for the given gigabytes.
func GB(n int64) int64 {
	return n * 1024 * 1024 * 1024
}

// TB returns the number of bytes for the given terabytes.
func TB(n int64) int64 {
	return n * 1024 * 1024 * 1024 * 1024
}
