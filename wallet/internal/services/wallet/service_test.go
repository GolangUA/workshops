package wallet

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_GetPaymentCapacity(t *testing.T) {
	type fields struct {
		repo Repository
	}
	type args struct {
		walletID string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := NewMockRepository(ctrl)
	repoMock.EXPECT().GetBalance("858f91c2-8bf8-4784-abd0-3181aff30001").Return(int64(0), errors.New("GetBalance error"))
	repoMock.EXPECT().GetBalance("858f91c2-8bf8-4784-abd0-3181aff30002").Return(int64(50), nil).Times(2)
	repoMock.EXPECT().GetCreditLimit("858f91c2-8bf8-4784-abd0-3181aff30002").Return(int64(0), errors.New("GetCreditLimit error"))
	repoMock.EXPECT().GetCreditLimit("858f91c2-8bf8-4784-abd0-3181aff30002").Return(int64(100), nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr error
	}{
		{
			name: "GetBalance error",
			fields: fields{
				repo: repoMock,
			},
			args: args{
				walletID: "858f91c2-8bf8-4784-abd0-3181aff30001",
			},
			want:    0,
			wantErr: errors.New("GetBalance error"),
		},
		{
			name: "GetCreditLimit error",
			fields: fields{
				repo: repoMock,
			},
			args: args{
				walletID: "858f91c2-8bf8-4784-abd0-3181aff30002",
			},
			want:    0,
			wantErr: errors.New("GetCreditLimit error"),
		},
		{
			name: "Success",
			fields: fields{
				repo: repoMock,
			},
			args: args{
				walletID: "858f91c2-8bf8-4784-abd0-3181aff30002",
			},
			want:    150,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repo: tt.fields.repo,
			}
			got, err := s.GetPaymentCapacity(tt.args.walletID)
			if tt.wantErr != nil {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}

			if got != tt.want {
				t.Errorf("GetPaymentCapacity() got = %v, want %v", got, tt.want)
			}
		})
	}
}
