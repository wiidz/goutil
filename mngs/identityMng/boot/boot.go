package boot

import (
	"sync"

	idauthn "github.com/wiidz/goutil/mngs/identityMng/authn"
	idfacade "github.com/wiidz/goutil/mngs/identityMng/facade"
	idsetup "github.com/wiidz/goutil/mngs/identityMng/setup"
	idsubhandler "github.com/wiidz/goutil/mngs/identityMng/subject/handler"
	idsubgorm "github.com/wiidz/goutil/mngs/identityMng/subject/repo/gorm"
	idsubsvc "github.com/wiidz/goutil/mngs/identityMng/subject/service"
	"gorm.io/gorm"
)

var once sync.Once

func ensureSetup() { once.Do(func() { idsetup.Init() }) }

// NewFacade builds a ready-to-use identity facade with default subject type.
// It also ensures Sa-Token is initialized once.
func NewFacade(db *gorm.DB, defaultSubjectType string) *idfacade.Service {
	ensureSetup()
	auth := idauthn.New()
	subRepo := idsubgorm.NewSubjectRepository(db)
	return idfacade.New(auth, subRepo, defaultSubjectType)
}

// NewSubjectService returns a subject service wired with gorm repository.
func NewSubjectService(db *gorm.DB) *idsubsvc.Service {
	return idsubsvc.New(idsubgorm.NewSubjectRepository(db))
}

// NewSubjectHandler returns a Gin handler for subject directory.
func NewSubjectHandler(db *gorm.DB) *idsubhandler.Handler {
	return idsubhandler.New(NewSubjectService(db))
}
