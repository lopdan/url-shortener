package msgpack

import (
	"github.com/pkg/errors"
	"github.com/lopdan/url-shortener/src/shortener"
	"github.com/vmihailenco/msgpack"
)

type Redirect struct{}

/** Decodes a message and puts it into the Redirect struct */
func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}
	return redirect, nil
}

/** Given a message, returns it encoded */
func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	return rawMsg, nil
}