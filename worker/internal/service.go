package internal

import (
	prmt "github.com/gitchander/permutation"
	"github.com/mikhailvzhzhv/crack-hash/shared/models"
)

type Service struct {
}

func (s *Service) ProcessTask(task *models.Task) {
	a := []int{1, 2, 3}
	r := task
	prmt.NewSlicePermutator(a)
}
