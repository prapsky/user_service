package zerolog_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prapsky/user_service/common/logger/zerolog"
)

func TestNewZeroLog(t *testing.T) {
	logger := zerolog.NewZeroLog()

	assert.NotNil(t, logger)
}

func TestNewZeroLog_ErrorWithContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), zerolog.ContextKeyRequestId, "1234-abcd-5678-efgh")
	err := errors.New("some error")
	logger := zerolog.NewZeroLog()

	assert.NotNil(t, logger)
	logger.ErrorWithContext(ctx, err, "log error")
}

func TestNewZeroLog_ErrorfWithContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), zerolog.ContextKeyRequestId, "1234-abcd-5678-efgh")
	err := errors.New("some error")
	logger := zerolog.NewZeroLog()

	assert.NotNil(t, logger)
	logger.ErrorfWithContext(ctx, err, "log error %s", "with formatted string")
}

func TestNewZeroLog_WarnfWithContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), zerolog.ContextKeyRequestId, "1234-abcd-5678-efgh")
	err := errors.New("some error")
	logger := zerolog.NewZeroLog()

	assert.NotNil(t, logger)
	logger.WarnfWithContext(ctx, err, "log with warning")
}

func TestNewZeroLog_InfofWithContext(t *testing.T) {
	ctx := context.WithValue(context.TODO(), zerolog.ContextKeyRequestId, "1234-abcd-5678-efgh")
	ctx = context.WithValue(ctx, zerolog.ContextKeyEventId, "event-id")
	logger := zerolog.NewZeroLog()

	assert.NotNil(t, logger)
	logger.InfofWithContext(ctx, "log info %s", "with detail info")
}

func TestNewZeroLog_WithHandlerName(t *testing.T) {
	logger := zerolog.NewZeroLog().WithHandlerName("handlerName")

	assert.NotNil(t, logger)
	logger.Info("with handler name log")
}

func TestNewZeroLog_WithServiceName(t *testing.T) {
	logger := zerolog.NewZeroLog().WithServiceName("serviceName")

	assert.NotNil(t, logger)
	logger.Info("with service name log")
}

func TestNewZeroLog_WithRepositoryName(t *testing.T) {
	logger := zerolog.NewZeroLog().WithRepositoryName("repositoryName")

	assert.NotNil(t, logger)
	logger.Info("with repository name log")
}
