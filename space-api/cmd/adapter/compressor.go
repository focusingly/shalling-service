package adapter

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

type httpCompressResponseWrapper struct {
	http.ResponseWriter                // 确保接口兼容性的占位
	compressWriter      io.WriteCloser // 实际负责压缩和写入的实例
}

// Write 覆盖实际的 Write 写入方法
func (w *httpCompressResponseWrapper) Write(b []byte) (int, error) {
	return w.compressWriter.Write(b)
}

func newHttpResponseCompressWrapper(originResponseWriter http.ResponseWriter, compressWriter io.WriteCloser) http.ResponseWriter {
	return &httpCompressResponseWrapper{
		ResponseWriter: originResponseWriter,
		compressWriter: compressWriter,
	}
}

// 压缩配置
type compressConfig struct {
	encoding                 string
	getCompressorWriteCloser func(io.Writer) (io.WriteCloser, error)
}

// 默认支持的响应体压缩方案(性能考虑, 均只使用最快的压缩级别)
var presetCompressors = map[string]*compressConfig{
	// 优先使用
	"br": {
		encoding: "br",
		getCompressorWriteCloser: func(w io.Writer) (io.WriteCloser, error) {
			return brotli.NewWriterOptions(w, brotli.WriterOptions{
				Quality: brotli.BestSpeed,
			}), nil
		},
	},
	"zstd": {
		encoding: "zstd",
		getCompressorWriteCloser: func(w io.Writer) (io.WriteCloser, error) {
			return zstd.NewWriter(w, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
		},
	},
	"gzip": {
		encoding: "gzip",
		getCompressorWriteCloser: func(w io.Writer) (io.WriteCloser, error) {
			return gzip.NewWriterLevel(w, gzip.BestSpeed)
		},
	},
	// 兼容性压缩方案
	"deflate": {
		encoding: "deflate",
		getCompressorWriteCloser: func(w io.Writer) (io.WriteCloser, error) {
			return zlib.NewWriterLevel(w, zlib.BestSpeed)
		},
	},
}

// AdaptiveCompressHandler 处理多种压缩算法
type AdaptiveCompressHandler struct {
	handler http.Handler
}

// 选择最佳压缩方法
func (h *AdaptiveCompressHandler) selectClientSupportCompressor(r *http.Request) *compressConfig {
	acceptEncoding := r.Header.Get("Accept-Encoding")
	// 按优先级检查支持的压缩方法
	for _, enc := range []string{"br", "zstd", "gzip", "deflate"} {
		if strings.Contains(acceptEncoding, enc) {
			if cfg, ok := presetCompressors[enc]; ok {
				return cfg
			}
		}
	}

	return nil
}

func (h *AdaptiveCompressHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 尝试寻找最佳的客户端支持的压缩算法
	compressConf := h.selectClientSupportCompressor(r)
	// 如果客户端不支持任何匹配的压缩算法, 那么原样返回
	if compressConf == nil {
		h.handler.ServeHTTP(w, r)
		return
	}

	// 创建压缩数据写入器(为每次请求都会创建一个新的压缩器, 确保无竞争问题)
	compressorWriter, err := compressConf.getCompressorWriteCloser(w)
	if err != nil {
		h.handler.ServeHTTP(w, r)
		return
	}
	defer compressorWriter.Close()

	// 设置响应头
	w.Header().Set("Content-Encoding", compressConf.encoding)
	w.Header().Set("Vary", "Accept-Encoding")
	// 使用压缩器写入数据
	h.handler.ServeHTTP(
		newHttpResponseCompressWrapper(w, compressorWriter),
		r,
	)
}

func NewAdaptiveCompressionHttpWriter(originHandler http.Handler) http.Handler {
	return &AdaptiveCompressHandler{
		handler: originHandler,
	}
}
