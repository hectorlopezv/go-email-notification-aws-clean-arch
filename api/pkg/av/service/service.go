package av

import (
	av "av-send-email/api/pkg/av/repository"
	"av-send-email/api/pkg/entities"
	"encoding/json"
	"fmt"
)

//Interface from which our module can access our repossitory of all our models
type Service interface{
	ClamAvScan(paths []string)(*entities.ScanResult, error)
}


type service struct{
	repository av.Repository

	
}

func(s *service) ClamAvScan(paths []string)(*entities.ScanResult, error){
	result, err := s.repository.ExecuteClamAvScan(paths)
	if err != nil {
	return nil, err
	}
	
	jsonResult, err := json.Marshal(result)
	fmt.Printf("Result: %+v\n", jsonResult)

	resultProcessScan,errScan := s.repository.ProcessScanResult(result)
	if errScan != nil {
		return nil, errScan
	}
	fmt.Printf("Result: %+v\n", resultProcessScan)

	return result, nil
}

func NewService(repository av.Repository) Service{
	return &service{repository}
}