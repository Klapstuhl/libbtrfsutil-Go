/*
 * Copyright (C) 2022 Jan-Oliver Rettig
 *
 * This file is part of libbtrfsutil-go.
 *
 * libbtrfsutil-go is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * libbtrfsutil-go is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with libbtrfsutil-go.  If not, see <http://www.gnu.org/licenses/>.
 */

package btrfsutil

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestSync(t *testing.T) {
	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	old_generation, err := superGeneration(mountpoint)
	if err != nil {
		t.Fatal(err)
	}

	touch(mountpoint.path)
	if err := Sync(mountpoint.path); err != nil {
		t.Errorf("Sync() error = %v", err)
	}

	new_generation, err := superGeneration(mountpoint)
	if err != nil {
		t.Fatal(err)
	}
	if new_generation <= old_generation {
		t.Errorf("Sync failed: New Generation '%d' <= Old Generation '%d'", new_generation, old_generation)
	}
}

func TestStartSync(t *testing.T) {
	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	old_generation, err := superGeneration(mountpoint)
	if err != nil {
		t.Fatal(err)
	}
	touch(mountpoint.path)

	transid, err := StartSync(mountpoint.path)
	if err != nil {
		t.Errorf("StartSync() error = %v", err)
	}
	if transid <= old_generation {
		t.Errorf("Sync failed: StartSync() = %v <= %v", transid, old_generation)
	}
}

func TestWaitSync(t *testing.T) {
	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	old_generation, err := superGeneration(mountpoint)
	if err != nil {
		t.Fatal(err)
	}
	touch(mountpoint.path)

	transid, err := StartSync(mountpoint.path)
	if err != nil {
		t.Errorf("StartSync() error = %v", err)
	}

	if err := WaitSync(mountpoint.path, transid); err != nil {
		t.Errorf("WaitSync() error = %v", err)
	}
	new_generation, err := superGeneration(mountpoint)
	if err != nil {
		t.Fatal(err)
	}
	if new_generation <= old_generation {
		t.Errorf("Sync failed. New Generation '%d' <= Old Generation '%d'", new_generation, old_generation)
	}
}

