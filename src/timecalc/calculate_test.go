package timecalc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tpoulsen/pcc-timebot/src/helpers"
)

func Test_buildTimeDelta(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "standard time format AM",
			timeStr: "9:00am",
			want:    9 * time.Hour,
			wantErr: false,
		},
		{
			name:    "standard time format PM",
			timeStr: "1:30pm",
			want:    13*time.Hour + 30*time.Minute,
			wantErr: false,
		},
		{
			name:    "short format AM",
			timeStr: "9am",
			want:    9 * time.Hour,
			wantErr: false,
		},
		{
			name:    "short format PM",
			timeStr: "2pm",
			want:    14 * time.Hour,
			wantErr: false,
		},
		{
			name:    "midnight",
			timeStr: "12:00am",
			want:    0,
			wantErr: false,
		},
		{
			name:    "noon",
			timeStr: "12:00pm",
			want:    12 * time.Hour,
			wantErr: false,
		},
		{
			name:    "invalid format",
			timeStr: "11:000am",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid hours",
			timeStr: "13:00pm",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid minutes",
			timeStr: "9:60am",
			want:    0,
			wantErr: true,
		},
		{
			name:    "missing meridiem",
			timeStr: "9:00",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDuration(tt.timeStr)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCalculateTime(t *testing.T) {
	type args struct {
		start string
		end   string
		less  string
		more  string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "standard workday with lunch",
			args: args{
				start: "9:00am",
				end:   "5:00pm",
				less:  "1",
				more:  "0",
			},
			want:    7.0,
			wantErr: false,
		},
		{
			name: "short format with half hour lunch",
			args: args{
				start: "8am",
				end:   "4pm",
				less:  "0.5",
				more:  "0",
			},
			want:    7.5,
			wantErr: false,
		},
		{
			name: "with extra time and lunch",
			args: args{
				start: "7:30am",
				end:   "3:30pm",
				less:  "1",
				more:  "0.5",
			},
			want:    7.5,
			wantErr: false,
		},
		{
			name: "midnight to noon",
			args: args{
				start: "12:00am",
				end:   "12:00pm",
				less:  "0",
				more:  "0",
			},
			want:    12.0,
			wantErr: false,
		},
		{
			name: "invalid end before start",
			args: args{
				start: "5:00pm",
				end:   "9:00am",
				less:  "0",
				more:  "0",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid lunch format",
			args: args{
				start: "9:00am",
				end:   "5:00pm",
				less:  "invalid",
				more:  "0",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid extra time format",
			args: args{
				start: "9:00am",
				end:   "5:00pm",
				less:  "1",
				more:  "invalid",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateTime(tt.args.start, tt.args.end, tt.args.less, tt.args.more)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0.0, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_round(t *testing.T) {
	type args struct {
		num    float64
		places int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "round to 2 decimal places",
			args: args{
				num:    1.2345,
				places: 2,
			},
			want: 1.23,
		},
		{
			name: "round up at .5",
			args: args{
				num:    1.235,
				places: 2,
			},
			want: 1.24,
		},
		{
			name: "round to whole number",
			args: args{
				num:    1.6,
				places: 0,
			},
			want: 2,
		},
		{
			name: "round negative number",
			args: args{
				num:    -1.234,
				places: 2,
			},
			want: -1.23,
		},
		{
			name: "round to 3 decimal places",
			args: args{
				num:    1.2345,
				places: 3,
			},
			want: 1.235,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.Round(tt.args.num, tt.args.places)
			assert.Equal(t, got, tt.want)
		})
	}
}
