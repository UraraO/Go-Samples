package example

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination=mock_db.go -package=example example DB

func TestGetUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockDB(ctrl)

	ctx := context.Background()
	userID := 1
	expectedName := "John Doe"

	// 设置模拟期望
	mockDB.EXPECT().GetUser(ctx, userID).Return(expectedName, nil)

	// 调用被测试的函数
	name, err := GetUserName(ctx, mockDB, userID)

	// 断言结果
	assert.NoError(t, err)
	assert.Equal(t, expectedName, name)
}
