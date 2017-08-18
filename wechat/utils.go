package wechat

import (
	"crypto/rand"
	"fmt"
	"io"
)

func LimitRead(r io.Reader, w io.Writer, max int64) error {

	// Do not read more than max + 1 bytes.
	r = &io.LimitedReader{
		R: r,
		N: max + 1,
	}

	n, err := io.Copy(w, r)

	if err != nil {

		return err

	}

	if n > max {

		return fmt.Errorf("Read excceed max (%v)", max)

	}

	return nil

}

// Generate a cryptographically secure hex string of length n.
func HexCryptoRandString(n int) string {

	buf := make([]byte, (n>>1)+(n&1))

	_, err := rand.Read(buf)

	if err != nil {

		panic(err)

	}

	return fmt.Sprintf("%x", buf)[:n]

}

// log.Logger like.
type Logger interface {
	Printf(format string, v ...interface{})
}

// A logger that log nothing.
type NopLogger bool

func (l NopLogger) Printf(format string, v ...interface{}) {
}
