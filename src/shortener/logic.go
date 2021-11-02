package shortener

import (
	"errors"
	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
	"time"
)

var (
	ErrRedirectNotFound = errors.New("Redirect Not Found")
	ErrRedirectInvalid  = errors.New("Redirect Invalid")
)

type redirectService struct {
	redirectRepo RedirectRepository
}

/** Creates a Redirect service */
func NewRedirectService(redirectRepo RedirectRepository) RedirectService {
	return &redirectService{
		redirectRepo,
	}
}

/** Given a code returns a Redirect */
func (r *redirectService) Find(code string) (*Redirect, error) {
	return r.redirectRepo.Find(code)
}

/** Creates a code for a given URL */
func (r *redirectService) Store(redirect *Redirect) error {
	// Check it is in URL format and not empty
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepo.Store(redirect)
}