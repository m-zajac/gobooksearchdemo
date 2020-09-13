package app

import (
	"errors"
	"strings"
)

// ExtractParagraph takes rune slice + index and tries to extract full paragraph around that index.
// Paragraph is a text that ends in \n\n or \r\n\r\n.
func ExtractParagraph(data []rune, charIdx int, maxSize uint) (string, error) {
	if charIdx >= len(data) {
		return "", errors.New("char index too big")
	}

	leftIdx, rightIdx := charIdx, charIdx
	peeker := newNewLinePeeker(data)

	// Move rightIdx to the right until new line is found.
	for i := charIdx; i < len(data) && rightIdx-leftIdx < int(maxSize); i++ {
		rightIdx++
		peeker.SetIdx(i)
		if peeker.NewLineOnRight() {
			break
		}
	}
	// Move leftIdx to the left until new line is found.
	for i := charIdx; i > 0 && rightIdx-leftIdx < int(maxSize); i-- {
		peeker.SetIdx(i)
		peeker.SetIdx(i)
		if peeker.NewLineOnLeft() {
			break
		}
		leftIdx--
	}

	// Return found substring with trimmed whitespace chars.
	result := string(data[leftIdx:rightIdx])
	result = strings.Trim(result, "\r\n\t")

	return result, nil
}

type newLinePeeker struct {
	data     []rune
	startIdx int
	idx      int
}

func newNewLinePeeker(data []rune) *newLinePeeker {
	return &newLinePeeker{
		data: data,
	}
}

func (rp *newLinePeeker) SetIdx(idx int) {
	rp.startIdx = idx
	rp.idx = idx
}

func (rp *newLinePeeker) NewLineOnRight() bool {
	for {
		ok, peek := rp.advanceRight()
		if peek == "\n\n" || peek == "\r\n\r\n" {
			return true
		}
		if !ok || len(peek) >= 4 {
			return false
		}
	}
}

func (rp *newLinePeeker) NewLineOnLeft() bool {
	for {
		ok, peek := rp.advanceLeft()
		if peek == "\n\n" || peek == "\r\n\r\n" {
			return true
		}
		if !ok || len(peek) >= 4 {
			return false
		}
	}
}

func (rp *newLinePeeker) advanceRight() (bool, string) {
	if rp.idx >= len(rp.data)-1 {
		return false, ""
	}
	currCh := rp.data[rp.idx]
	rp.idx++
	switch currCh {
	case '\n', '\r':
		return true, string(rp.data[rp.startIdx:rp.idx])
	default:
		return false, ""
	}
}

func (rp *newLinePeeker) advanceLeft() (bool, string) {
	if rp.idx <= 0 {
		return false, ""
	}
	rp.idx--
	currCh := rp.data[rp.idx]
	switch currCh {
	case '\n', '\r':
		return true, string(rp.data[rp.idx:rp.startIdx])
	default:
		return false, ""
	}
}
