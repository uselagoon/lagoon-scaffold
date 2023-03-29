package internal

import (
	"reflect"
	"testing"
)

func Test_resolveScaffolds(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]ScaffoldRepo
		wantErr bool
	}{
		{
			name: "Test resolveScaffolds",
			want: map[string]ScaffoldRepo{
				"laravel-init": {
					Name:             "laravel-init",
					GitRepo:          "https://github.com/bomoko/lagoon-laravel-dir.git",
					Branch:           "main",
					ShortDescription: "Will add a minimal set of files to an existing Laravel 10 installation",
					Description:      "Will add a minimal set of files to an existing Laravel 10 installation",
				},
				"drupal-9": {
					Name:             "drupal-9",
					GitRepo:          "https://github.com/lagoon-examples/drupal9-full.git",
					Branch:           "scaffold",
					ShortDescription: "Pulls and sets up a new Lagoon ready Drupal 9",
					Description:      "Pulls and sets up a new Lagoon ready Drupal 9",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveScaffolds()
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveScaffolds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveScaffolds() got = %v, want %v", got, tt.want)
			}
		})
	}
}
