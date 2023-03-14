package se2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suborbital/se2-go"
)

func TestNewClient2(t *testing.T) {
	type args struct {
		mode    se2.ServerMode
		ak      string
		options []se2.ClientOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "correct for production, valid key, no options",
			args: args{
				mode: se2.ModeProduction,
				ak:   "eyJrZXkiOjQwNywic2VjcmV0IjoiZWsvNFV3VTBnZ2VHUjdQanF1MmlyaWJacGR1MXZvcWNhMXl3eDE3aWhpTT0ifQ==",
			},
			wantErr: assert.NoError,
		},
		{
			name: "correct for staging, valid key, no options",
			args: args{
				mode: se2.ModeStaging,
				ak:   "eyJrZXkiOjQwNywic2VjcmV0IjoiZWsvNFV3VTBnZ2VHUjdQanF1MmlyaWJacGR1MXZvcWNhMXl3eDE3aWhpTT0ifQ==",
			},
			wantErr: assert.NoError,
		},
		{
			name: "correct for prod, empty key, no options",
			args: args{
				mode: se2.ModeProduction,
				ak:   "",
			},
			wantErr: assert.Error,
		},
		{
			name: "correct for prod, decode error key, no options",
			args: args{
				mode: se2.ModeProduction,
				ak:   "eJrZXkiOjQwNywic2VjcmV0IjoiZWs3VTBnZ2VHUjdQanF1MmlyaWJacGR1MXZvcWNheDE3aWhpTT0ifQ",
			},
			wantErr: assert.Error,
		},
		{
			name: "correct for prod, short key, no options",
			args: args{
				mode: se2.ModeProduction,
				ak:   "eJrZXkiOjQwNywic2VjcmV0IjoiZWs3VTB",
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := se2.NewClient(tt.args.mode, tt.args.ak, tt.args.options...)

			tt.wantErr(t, err)
		})
	}
}
