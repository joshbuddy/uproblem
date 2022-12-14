package assetserver

import (
	"net/http"
	"path/filepath"
	"sync"

	"github.com/wailsapp/mimetype"
)

var (
	cache = map[string]string{}
	mutex sync.Mutex
)

func GetMimetype(filename string, data []byte) string {
	mutex.Lock()
	defer mutex.Unlock()

	// short-circuit .js, .css to ensure the
	// browser evaluates them in the right context
	switch filepath.Ext(filename) {
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css; charset=utf-8"
	}

	result := cache[filename]
	if result != "" {
		return result
	}

	detect := mimetype.Detect(data)
	if detect == nil {
		result = http.DetectContentType(data)
	} else {
		result = detect.String()
	}

	if result == "" {
		result = "application/octet-stream"
	}

	cache[filename] = result
	return result
}
