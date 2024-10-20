package retoolsdk_test

import (
	"github.com/stretchr/testify/assert"
	retool "github.com/thoughtgears/retoolsdk"
	"testing"
)

func TestUpdateOperations_Validate(t *testing.T) {
	validOperations := []retool.UpdateOperations{
		{Op: retool.OpAdd, Path: "/email", Value: "new.email@example.com"},
		{Op: retool.OpReplace, Path: "/name", Value: "Jane Doe"},
		{Op: retool.OpRemove, Path: "/phone"},
	}

	for _, op := range validOperations {
		err := op.Validate()
		assert.NoError(t, err, "Valid operation should not return an error")
	}

	invalidOperations := []retool.UpdateOperations{
		{Op: "invalid_op", Path: "/email", Value: "new.email@example.com"},
		{Op: retool.OpAdd, Path: "", Value: "missing path"},
		{Op: retool.OpReplace, Path: "/email", Value: ""},
	}

	for _, op := range invalidOperations {
		err := op.Validate()
		assert.Error(t, err, "Invalid operation should return an error")
	}

	emptyOperation := retool.UpdateOperations{}
	err := emptyOperation.Validate()
	assert.Error(t, err, "Empty operation should return an error")
}
