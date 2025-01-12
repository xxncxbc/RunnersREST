package services

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

func