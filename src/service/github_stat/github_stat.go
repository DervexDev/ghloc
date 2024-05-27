package github_stat

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/DervexDev/ghloc/src/server/rest"
	"github.com/DervexDev/ghloc/src/service/loc_count"
	"github.com/DervexDev/ghloc/src/util"
	"github.com/rs/zerolog"
)

var ErrNoData = fmt.Errorf("no data")
var ErrRepoTooLarge = rest.BadRequest{Msg: "Too large repository"}

type LOCCacher interface {
	SetLOCs(ctx context.Context, user, repo, branch string, locs []loc_count.LOCForPath) error
	SetTooLarge(ctx context.Context, user, repo, branch string) error
	GetLOCs(ctx context.Context, user, repo, branch string) ([]loc_count.LOCForPath, error) // error may be ErrNoData
}

type TempStorage int

const (
	TempStorageFile TempStorage = iota
	TempStorageRam
)

type Opener = func() (io.ReadCloser, error)

type FileForPath struct {
	Path   string
	Opener Opener
}

type ContentProvider interface {
	GetContent(ctx context.Context, user, repo, branch, token string, tempStorage TempStorage) (_ []FileForPath, close func() error, _ error)
}

type Service struct {
	ContentProvider ContentProvider
}

func New(contentProvider ContentProvider) *Service {
	return &Service{
		ContentProvider: contentProvider,
	}
}

func (s *Service) GetStat(ctx context.Context, user, repo, branch, token string, filter, matcher *string, noLOCProvider bool, tempStorage TempStorage) (*loc_count.StatTree, error) {
	var locs *[]loc_count.LOCForPath

	filesForPaths, close, err := s.ContentProvider.GetContent(ctx, user, repo, branch, token, tempStorage)
	if err != nil {
		return nil, fmt.Errorf("get repo content: %w", err)
	}
	defer close()

	locCountingStart := time.Now()

	locCounter := loc_count.NewFilesLOCCounter()
	for _, file := range filesForPaths {
		err := func() error {
			fileReader, err := file.Opener()
			if err != nil {
				return err
			}
			defer fileReader.Close()
			return locCounter.AddFile(file.Path, fileReader)
		}()
		if err != nil {
			return nil, err
		}
	}

	locs = util.Pointer(locCounter.GetLOCsForPaths())
	zerolog.Ctx(ctx).Info().
		Float64("durationSec", time.Since(locCountingStart).Seconds()).
		Msg("LOCs counted")

	treeBuildingStart := time.Now()
	statTree := loc_count.BuildStatTree(*locs, filter, matcher)
	zerolog.Ctx(ctx).Info().
		Float64("durationSec", time.Since(treeBuildingStart).Seconds()).
		Msg("Statistics tree built")

	return statTree, nil
}
