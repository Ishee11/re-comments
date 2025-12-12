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
	defer ctrl.Finish()

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
				require.ErrorIs(t, err, errInternalServErr)
				require.Nil(t, got)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			got, err := useCase.CreateComment(context.Background(), 1, 1, "test")
			tc.want(t, got, err)
		})
	}
}

func TestUseCase_UpdateComment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := usecase_test.NewMockCommentRepository(ctrl)
	useCase := commentUseCase(repo)

	tests := []struct {
		name string
		mock func()
		want func(t *testing.T, err error)
	}{
		{
			name: "success",
			mock: func() {
				// UseCase.UpdateComment сначала вызывает GetCommentByID
				existing := &entity.Comment{
					ID:       1,
					UserID:   1,
					EntityID: 1,
					Text:     "old",
				}
				repo.EXPECT().GetCommentByID(int64(1)).Return(existing, nil)
				repo.EXPECT().
					UpdateComment(gomock.AssignableToTypeOf(&entity.Comment{})).
					Return(nil)
			},
			want: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "get comment error",
			mock: func() {
				repo.EXPECT().GetCommentByID(int64(1)).Return(nil, errInternalServErr)
			},
			want: func(t *testing.T, err error) {
				require.ErrorIs(t, err, errInternalServErr)
			},
		},
		{
			name: "update repo error",
			mock: func() {
				existing := &entity.Comment{
					ID:       1,
					UserID:   1,
					EntityID: 1,
					Text:     "old",
				}
				repo.EXPECT().GetCommentByID(int64(1)).Return(existing, nil)
				repo.EXPECT().
					UpdateComment(gomock.AssignableToTypeOf(&entity.Comment{})).
					Return(errInternalServErr)
			},
			want: func(t *testing.T, err error) {
				require.ErrorIs(t, err, errInternalServErr)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			err := useCase.UpdateComment(context.Background(), 1, 1, "new text")
			tc.want(t, err)
		})
	}
}
