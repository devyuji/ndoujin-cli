package utils

import (
	"regexp"
	"strings"
)

func SanitizeFilename(dirtyName string) string {
	// 1. Define a regex that matches illegal characters and control characters
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)

	// 2. Replace all illegal characters with an underscore (or an empty string "")
	safeName := re.ReplaceAllString(dirtyName, "_")

	// 3. Remove trailing periods and spaces (Windows hates these at the end of a folder name)
	safeName = strings.TrimRight(safeName, " .")

	// 4. Fallback in case the server sent a string that was entirely illegal characters
	if safeName == "" {
		return "default_folder_name"
	}

	// 5. Truncate if too long (Most file systems limit names to 255 bytes)
	if len(safeName) > 255 {
		safeName = safeName[:255]
	}

	return safeName
}
