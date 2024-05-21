package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vadskev/urlshort/internal/lib/logger/zp"
	"go.uber.org/zap"
)

type compressWriter struct {
	respWriter http.ResponseWriter
	zipWriter  *gzip.Writer
}

func newCompressWriter(rw http.ResponseWriter) *compressWriter {
	return &compressWriter{
		respWriter: rw,
		zipWriter:  gzip.NewWriter(rw),
	}
}

func (cw *compressWriter) Header() http.Header {
	return cw.respWriter.Header()
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	return cw.zipWriter.Write(p)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		cw.respWriter.Header().Set("Content-Encoding", "gzip")
	}
	cw.respWriter.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {
	return cw.zipWriter.Close()
}

/**/

type compressReader struct {
	reqReader io.ReadCloser
	zipReader *gzip.Reader
}

func newCompressReader(rc io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(rc)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		reqReader: rc,
		zipReader: zr,
	}, nil
}

func (cr *compressReader) Read(p []byte) (n int, err error) {
	return cr.zipReader.Read(p)
}

func (cr *compressReader) Close() error {
	if err := cr.reqReader.Close(); err != nil {
		return err
	}
	return cr.zipReader.Close()
}

func New(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			writer := w

			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportGzip := strings.Contains(acceptEncoding, "gzip")

			contentEncoding := r.Header.Get("Content-Encoding")
			sendGzip := strings.Contains(contentEncoding, "gzip")

			if sendGzip && supportGzip {
				cressWriter := newCompressWriter(w)
				writer = cressWriter

				defer func() {
					if err := cressWriter.Close(); err != nil {
						log.Info("Compress", zp.Err(err))
					}
				}()
			}

			if sendGzip && supportGzip {
				respReader, err := newCompressReader(r.Body)
				if err != nil {
					log.Info("Compress", zp.Err(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				r.Body = respReader

				defer func() {
					if err = respReader.Close(); err != nil {
						log.Info("Compress", zp.Err(err))
					}
				}()
			}

			ww := middleware.NewWrapResponseWriter(writer, r.ProtoMajor)
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
