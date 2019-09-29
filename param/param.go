package param

import "net/http"

type Param interface {
	Format(w http.ResponseWriter, r *http.Request) error
}
