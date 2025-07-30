package domain

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError_Add(t *testing.T) {
	err := NewError("[error.test.root]").SetCode("rootCode")
	err1 := Error{Code: "subCode0", Source: "[error.test.sub0]"}
	err2 := errors.New("some error 0")
	err3 := NewError("[error.test.sub1]").SetCode("subCode1")

	err = err.Add(err1).Add(err2).Add(err3)

	assert.Error(t, err)
	assert.Error(t, err1)
	assert.Error(t, err2)
	assert.Error(t, err3)

	assert.Equal(t, err.Code, "rootCode")
	assert.Equal(t, err.Source, "[error.test.root]")

	assert.Equal(t, len(err.Errors), 2)
	assert.Equal(t, len(err.Native), 1)

	for i, e := range err.Errors {
		assert.Error(t, e)
		assert.Equal(t, e.Code, fmt.Sprintf("subCode%d", i))
		assert.Equal(t, e.Source, fmt.Sprintf("[error.test.sub%d]", i))
	}

	for i, e := range err.Native {
		assert.Error(t, e)
		assert.Equal(t, e.Error(), fmt.Sprintf("some error %d", i))
	}
}

func TestError_Error(t *testing.T) {
	err := &Error{Code: "someCode"}

	assert.Equal(t, err.Error(), ":someCode: ")
}

func TestError_SetCode(t *testing.T) {
	err := &Error{}
	assert.Equal(t, err.Code, "")

	err = err.SetCode("someCode")
	assert.Equal(t, err.Code, "someCode")
}

func TestError_SetSource(t *testing.T) {
	err := &Error{}
	assert.Equal(t, err.Source, "")

	err = err.SetSource("[error.test.root]")
	assert.Equal(t, err.Source, "[error.test.root]")
}
