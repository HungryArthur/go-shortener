package url

import (
	"fmt"
	"net/http"
)

func jsonErrResp(w http.ResponseWriter, statusCode int, errText string) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(fmt.Sprintf(`
	{
		"error": "%s"
	}
	`, errText)))
}
