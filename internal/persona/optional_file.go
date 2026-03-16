package persona

import "os"

// readOptionalFile reads the file at path and returns its content as a string.
// If the file does not exist, it returns an empty string and nil error.
// Other errors (e.g. permission denied) are returned as-is.
func readOptionalFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(content), nil
}
