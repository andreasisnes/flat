package flat

import (
	"fmt"
)

func concatKey(key, idx string, options *Options) string {
	if key == "" {
		return idx
	}

	return fmt.Sprint(options.Fold(key), options.Delimiter, options.Fold(idx))
}
