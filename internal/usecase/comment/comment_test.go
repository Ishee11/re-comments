package comment

import (
	"context"
	"errors"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	usecase_test "github.com/evrone/go-clean-template/internal/usecase/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errInternalServErr = errors.New("internal server error")

type testCase struct {
	name string
	mock func()
	want func(t *testing.T, got *entity.Comment, err error)
}

func commentUseCase(repo *usecase_test.MockCommentRepository) *UseCase {
	// хелпер только создаёт usecase на основе уже созданного мока
	return New(repo)
}

func TestUseCase_CreateComment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	repo := usecase_test.NewMockCommentRepository(ctrl)
	useCase := commentUseCase(repo)

	tests := []testCase{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().
					CreateComment(gomock.AssignableToTypeOf(&entity.Comment{})).
					Return(nil)
			},
			want: func(t *testing.T, got *entity.Comment, err error) {
				t.Helper()
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, int64(1), got.UserID)
				require.Equal(t, int64(1), got.EntityID)
				require.Equal(t, "test", got.Text)
			},
		},
		{
			name: "repo error",
			mock: func() {
				repo.EXPECT().
					CreateComment(gomock.AssignableToTypeOf(&entity.Comment{})).
					Return(errInternalServErr)
			},
			want: func(t *testing.T, got *entity.Comment, err error) {
				t.Helper()
				require.ErrorIs(t, err, errInternalServErr)
				require.Nil(t, got)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.mock()
			got, err := useCase.CreateComment(context.Background(), 1, 1, "test")
			tc.want(t, got, err)
		})
	}
}

func runUpdateCommentTest(t *testing.T, name string, mock func(), want func(t *testing.T, err error)) {
	t.Helper() // сразу вызываем, это хелпер

	t.Run(name, func(t *testing.T) {
		t.Parallel() // под-тесты будут выполняться параллельно

		ctrl := gomock.NewController(t)
		t.Cleanup(func() { ctrl.Finish() }) // безопасно для параллельных тестов

		repo := usecase_test.NewMockCommentRepository(ctrl)
		useCase := commentUseCase(repo)

		mock()
		err := useCase.UpdateComment(context.Background(), 1, 1, "new text")
		want(t, err)
	})
}

func TestUseCase_UpdateComment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mock func()
		want func(t *testing.T, err error)
	}{
		{
			name: "success",
			mock: func() {
				existing := &entity.Comment{ID: 1, UserID: 1, EntityID: 1, Text: "old"}
				repo := usecase_test.NewMockCommentRepository(nil)
				repo.EXPECT().GetCommentByID(int64(1)).Return(existing, nil)
				repo.EXPECT().UpdateComment(gomock.AssignableToTypeOf(&entity.Comment{})).Return(nil)
			},
			want: func(t *testing.T, err error) {
				t.Helper()
				require.NoError(t, err)
			},
		},
		// остальные кейсы
	}

	for _, tc := range tests {
		runUpdateCommentTest(t, tc.name, tc.mock, tc.want)
	}
}
