package app

import (
	"errors"
	"strings"
)

// ExtractParagraph takes rune slice + index and tries to extract full paragraph around that index.
// Paragraph is a text that ends in \n\n.
func ExtractParagraph(data []rune, charIdx int, maxSize uint) (string, error) {
	if charIdx >= len(data) {
		return "", errors.New("char index too big")
	}

	leftIdx, rightIdx := charIdx, charIdx
	peeker := newPeeker(data)

	for i := charIdx; i < len(data) && rightIdx-leftIdx < int(maxSize); i++ {
		rightIdx++
		peeker.SetIdx(i)
		if peeker.NewLineOnRight() {
			break
		}
	}

	for i := charIdx; i > 0 && rightIdx-leftIdx < int(maxSize); i-- {
		peeker.SetIdx(i)
		if peeker.NewLineOnLeft() {
			break
		}
		leftIdx--
	}

	result := string(data[leftIdx:rightIdx])
	result = strings.Trim(result, "\n\t")

	return result, nil
}

type runePeeker struct {
	data     []rune
	startIdx int
	idx      int
}

func newPeeker(data []rune) *runePeeker {
	return &runePeeker{
		data: data,
	}
}
func (rp *runePeeker) SetIdx(idx int) {
	rp.startIdx = idx
	rp.idx = idx
}

func (rp *runePeeker) PeekLeft() (bool, string) {
	if rp.idx <= 0 {
		return false, ""
	}
	rp.idx--

	return true, string(rp.data[rp.idx:rp.startIdx])
}

func (rp *runePeeker) NewLineOnLeft() bool {
	for {
		ok, peek := rp.PeekLeft()
		if peek == "\n\n" || peek == "\r\n\r\n" {
			return true
		}
		if !ok || len(peek) >= 4 {
			return false
		}
	}
}

func (rp *runePeeker) PeekRight() (bool, string) {
	if rp.idx >= len(rp.data)-1 {
		return false, ""
	}
	rp.idx++

	return true, string(rp.data[rp.startIdx:rp.idx])
}

func (rp *runePeeker) NewLineOnRight() bool {
	for {
		ok, peek := rp.PeekRight()
		if peek == "\n\n" || peek == "\r\n\r\n" {
			return true
		}
		if !ok || len(peek) >= 4 {
			return false
		}
	}
}
