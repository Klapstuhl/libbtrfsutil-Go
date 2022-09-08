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
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type btrfsMountpoint struct {
	path  string
	image *os.File
}

func mountBtrfs() (*btrfsMountpoint, error) {
	if !hasPrivileges() {
		return nil, errors.New("must be run as root")
	}

	path, err := os.MkdirTemp(os.TempDir(), "btrfsutil-")
	if err != nil {
		return nil, err
	}

	image, err := os.CreateTemp(os.TempDir(), "btrfsutil-")
	if err != nil {
		os.Remove(path)
		return nil, err
	}

	if err := image.Truncate(1024 * 1024 * 1024); err != nil {
		os.Remove(path)
		image.Close()
		os.Remove(image.Name())
		return nil, err
	}

	err = exec.Command("mkfs.btrfs", "-q", image.Name()).Run()
	if err != nil {
		os.Remove(path)
		image.Close()
		os.Remove(image.Name())
		return nil, err
	}

	err = exec.Command("mount", "-o", "loop,user_subvol_rm_allowed", image.Name(), path).Run()
	if err != nil {
		os.Remove(path)
		image.Close()
		os.Remove(image.Name())
		return nil, err
	}

	return &btrfsMountpoint{path, image}, nil
}

func cleanup(mp *btrfsMountpoint) {
	exec.Command("umount", "-R", mp.path).Run()
	os.Remove(mp.path)
	mp.image.Close()
	os.Remove(mp.image.Name())
}

func touch(path string) {
	now := time.Now()
	os.Chtimes(path, now, now)
}

func hasPrivileges() bool {
	return os.Geteuid() == 0
}

func superGeneration(mp *btrfsMountpoint) (uint64, error) {
	bytes := make([]byte, 8)
	_, err := mp.image.ReadAt(bytes, 65536+32+16+8+8)
	if err != nil {
		return 0, err
	}
	if string(bytes) != "_BHRfS_M" {
		return 0, fmt.Errorf("wrong Magic value: got '%s' want '_BHRfS_M'", string(bytes))
	}
	_, err = mp.image.ReadAt(bytes, 65536+32+16+8+8+8)
	return binary.LittleEndian.Uint64(bytes), err
}

func compareSubvolumeInfo(got, want *SubvolumeInfo) string {
	res := "SubvolumeInfo mismatch:"
	if got.Id != want.Id {
		res += fmt.Sprintf("\n\tid: got %d, want %d", got.Id, got.Id)
	}
	if got.ParentId != want.ParentId {
		res += fmt.Sprintf("\n\tparent_id: got %d, want %d", got.ParentId, want.ParentId)
	}
	if got.DirId != want.DirId {
		res += fmt.Sprintf("\n\tdir_id: got %d, want %d", got.DirId, want.DirId)
	}
	if got.Flags != want.Flags {
		res += fmt.Sprintf("\n\tflags: got %d, want %d", got.Flags, want.Flags)
	}
	if got.Generation != want.Generation {
		res += fmt.Sprintf("\n\tgeneration: got %d, want %d", got.Generation, want.Generation)
	}
	if got.Ctransid < want.Ctransid {
		res += fmt.Sprintf("\n\tctransid: got %d, want >= %d", got.Ctransid, want.Ctransid)
	}
	if got.Rtransid != want.Rtransid {
		res += fmt.Sprintf("\n\trtransid: got %d, want %d", got.Rtransid, want.Rtransid)
	}
	if got.Stransid != want.Stransid {
		res += fmt.Sprintf("\n\tstransid: got %d, want %d", got.Stransid, want.Stransid)
	}
	if !want.Ctime.Before(got.Ctime) {
		res += fmt.Sprintf("\n\tctime: got %s, want after %s", got.Ctime.String(), want.Ctime.String())
	}
	if !want.Otime.Before(got.Otime) {
		res += fmt.Sprintf("\n\totime: got %s, want after %s", got.Otime.String(), want.Otime.String())
	}
	if !got.Stime.Equal(want.Stime) {
		res += fmt.Sprintf("\n\tstime: got %s, want %s", got.Stime.String(), want.Stime.String())
	}
	if !got.Rtime.Equal(want.Stime) {
		res += fmt.Sprintf("\n\trtime: got %s, want %s", got.Rtime.String(), want.Rtime.String())
	}

	if len(res) == 23 {
		return ""
	}
	return res
}