func TestIsSubvolume(t *testing.T) {
	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	foo := filepath.Join(mountpoint.path, "foo")
	os.Mkdir(foo, 0770)

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"<FS_TREE>", args{path: mountpoint.path}, true, false},
		{"foo", args{path: foo}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsSubvolume(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsSubvolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsSubvolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubvolumeId(t *testing.T) {
	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	foo := filepath.Join(mountpoint.path, "foo")
	os.Mkdir(foo, 0770)

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"<FS_TREE>", args{path: mountpoint.path}, 5, false},
		{"foo", args{path: foo}, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SubvolumeId(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubvolumeId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SubvolumeId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateSubvolume(t *testing.T) {
	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	foo := filepath.Join(mountpoint.path, "foo")
	os.Mkdir(foo, 0770)

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"<FS_TREE>", args{path: ""}, true},
		{"subvol1", args{path: "subvol1"}, false},
		{"subvol2", args{path: "foo/subvol2"}, false},
		{"subvol3", args{path: "foo/subvol2/subvol3"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subvol := filepath.Join(mountpoint.path, tt.args.path)
			if err := CreateSubvolume(subvol); (err != nil) != tt.wantErr {
				t.Errorf("CreateSubvolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubvolumePath(t *testing.T) {
	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	foo := filepath.Join(mountpoint.path, "foo")
	os.Mkdir(foo, 0770)

	if CreateSubvolume(filepath.Join(mountpoint.path, "foo/subvol1")) != nil {
		t.Error("Failed to create subvolumes")
	}

	type args struct {
		path string
		id   uint64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"foo", args{path: "foo", id: 0}, "", true},

		{"<FS_TREE>-1", args{path: "", id: 0}, "", false},
		{"<FS_TREE>-2", args{path: "", id: 5}, "", false},

		{"subvol1-1", args{path: "", id: 256}, "foo/subvol1", false},
		{"subvol1-2", args{path: "foo/subvol1", id: 0}, "foo/subvol1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subvol := filepath.Join(mountpoint.path, tt.args.path)
			got, err := SubvolumePath(subvol, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubvolumePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SubvolumePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSubvolumeInfo(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	if CreateSubvolume(filepath.Join(mountpoint.path, "subvol1")) != nil {
		t.Error("Failed to create subvolumes")
	}

	type args struct {
		path string
		id   uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *SubvolumeInfo
		wantErr bool
	}{
		{
			"<FS_TREE>",
			args{path: "", id: 0},
			&SubvolumeInfo{
				Id:         5,
				ParentId:   0,
				DirId:      0,
				Flags:      0,
				Generation: 7,
				Ctransid:   0,
				Otransid:   0,
				Stransid:   0,
				Rtransid:   0,
				Ctime:      now,
				Otime:      time.Unix(0, 0),
				Stime:      time.Unix(0, 0),
				Rtime:      time.Unix(0, 0),
			},
			false,
		},
		{
			"subvol1",
			args{path: "", id: 256},
			&SubvolumeInfo{
				Id:         256,
				ParentId:   5,
				DirId:      256,
				Flags:      0,
				Generation: 7,
				Ctransid:   0,
				Otransid:   0,
				Stransid:   0,
				Rtransid:   0,
				Ctime:      now,
				Otime:      now,
				Stime:      time.Unix(0, 0),
				Rtime:      time.Unix(0, 0),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subvol := filepath.Join(mountpoint.path, tt.args.path)
			got, err := GetSubvolumeInfo(subvol, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubvolumeInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if res := compareSubvolumeInfo(got, tt.want); res != "" {
				t.Error(res)
			}
		})
	}
}

func TestSubvolumeReadOnly(t *testing.T) {
	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	foo := filepath.Join(mountpoint.path, "foo")
	os.Mkdir(foo, 0770)

	type args struct {
		path      string
		read_only bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"<FS_TREE> ro", args{path: "", read_only: true}, true, false},
		{"foo ro", args{path: "foo", read_only: true}, false, true},
		{"<FS_TREE> rw", args{path: "", read_only: false}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subvol := filepath.Join(mountpoint.path, tt.args.path)
			if err := SetSubvolumeReadOnly(subvol, tt.args.read_only); (err != nil) != tt.wantErr {
				t.Errorf("SetSubvolumeReadOnly() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := GetSubvolumeReadOnly(subvol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubvolumeReadOnly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSubvolumeReadOnly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultSubvolume(t *testing.T) {
	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	foo := filepath.Join(mountpoint.path, "foo")
	os.Mkdir(foo, 0770)

	if CreateSubvolume(filepath.Join(mountpoint.path, "subvol1")) != nil {
		t.Error("Failed to create subvolumes")
	}

	type args struct {
		path string
		id   uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{"<FS_TREE>", args{path: "", id: 0}, 5, false},
		{"foo", args{path: "foo", id: 0}, 5, true},
		{"subvol1-1", args{path: "subvol1", id: 0}, 256, false},
		{"subvol1-2", args{path: "subvol1", id: 5}, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subvol := filepath.Join(mountpoint.path, tt.args.path)

			if err := SetDefaultSubvolume(subvol, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultSubvolume() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := GetDefaultSubvolume(subvol)
			if err != nil {
				t.Errorf("GetDefaultSubvolume() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("GetDefaultSubvolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateSnapshot(t *testing.T) {
	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	foo := filepath.Join(mountpoint.path, "foo")
	os.Mkdir(foo, 0770)

	if CreateSubvolume(filepath.Join(mountpoint.path, "subvol1")) != nil {
		t.Error("Failed to create subvolumes")
	}
	if CreateSubvolume(filepath.Join(mountpoint.path, "subvol1/subvol2")) != nil {
		t.Error("Failed to create subvolumes")
	}

	type args struct {
		source    string
		path      string
		recursive bool
		read_only bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"foo", args{source: "foo", path: "foo1", recursive: false, read_only: false}, false},
		{"snap1", args{source: "subvol1", path: "foo/snap1", recursive: false, read_only: false}, false},
		{"snap2", args{source: "subvol1", path: "foo/snap2", recursive: false, read_only: true}, false},
		{"snap3", args{source: "subvol1", path: "foo/snap3", recursive: true, read_only: false}, false},
		{"snap4", args{source: "subvol1", path: "foo/snap4", recursive: true, read_only: true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := filepath.Join(mountpoint.path, tt.args.source)
			path := filepath.Join(mountpoint.path, tt.args.path)
			if err := CreateSnapshot(source, path, tt.args.recursive, tt.args.read_only); (err != nil) != tt.wantErr {
				t.Errorf("CreateSnapshot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteSubvolume(t *testing.T) {
	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	if CreateSubvolume(filepath.Join(mountpoint.path, "subvol1")) != nil {
		t.Error("Failed to create subvolumes")
	}
	if CreateSubvolume(filepath.Join(mountpoint.path, "subvol1/subvol2")) != nil {
		t.Error("Failed to create subvolumes")
	}
	if CreateSubvolume(filepath.Join(mountpoint.path, "subvol1/subvol3")) != nil {
		t.Error("Failed to create subvolumes")
	}

	type args struct {
		path      string
		recursive bool
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"subvol3", args{path: "subvol1/subvol3", recursive: false}, false},
		{"subvol1", args{path: "subvol1", recursive: false}, true},
		{"subvol1-r", args{path: "subvol1", recursive: true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subvol := filepath.Join(mountpoint.path, tt.args.path)
			if err := DeleteSubvolume(subvol, tt.args.recursive); (err != nil) != tt.wantErr {
				t.Errorf("DeleteSubvolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeletedSubvolumes(t *testing.T) {
	if !hasPrivileges() {
		t.Skipf("must be run as root")
	}

	mountpoint, err := mountBtrfs()
	if err != nil {
		t.Skip(err)
	}
	defer cleanup(mountpoint)

	test := struct {
		subvols []string
		want    []uint64
		wantErr bool
	}{
		[]string{"subvol1", "subvol2"},
		[]uint64{256, 257},
		false,
	}

	t.Run("", func(t *testing.T) {
		for _, subvol := range test.subvols {
			subvol = filepath.Join(mountpoint.path, subvol)
			if CreateSubvolume(subvol) != nil {
				t.Error("Failed to create subvolumes")
			}
			if DeleteSubvolume(subvol, false) != nil {
				t.Error("Failed to delete subvolumes")
			}
		}

		got, err := DeletedSubvolumes(mountpoint.path)
		if (err != nil) != test.wantErr {
			t.Errorf("DeletedSubvolumes() error = %v, wantErr %v", err, test.wantErr)
			return
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("DeletedSubvolumes() = %v, want %v", got, test.want)
		}
	})
}
