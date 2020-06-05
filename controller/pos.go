package pos

import (
	"github.com/gorilla/mux"
	"github.com/tiuriandy/ITISAssignment2/storage"
)

type PosEngine struct {
	Storage *storage.Storage
}

func (pos PosEngine) Route(r *mux.Router) {

}
