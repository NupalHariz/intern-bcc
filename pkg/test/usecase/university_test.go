package test

import (
	"intern-bcc/domain"
	"intern-bcc/internal/usecase"
	test "intern-bcc/pkg/test/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUniversity(t *testing.T) {
	
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := test.NewMockIUniversityRepository(ctrl)
	// fmt.Println(c)

	// type args struct {
	// 	input domain.Universities
	// }

	mockParamUniversity := domain.Universities{
		University: "Brawijaya",
	}

	type args struct {
		universityRequest domain.Universities
	}

	type mockObjects struct{
		universityRepo *test.MockIUniversityRepository
	}

	mockObject := mockObjects{
		universityRepo: c,
	}

	mockArgs := args{
		universityRequest: mockParamUniversity,
	}

	tests := []struct {
		name string
		args
		beforeTest func(mockObject mockObjects)
		wantErr    bool
	}{
		{
			name: "failed create university",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.universityRepo.EXPECT().CreateUniversity(&mockParamUniversity).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.universityRepo.EXPECT().CreateUniversity(&mockParamUniversity).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//mockUniversityRepo := test.NewMockIUniversityRepository(ctrl)
			w := usecase.NewUniversityUsecase(c)

			// if tt.beforeTest != nil {
			tt.beforeTest(mockObject)
			// }

			err := w.CreateUniversity(tt.args.universityRequest)
			if err != nil && !tt.wantErr {
				t.Errorf("error = %v", err)
			}
		})
	}
}

