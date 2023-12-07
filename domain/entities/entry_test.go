package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEntry(t *testing.T) {
	type args struct {
		id        int64
		account   string
		direction Direction
		amount    int64
		createdAt *time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    *Entry
		wantErr bool
	}{
		{
			name: "should create a new entry - Debit",
			args: args{
				id:        1,
				account:   "account",
				direction: Debit,
				amount:    100,
				createdAt: nil,
			},
			want: &Entry{
				ID:        1,
				Account:   "account",
				Direction: Debit,
				Amount:    100,
				CreatedAt: nil,
			},
			wantErr: false,
		},
		{
			name: "should create a new entry - Credit",
			args: args{
				id:        1,
				account:   "account",
				direction: Credit,
				amount:    100,
				createdAt: nil,
			},
			want: &Entry{
				ID:        1,
				Account:   "account",
				Direction: Credit,
				Amount:    100,
				CreatedAt: nil,
			},
			wantErr: false,
		},
		{
			name: "should fail to create a new entry with invalid direction",
			args: args{
				id:        1,
				account:   "account",
				direction: "invalid",
				amount:    100,
				createdAt: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEntry(tt.args.id, tt.args.account, tt.args.direction, tt.args.amount, tt.args.createdAt)

			if !assert.Equal(t, tt.wantErr, err != nil) {
				t.Errorf("NewEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)

		})
	}
}

func TestNewDoubleEntry(t *testing.T) {
	type args struct {
		debit  Entry
		credit Entry
	}
	tests := []struct {
		name    string
		args    args
		want    *DoubleEntry
		wantErr bool
	}{
		{
			name: "should create a new double entry",
			args: args{
				debit: Entry{
					ID:        1,
					Account:   "account",
					Direction: Debit,
					Amount:    100,
					CreatedAt: nil,
				},
				credit: Entry{
					ID:        2,
					Account:   "account",
					Direction: Credit,
					Amount:    100,
					CreatedAt: nil,
				},
			},
			want: &DoubleEntry{
				Debit: Entry{
					ID:        1,
					Account:   "account",
					Direction: Debit,
					Amount:    100,
					CreatedAt: nil,
				},
				Credit: Entry{
					ID:        2,
					Account:   "account",
					Direction: Credit,
					Amount:    100,
					CreatedAt: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "should fail to create a new double entry with different amounts",
			args: args{
				debit: Entry{
					ID:        1,
					Account:   "account",
					Direction: Debit,
					Amount:    100,
					CreatedAt: nil,
				},
				credit: Entry{
					ID:        2,
					Account:   "account",
					Direction: Credit,
					Amount:    200,
					CreatedAt: nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDoubleEntry(tt.args.debit, tt.args.credit)
			if !assert.Equal(t, tt.wantErr, err != nil) {
				t.Errorf("NewDoubleEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
