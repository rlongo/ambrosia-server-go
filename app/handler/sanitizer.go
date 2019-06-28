package handler

import (
	"strings"
	
	"github.com/rlongo/ambrosia-server-go/api"
)

func Sanitize(r *api.Recipe) {
	for i, val := range r.Tags {
		r.Tags[i] = strings.ToLower(val)
	}
}