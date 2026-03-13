package service

import "github.com/mikhailvzhzhv/crack-hash/shared/v2/models"

type ResultSender interface {
	SendResult(result *models.Result)
}
