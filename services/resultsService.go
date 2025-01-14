package services

import (
	"net/http"
	"runner/models"
	"time"
)

type ResultsService struct {
	runnersRepository *repositories.RunnersRepository
	resultsRepository *repositories.ResultsRepository
}

func NewResultsService(runnersRepository *repositories.RunnersRepository,
	resultsRepository *repositories.ResultsRepository) *ResultsService {
	return &ResultsService{
		runnersRepository: runnersRepository,
		resultsRepository: resultsRepository,
	}
}

func parseRaceResult(timeString string) (time.Duration, error) {
	return time.ParseDuration(
		timeString[0:2] + "h" +
			timeString[3:5] + "m" +
			timeString[6:8] + "s")
}

func (rs ResultsService) CreateResult(result *models.Result) (*models.Result, *models.ResponseError) {
	currentYear := time.Now().Year()

	if result.RunnerID == "" {
		return nil, &models.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid runner ID",
		}
	}
	if result.RaceResult == "" {
		return nil, &models.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid race result",
		}
	}
	if result.Location == "" {
		return nil, &models.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid race location",
		}
	}
	if result.Year > currentYear || result.Year < 0 {
		return nil, &models.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid year",
		}
	}
	raceResult, err := parseRaceResult(result.RaceResult)
	if err != nil {
		return nil, &models.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid race result",
		}
	}
	response, responseErr := rs.resultsRepository.CreateResult(result)
	if responseErr != nil {
		return nil, responseErr
	}
	runner, responseErr := rs.runnersRepository.GetRunner(result.RunnerID)
	if responseErr != nil {
		return nil, responseErr
	}
	if runner == nil {
		return nil, &models.ResponseError{
			Status:  http.StatusNotFound,
			Message: "Runner not found",
		}
	}
	//update runner's personal best
	if runner.PersonalBest == "" {
		runner.PersonalBest = result.RaceResult
	} else {
		personalBest, err := parseRaceResult(runner.PersonalBest)
		if err != nil {
			return nil, &models.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to parse personal best result",
			}
		}
		if raceResult < personalBest {
			runner.PersonalBest = result.RaceResult
		}
	}
	//update runner's season best
	if result.Year == currentYear {
		if runner.SeasonBest == "" {
			runner.SeasonBest = result.RaceResult
		}
		seasonBest, err := parseRaceResult(runner.SeasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to parse personal best result",
			}
		}
		if raceResult < seasonBest {
			runner.SeasonBest = result.RaceResult
		}
	}
	responseErr = rs.runnersRepository.UpdateRunnerResults(runner)
	if responseErr != nil {
		return nil, responseErr
	}
	return response, nil
}

func (rs ResultsService) DeleteResult(resultID string) *models.ResponseError {
	if resultID == "" {
		return &models.ResponseError{
			Status:  http.StatusBadRequest,
			Message: "Invalid result ID",
		}
	}
	result, responseErr := rs.resultsRepository.DeleteRunner(resultID)
	if responseErr != nil {
		return responseErr
	}
	runner, responseErr := rs.runnersRepository.GetRunner(resultID)
	if responseErr != nil {
		return responseErr
	}
	//Checking if deleted result is personal best for the runner
	if runner.PersonalBest == result.RaceResult {
		personalBest, responseErr := rs.resultsRepository.GetPersonalBestResults(result.RunnerID)
		if responseErr != nil {
			return responseErr
		}
		runner.PersonalBest = personalBest
	}
	//Checking if deleted result is season best for the runner
	currentYear := time.Now().Year()
	if result.Year == currentYear && runner.SeasonBest == result.RaceResult {
		seasonBest, responseErr := rs.resultsRepository.GetSeasonBestResults(result.RunnerID, result.Year)
		if responseErr != nil {
			return responseErr
		}
		runner.SeasonBest = seasonBest
	}
	responseErr = rs.runnersRepository.UpdateRunnerResults(runner)
	if responseErr != nil {
		return responseErr
	}
	return nil
}
