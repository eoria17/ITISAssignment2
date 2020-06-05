package pos

import (
	"github.com/ITISAssignment2/storage"
	"github.com/gorilla/mux"
)

type PosEngine struct {
	Storage *storage.Storage
}

func (pos PosEngine) Route(r *mux.Router) {

}
