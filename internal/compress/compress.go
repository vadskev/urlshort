package compress

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type compressWriter struct {
	respWriter  http.ResponseWriter
	gzipWriter  *gzip.Writer
	contentType map[string]bool
	isCompress  bool
}

func newCompressWriter(rw http.ResponseWriter) *compressWriter {
	return &compressWriter{
		respWriter:  rw,
		gzipWriter:  gzip.NewWriter(rw),
		contentType: map[string]bool{"application/json": true, "text/html": true},
		isCompress:  false,
	}
}

func (cw *compressWriter) Header() http.Header {
	return cw.respWriter.Header()
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	return cw.writer().Write(p)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	contentType := cw.respWriter.Header().Get("Content-Type")

	if _, ok := cw.contentType[contentType]; ok {
		cw.respWriter.Header().Set("Content-Encoding", "gzip")
		cw.isCompress = true
	}
	cw.respWriter.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {
	if cw.isCompress {
		return cw.writer().(io.WriteCloser).Close()
	}
	return nil
}

func (cw *compressWriter) writer() io.Writer {
	if cw.isCompress {
		return cw.gzipWriter
	} else {
		return cw.respWriter
	}
}

/**/

type compressReader struct {
	respReadCloser io.ReadCloser
	gzipReader     *gzip.Reader
}

func newCompressReader(rc io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(rc)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		respReadCloser: rc,
		gzipReader:     zr,
	}, nil
}

func (cr compressReader) Read(p []byte) (n int, err error) {
	return cr.gzipReader.Read(p)
}

func (cr *compressReader) Close() error {
	if err := cr.respReadCloser.Close(); err != nil {
		return err
	}
	return cr.gzipReader.Close()
}

func RequestCompress(h http.Handler) http.Handler {
	return http.HandlerFunc(func(irw http.ResponseWriter, r *http.Request) {
		tmpw := irw
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		log.Println(supportsGzip)
		if supportsGzip {
			cw := newCompressWriter(irw)
			tmpw = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {

			cr, err := newCompressReader(r.Body)
			if err != nil {
				irw.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr

			defer func() {
				err = cr.Close()
				if err != nil {
					irw.WriteHeader(http.StatusInternalServerError)
					return
				}
			}()
		}

		h.ServeHTTP(tmpw, r)
	})
}
