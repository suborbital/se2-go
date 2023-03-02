package se2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suborbital/se2-go"
)

func TestNewClient2(t *testing.T) {
	type args struct {
		host    se2.ServerURL
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
				host: se2.HostProduction,
				ak:   "eyJrZXkiOjQwNywic2VjcmV0IjoiZWsvNFV3VTBnZ2VHUjdQanF1MmlyaWJacGR1MXZvcWNhMXl3eDE3aWhpTT0ifQ==",
			},
			wantErr: assert.NoError,
		},
		{
			name: "correct for staging, valid key, no options",
			args: args{
				host: se2.HostStaging,
				ak:   "eyJrZXkiOjQwNywic2VjcmV0IjoiZWsvNFV3VTBnZ2VHUjdQanF1MmlyaWJacGR1MXZvcWNhMXl3eDE3aWhpTT0ifQ==",
			},
			wantErr: assert.NoError,
		},
		{
			name: "correct for prod, empty key, no options",
			args: args{
				host: se2.HostProduction,
				ak:   "",
			},
			wantErr: assert.Error,
		},
		{
			name: "correct for prod, decode error key, no options",
			args: args{
				host: se2.HostProduction,
				ak:   "eJrZXkiOjQwNywic2VjcmV0IjoiZWs3VTBnZ2VHUjdQanF1MmlyaWJacGR1MXZvcWNheDE3aWhpTT0ifQ",
			},
			wantErr: assert.Error,
		},
		{
			name: "correct for prod, short key, no options",
			args: args{
				host: se2.HostProduction,
				ak:   "eJrZXkiOjQwNywic2VjcmV0IjoiZWs3VTB",
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := se2.NewClient2(tt.args.host, tt.args.ak, tt.args.options...)

			tt.wantErr(t, err)
		})
	}
}
