package services

import (
	"net/http"
	"runners/models"
	"runners/repositories"
	"strconv"
	"time"
)

type RunnersService struct {
	runnersRepository *repositories.RunnersRepository
	resultsRepository *repositories.ResultsRepository
}

func NewRunnersService(
	runnersRepository *repositories.RunnersRepository,
	resultsRepository *repositories.ResultsRepository,
) *RunnersService {
	return &RunnersService{
		runnersRepository: runnersRepository,
		resultsRepository: resultsRepository,
	}
}

func (rs RunnersService) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	resErr := validateRunner(runner)
	if resErr != nil {
		return nil, resErr
	}
	return rs.runnersRepository.CreateRunner(runner)
}

func (rs RunnersService) UpdateRunner(runner *models.Runner) *models.ResponseError {
	resErr := valdateRunnerId(runner.ID)
	if resErr != nil {
		return resErr
	}

	resErr = validateRunner(runner)
	if resErr != nil {
		return resErr
	}
	return rs.runnersRepository.UpdateRunner(runner)
}

func (rs RunnersService) DeleteRunner(runnerId string) *models.ResponseError {
	resErr := valdateRunnerId(runnerId)
	if resErr != nil {
		return resErr
	}
	return rs.runnersRepository.DeleteRunner(runnerId)
}

func (rs RunnersService) GetRunner(runnerId string) (*models.Runner, *models.ResponseError) {
	resErr := valdateRunnerId(runnerId)
	if resErr != nil {
		return nil, resErr
	}

	runner, resErr := rs.runnersRepository.GetRunner(runnerId)
	if resErr != nil {
		return nil, resErr
	}

	results, resErr := rs.resultsRepository.GetAllRunnersResults(runnerId)
	if resErr != nil {
		return nil, resErr
	}
	runner.Results = results
	return runner, nil
}

func (rs RunnersService) GetRunnersBatch(country string, year string) ([]*models.Runner, *models.ResponseError) {
	if country != "" && year != "" {
		return nil, &models.ResponseError{
			Message: "Only one parameter can be passed",
			Status:  http.StatusBadRequest,
		}
	}

	if country != "" {
		return rs.runnersRepository.GetRunnersByCountry(country)
	}

	if year != "" {
		intYear, err := strconv.Atoi(year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}
		currentYear := time.Now().Year()
		if intYear < 0 || intYear > currentYear {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}

		return rs.runnersRepository.GetRunnersByYear(intYear)
	}

	return rs.runnersRepository.GetAllRunners()
}

func validateRunner(runner *models.Runner) *models.ResponseError {
	if runner.FirstName == "" {
		return &models.ResponseError{
			Message: "Invalid first name",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.LastName == "" {
		return &models.ResponseError{
			Message: "Invalid last name",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.Age <= 16 || runner.Age > 125 {
		return &models.ResponseError{
			Message: "Invalid age",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.Country == "" {
		return &models.ResponseError{
			Message: "Invalid country",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}

func valdateRunnerId(runnerId string) *models.ResponseError {
	if runnerId == "" {
		return &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}
	return nil
}
