package module

import (
	"math/rand"
	"time"
	"unicode/utf8"
)

const (
	UPPER_ALPHABETS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LOWER_ALPHABETS = "abcdefghijklmnopqrstuvwxyz"
	ALL_ALPHABETS   = UPPER_ALPHABETS + LOWER_ALPHABETS
	NUMBERS         = "0123456789"
	ALPHA_NUMBERIC  = ALL_ALPHABETS + NUMBERS
)

var ALPHA_NUMBERIC_LENGTH int = utf8.RuneCountInString(ALPHA_NUMBERIC)

type GenerateRandom struct {
	r *rand.Rand
}

func (g *GenerateRandom) Init() {
	g.r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (g *GenerateRandom) GetRandom(max int, min int) int {
	return g.r.Intn(max-min) + min
}

func (g *GenerateRandom) GetAlphanumberic(length int) string {
	randStr := make([]byte, 0)
	for i := 0; i < length; i++ {
		randStr = append(randStr, ALPHA_NUMBERIC[g.GetRandom(ALPHA_NUMBERIC_LENGTH, 0)])
	}
	return string(randStr)
}
