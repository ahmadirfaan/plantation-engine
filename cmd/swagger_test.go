package main

import (
	"testing"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/stretchr/testify/assert"
)

func TestGetSwagger_Success(t *testing.T) {
	swagger, err := generated.GetSwagger()
	assert.NotNil(t, swagger)
	assert.Nil(t, err)

	assert.Equal(t, "3.0.0", swagger.OpenAPI)
}
