package fileupload

import (
	"fmt"
	"os"
	"path"

	"github.com/majid-cj/go-docker-mongo/util"

	"github.com/twinj/uuid"
)

// FormatFile ...
func FormatFile(filepath string) string {
	ext := path.Ext(filepath)

	return fmt.Sprintf("%s-%s%s", util.GetTimeNow().Format(os.Getenv("TIME_FORMATE")), uuid.NewV4().String(), ext)
}
