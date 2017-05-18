// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"testing"
)

func TestConfigDefaultFileSettingsDirectory(t *testing.T) {
	c1 := Config{}
	c1.SetDefaults()

	if c1.FileSettings.Directory != "./data/" {
		t.Fatal("FileSettings.Directory should default to './data/'")
	}
}

func TestConfigIsValid(t *testing.T) {
	config := Config{}
	config.SetDefaults()

	// These should be removed as we move these settings to use pointers
	config.ServiceSettings.ListenAddress = ":8065"
	config.ServiceSettings.MaximumLoginAttempts = 5
	config.TeamSettings.MaxUsersPerTeam = 100
	config.SqlSettings.DriverName = "mysql"
	config.SqlSettings.DataSource = "localhost"
	config.SqlSettings.MaxIdleConns = 10
	config.SqlSettings.MaxOpenConns = 10
	config.FileSettings.DriverName = "local"
	config.FileSettings.PreviewWidth = 100
	config.FileSettings.ProfileHeight = 100
	config.FileSettings.ProfileWidth = 100
	config.FileSettings.ThumbnailHeight = 100
	config.FileSettings.ThumbnailWidth = 100
	config.RateLimitSettings.PerSec = 100
	config.RateLimitSettings.MemoryStoreSize = 100

	if err := config.IsValid(); err != nil {
		t.Fatal("Default config should be valid", err)
	}

	// tests for other fields should be added

	*config.DataRetentionSettings.RunAtHour = 0
	if err := config.IsValid(); err != nil {
		t.Fatal(err)
	}

	*config.DataRetentionSettings.RunAtHour = 23
	if err := config.IsValid(); err != nil {
		t.Fatal(err)
	}

	*config.DataRetentionSettings.RunAtHour = -1
	if err := config.IsValid(); err == nil {
		t.Fatal("shouldn't allow negative data retention schedule")
	}

	*config.DataRetentionSettings.RunAtHour = 24
	if err := config.IsValid(); err == nil {
		t.Fatal("shouldn't allow data retention hour greater than 23")
	}

	*config.DataRetentionSettings.RunAtHour = 2

	*config.DataRetentionSettings.PostRetentionPeriod = -1
	if err := config.IsValid(); err == nil {
		t.Fatal("shouldn't allow negative post retention period")
	}

	*config.DataRetentionSettings.PostRetentionPeriod = 0
	*config.DataRetentionSettings.FileRetentionPeriod = -1
	if err := config.IsValid(); err == nil {
		t.Fatal("shouldn't allow negative file retention period")
	}

	*config.DataRetentionSettings.PostRetentionPeriod = 0
	*config.DataRetentionSettings.FileRetentionPeriod = 10
	if err := config.IsValid(); err != nil {
		t.Fatal("should allow finite file retention with indefinite post retention", err)
	}

	*config.DataRetentionSettings.PostRetentionPeriod = 10
	*config.DataRetentionSettings.FileRetentionPeriod = 0
	if err := config.IsValid(); err == nil {
		t.Fatal("shouldn't allow finite post retention with indefinite file retention")
	}

	*config.DataRetentionSettings.PostRetentionPeriod = 10
	*config.DataRetentionSettings.FileRetentionPeriod = 10
	if err := config.IsValid(); err != nil {
		t.Fatal("should allow same file retention window as post retention window", err)
	}

	*config.DataRetentionSettings.PostRetentionPeriod = 10
	*config.DataRetentionSettings.FileRetentionPeriod = 5
	if err := config.IsValid(); err != nil {
		t.Fatal("should allow shorter file retention window than post retention window", err)
	}

	*config.DataRetentionSettings.PostRetentionPeriod = 5
	*config.DataRetentionSettings.FileRetentionPeriod = 10
	if err := config.IsValid(); err == nil {
		t.Fatal("shouldn't allow shorter post retention window than file retention window")
	}

	*config.DataRetentionSettings.PostRetentionPeriod = 00
	*config.DataRetentionSettings.FileRetentionPeriod = 0
}
