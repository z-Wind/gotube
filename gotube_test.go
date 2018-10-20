package gotube

import "testing"

func TestYoutube_GetVideo(t *testing.T) {
	tests := []struct {
		name    string
		y       *Youtube
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", &Youtube{VideoURL: "https://www.youtube.com/watch?v=5yAU52qfYuU"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.y.GetVideo(); (err != nil) != tt.wantErr {
				t.Errorf("Youtube.GetVideo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
