package util

import (
	"math/big"

	"github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
)

type WordGenerator struct {
	alphabet     []rune
	maxWordLen   int
	startIndex   *big.Int
	endIndex     *big.Int
	currentIndex *big.Int
	alphabetSize int
}

func NewWordGenerator(task *models.Task, batchSize int) *WordGenerator {
	if task.PartCount <= 0 {
		task.PartCount = 1
	}
	if task.PartNumber < 0 {
		task.PartNumber = 0
	}
	if task.PartNumber >= task.PartCount {
		task.PartNumber = task.PartCount - 1
	}

	alphabetRunes := []rune(task.Alphabet)
	alphabetSize := len(alphabetRunes)

	totalWords := calculateTotalWords(alphabetSize, task.MaxWordLen)
	partSize := big.NewInt(int64(batchSize))
	startIndex := new(big.Int).Mul(big.NewInt(int64(task.PartNumber)), partSize)
	endIndex := new(big.Int).Mul(big.NewInt(int64(task.PartNumber+1)), partSize)
	if task.PartNumber == task.PartCount-1 {
		endIndex = new(big.Int).Set(totalWords)
	}

	currentIndex := new(big.Int).Set(startIndex)

	return &WordGenerator{
		alphabet:     alphabetRunes,
		maxWordLen:   task.MaxWordLen,
		startIndex:   startIndex,
		endIndex:     endIndex,
		currentIndex: currentIndex,
		alphabetSize: alphabetSize,
	}
}

func calculateTotalWords(alphabetSize int, maxWordLen int) *big.Int {
	if alphabetSize == 0 {
		return big.NewInt(0)
	}
	if alphabetSize == 1 {
		return big.NewInt(int64(maxWordLen))
	}

	r := big.NewInt(int64(alphabetSize))
	rn := new(big.Int).Exp(r, big.NewInt(int64(maxWordLen)), nil)
	numerator := new(big.Int).Sub(rn, big.NewInt(1))
	numerator.Mul(numerator, r)
	denominator := new(big.Int).Sub(r, big.NewInt(1))

	return new(big.Int).Div(numerator, denominator)
}

func (wg *WordGenerator) Next() (string, bool) {
	if wg.currentIndex.Cmp(wg.endIndex) >= 0 {
		return "", false
	}

	word := wg.indexToWord(wg.currentIndex)
	wg.currentIndex.Add(wg.currentIndex, big.NewInt(1))

	return word, true
}

func (wg *WordGenerator) indexToWord(index *big.Int) string {
	length := 1
	remaining := new(big.Int).Set(index)

	for length <= wg.maxWordLen {
		count := wg.wordsOfLength(length)

		if remaining.Cmp(count) < 0 {
			break
		}

		remaining.Sub(remaining, count)
		length++
	}

	word := make([]rune, length)

	for i := length - 1; i >= 0; i-- {
		divider := big.NewInt(int64(wg.alphabetSize))
		charIndex := new(big.Int).Mod(remaining, divider)
		word[i] = wg.alphabet[charIndex.Int64()]

		remaining.Div(remaining, divider)
	}

	return string(word)
}

func (wg *WordGenerator) wordsOfLength(length int) *big.Int {
	return new(big.Int).Exp(big.NewInt(int64(wg.alphabetSize)), big.NewInt(int64(length)), nil)
}
