package services

import (
	"database/sql"
	"github.com/development-raul/footy-predictor/src/domains/seasons"
	"github.com/development-raul/footy-predictor/src/providers/api_sports_provider"
	"github.com/development-raul/footy-predictor/src/utils/resterror"
	"github.com/development-raul/footy-predictor/src/zlog"
)

type SeasonServiceI interface {
	Create(id int64) resterror.RestErrorI
	Find(id int64) (*seasons.Season, resterror.RestErrorI)
	List(req *seasons.ListSeasonInput) ([]seasons.Season, resterror.RestErrorI)
	Delete(id int64) resterror.RestErrorI
	Sync() resterror.RestErrorI
}

type seasonService struct{}

var SeasonService SeasonServiceI = &seasonService{}

func (s *seasonService) Create(id int64) resterror.RestErrorI {
	if err := seasons.SeasonDao.Create(id); err != nil {
		return resterror.NewStandardInternalServerError()
	}
	return nil
}

func (s *seasonService) Find(id int64) (*seasons.Season, resterror.RestErrorI) {
	res, err := seasons.SeasonDao.Find(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, resterror.NewStandardInternalServerError()
	}
	return res, nil
}

func (s *seasonService) List(req *seasons.ListSeasonInput) ([]seasons.Season, resterror.RestErrorI) {
	results, err := seasons.SeasonDao.List(req)
	if err != nil && err != sql.ErrNoRows {
		return nil, resterror.NewStandardInternalServerError()
	}

	return results, nil
}

func (s *seasonService) Delete(id int64) resterror.RestErrorI {
	if err := seasons.SeasonDao.Delete(id); err != nil {
		return resterror.NewStandardInternalServerError()
	}
	return nil
}

func (s *seasonService) Sync() resterror.RestErrorI {
	zlog.Logger.Info("Sync Seasons Start")
	// Get existing seasons
	results, err := seasons.SeasonDao.List(&seasons.ListSeasonInput{Order: "asc"})
	if err != nil && err != sql.ErrNoRows {
		return resterror.NewStandardInternalServerError()
	}
	// Create a map with existing seasons, so we can easily identify the already existing seasons
	existingSeasons := make(map[int64]int64, len(results))
	for _, v := range results {
		existingSeasons[v.ID] = v.ID
	}

	// Get the list of seasons from API Sports
	res, apiErr := api_sports_provider.GetSeasons()
	if apiErr != nil {
		return resterror.NewStandardInternalServerError()
	}

	for _, id := range res {
		// Check if the season already exists
		if _, exists := existingSeasons[id]; exists {
			continue
		}
		// Create the season if it does not exist
		if err := seasons.SeasonDao.Create(id); err != nil {
			zlog.Logger.Warn("could not create season: ", id)
			continue
		}
		zlog.Logger.Info("created new season: ", id)
	}
	zlog.Logger.Info("Sync Seasons End")
	return nil
}
