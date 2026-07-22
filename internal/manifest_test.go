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
			name: "Test getDefaultScaffold",
			want: map[string]ScaffoldRepo{
				"drupal-9": {
					Name:             "drupal-9",
					GitRepo:          "https://github.com/lagoon-examples/drupal9-full.git",
					Branch:           "scaffold",
					ShortDescription: "Pulls and sets up a new Lagoon ready Drupal 9",
					Description:      "Pulls and sets up a new Lagoon ready Drupal 9",
				},
				"rails-init": {
					Name:             "rails-init",
					GitRepo:          "https://github.com/CGoodwin90/lagoon-rails-dir.git",
					Branch:           "main",
					ShortDescription: "Will add a minimal set of files to an existing Rails 7 installation",
					Description:      "Will add a minimal set of files to an existing Rails 7 installation",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDefaultScaffold()
			if (err != nil) != tt.wantErr {
				t.Errorf("getDefaultScaffold() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDefaultScaffold() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetScaffolds(t *testing.T) {
	type args struct {
		localmanifest string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]ScaffoldRepo
		wantErr bool
	}{
		{
			name: "Test 1 - Passing manifest file manually",
			args: args{
				localmanifest: "./testassets/manifest_test_1.yml",
			},
			want: map[string]ScaffoldRepo{
				"test1": {
					Name:             "test1",
					GitRepo:          "https://github.com/lagoon-examples/test1.git",
					Branch:           "test1_branch",
					Description:      "test1_description",
					ShortDescription: "test1_shortDescription",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetScaffolds(tt.args.localmanifest)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetScaffolds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetScaffolds() got = %v, want %v", got, tt.want)
			}
		})
	}
}
