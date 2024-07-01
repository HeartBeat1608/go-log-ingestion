package internals

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func RootHandler(cm *DBConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var service string
		for k := range cm.connections {
			service = k
			break
		}
		if service == "" {
			WriteError(w, 500, errors.New("No service is exporting logs. Please export some logs first"))
			return
		}
		service = strings.Split(service, ".")[0]
		http.Redirect(w, r, fmt.Sprintf("/%s/view/all", service), http.StatusSeeOther)
	}
}
