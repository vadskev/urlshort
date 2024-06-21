package filestorage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vadskev/urlshort/internal/storage"
	"go.uber.org/zap"
)

func TestMemStorage_SaveURL(t *testing.T) {
	ctx := context.Background()
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}
	log := zap.Must(cfg.Build())

	store := NewFileStorage("/tmp/file.json", log)

	err := store.SaveURL(ctx, storage.URLData{
		URL:    "https://ya.ru/4",
		ResURL: "https://ya.ru/sdfsdf",
		Alias:  "sdfsdf",
	})
	require.NoError(t, err)

}

func TestMemStorage_GetURL(t *testing.T) {
	ctx := context.Background()
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}
	log := zap.Must(cfg.Build())

	store := NewFileStorage("/tmp/file.json", log)
	err := store.Get(ctx, store)
	require.NoError(t, err)

	err = store.SaveURL(ctx, storage.URLData{
		URL:    "https://ya.ru/2",
		ResURL: "https://ya.ru/sdfsdf",
		Alias:  "sdfsdf",
	})
	require.NoError(t, err)

	_, err = store.GetURL(ctx, "sdfsdf")
	require.NoError(t, err)
}

func TestNew(t *testing.T) {
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}
	log := zap.Must(cfg.Build())

	store := NewFileStorage("/tmp/file.json", log)
	require.NotNil(t, store)
}
