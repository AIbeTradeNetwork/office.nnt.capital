package domain

import (
	"crypto/rand"
	"encoding/base32"
	mathRand "math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/anandvarma/namegen"

	"server/internal/config"
)

const (
	SystemUserUID = "SYSTEM"
)

type Int interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

func GenUID(l int) string {
	byt := make([]byte, 20)
	_, err := rand.Read(byt)
	if err != nil {
		return ""
	}
	return base32.HexEncoding.EncodeToString(byt)[:l]
}

func GenNickname() string {
	nameSchema := []namegen.DictType{
		namegen.Adjectives,
		namegen.Colors,
		namegen.Animals,
	}
	ngen := namegen.NewWithPostfixId(nameSchema, namegen.Numeric, 4)
	return ngen.Get()
}

func Pow[I Int](base, exp I) I {
	result := I(1)
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}

	return result
}

func Percent[I Int](base, percent I) I {
	return I(float64(base) / 10000.0 * float64(percent))
}

func Min[I Int](val1, val2 I) I {
	if val1 > val2 {
		return val2
	}
	return val1
}

func Diff[I Int](val1, val2 I) I {
	if val1 > val2 {
		return val1 - val2
	}
	return val2 - val1
}

func DivCeil[I Int](val, div I) I {
	return (val + (div - 1)) / div
}

func BeginningOfWeek(t time.Time) time.Time {
	nowConfig := config.GetNowConfig()
	return nowConfig.With(t).BeginningOfWeek()
}

func BeginningOfMonth(t time.Time) time.Time {
	nowConfig := config.GetNowConfig()
	return nowConfig.With(t).BeginningOfMonth()
}

func BeginningOfNextMonth(t time.Time) time.Time {
	nowConfig := config.GetNowConfig()
	return nowConfig.With(nowConfig.With(t).EndOfMonth().Add(24 * time.Hour)).BeginningOfMonth()
}

func EndOfMonth(t time.Time) time.Time {
	nowConfig := config.GetNowConfig()
	return nowConfig.With(t).EndOfMonth()
}

func EndOfNextMonth(t time.Time) time.Time {
	nowConfig := config.GetNowConfig()
	return nowConfig.With(nowConfig.With(t).EndOfMonth().Add(24 * time.Hour)).EndOfMonth()
}

func EndOfWeek(t time.Time) time.Time {
	nowConfig := config.GetNowConfig()
	return nowConfig.With(t).EndOfWeek()
}

func BeginningOfNextWeek(t time.Time) time.Time {
	return BeginningOfWeek(t.Add(7 * 24 * time.Hour))
}

func GenSafeSecret(l int) string {
	sec := make([]string, 0, l)
	for i := 0; i < l; i++ {
		sec = append(sec, "*")
	}
	return strings.Join(sec, ",")
}

func GenSafeCode(l int) string {
	numbers := make([]string, 0, l)
	for i := 0; i < l; i++ {
		numbers = append(numbers, strconv.Itoa(GenRandomInt(1, 12)))
	}
	return strings.Join(numbers, ",")
}

func CountHackedSafeCode(secret string) int {
	secretNumbers := strings.Split(secret, ",")
	hacked := 0
	for i := 0; i < len(secretNumbers); i++ {
		if secretNumbers[i] != "*" {
			hacked++
		}
	}
	return hacked
}

func HackSafeCode(code string, secret string) string {
	if code == secret {
		return secret
	}
	codeNumbers := strings.Split(code, ",")
	secretNumbers := strings.Split(secret, ",")
	hiddenNumbers := make([]int, 0, len(codeNumbers))
	for i := 0; i < len(codeNumbers); i++ {
		if codeNumbers[i] != secretNumbers[i] {
			hiddenNumbers = append(hiddenNumbers, i)
		}
	}
	intRand := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	randNumberPos := intRand.Intn(len(hiddenNumbers))
	secretNumbers[hiddenNumbers[randNumberPos]] = codeNumbers[hiddenNumbers[randNumberPos]]
	return strings.Join(secretNumbers, ",")
}

func GenRandomInt(min, max int) int {
	intRand := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	return intRand.Intn(max-min+1) + min
}
