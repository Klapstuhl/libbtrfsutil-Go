/*
 * Copyright (C) 2022 Jana Marlou Rettig
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
)

func TestCreateSubvolumeIterator(t *testing.T) {
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
		path       string
		top        uint64
		post_order bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"<FS_TREE>", args{path: mountpoint.path, top: 0, post_order: false}, false},
		{"foo", args{path: foo, top: 0, post_order: false}, true},
		{"TOP=256", args{path: mountpoint.path, top: 256, post_order: false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iter, err := CreateSubvolumeIterator(tt.args.path, tt.args.top, tt.args.post_order)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSubvolumeIterator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer iter.Destroy()
		})
	}
}

func TestCreateSubvolumeInfoIterator(t *testing.T) {
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
		path       string
		top        uint64
		post_order bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"<FS_TREE>", args{path: mountpoint.path, top: 0, post_order: false}, false},
		{"foo", args{path: foo, top: 0, post_order: false}, true},
		{"TOP=256", args{path: mountpoint.path, top: 256, post_order: false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iter, err := CreateSubvolumeInfoIterator(tt.args.path, tt.args.top, tt.args.post_order)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSubvolumeInfoIterator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer iter.Destroy()
		})
	}
}

func TestSubvolumeIterator(t *testing.T) {
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
		path       string
		top        uint64
		post_order bool
	}

	type iterNext struct {
		path string
		id   uint64
	}

	tests := []struct {
		name    string
		args    args
		want    []iterNext
		wantErr bool
	}{
		{
			"top=0 pre-order",
			args{path: mountpoint.path, top: 0, post_order: false},
			[]iterNext{
				{"subvol1", 256},
				{"subvol1/subvol2", 257},
				{"subvol1/subvol3", 258},
			},
			false,
		},
		{
			"top=0 post-order",
			args{path: mountpoint.path, top: 0, post_order: true},
			[]iterNext{
				{"subvol1/subvol2", 257},
				{"subvol1/subvol3", 258},
				{"subvol1", 256},
			},
			false,
		},
		{
			"top=256 pre-order",
			args{path: mountpoint.path, top: 256, post_order: false},
			[]iterNext{
				{"subvol2", 257},
				{"subvol3", 258},
			},
			false,
		},
	}
	t.Run("SubvolumeIterator", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				iter, err := CreateSubvolumeIterator(tt.args.path, tt.args.top, tt.args.post_order)
				if (err != nil) != tt.wantErr {
					t.Errorf("CreateSubvolumeIterator() error = %v, wantErr %v", err, tt.wantErr)
				}
				defer iter.Destroy()

				var got []iterNext
				for iter.HasNext() {
					result, err := iter.GetNext()
					if (err != nil) != tt.wantErr {
						t.Errorf("SubvolumeIterator.GetNext() error = %v, wantErr %v", err, tt.wantErr)
					}

					got = append(got, iterNext{result.Path, result.Id})
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("\n\tgot  %v\n\twant %v", got, tt.want)
				}
			})
		}
	})
	t.Run("SubvolumeInfoIterator", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				iter, err := CreateSubvolumeInfoIterator(tt.args.path, tt.args.top, tt.args.post_order)
				if (err != nil) != tt.wantErr {
					t.Errorf("CreateSubvolumeInfoIterator() error = %v, wantErr %v", err, tt.wantErr)
				}
				defer iter.Destroy()

				var got []iterNext
				for iter.HasNext(){
					result, err := iter.GetNext()
					if err == ErrStopIteration {
						break
					}
					if (err != nil) != tt.wantErr {
						t.Errorf("SubvolumeInfoIterator.GetNext() error = %v, wantErr %v", err, tt.wantErr)
					}

					got = append(got, iterNext{result.Path, result.Info.Id})
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("\n\tgot  %v\n\twant %v", got, tt.want)
				}
			})
		}
	})
}
