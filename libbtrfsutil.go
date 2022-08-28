/*
 * Copyright (C) 2022 Jan-Oliver Rettig
 *
 * This file is part of libbtrfsutil-Go.
 *
 * libbtrfsutil-Go is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * libbtrfsutil-Go is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with libbtrfsutil-Go.  If not, see <http://www.gnu.org/licenses/>.
 */

package libbtrfsutil

// #cgo LDFLAGS: -lbtrfsutil
// #include <stdlib.h>
// #include <btrfsutil.h>
import "C"
import (
	"errors"
	"time"
	"unsafe"
)

type subvolumeInfo struct {
	id            uint64
	parent_id     uint64
	dir_id        uint64
	flags         uint64
	uuid          []uint8
	parent_uuid   []uint8
	received_uuid []uint8
	generation    uint64
	ctransid      uint64
	otransid      uint64
	stransid      uint64
	rtransid      uint64
	ctime         time.Time
	otime         time.Time
	stime         time.Time
	rtime         time.Time
}

func newSubvolumeInfo(info *C.struct_btrfs_util_subvolume_info) subvolumeInfo {
	subvol := subvolumeInfo{
		id:            uint64(info.id),
		parent_id:     uint64(info.parent_id),
		dir_id:        uint64(info.dir_id),
		flags:         uint64(info.flags),
		uuid:          (*[16]uint8)(unsafe.Pointer(&info.uuid))[:16:16],
		parent_uuid:   (*[16]uint8)(unsafe.Pointer(&info.parent_uuid))[:16:16],
		received_uuid: (*[16]uint8)(unsafe.Pointer(&info.received_uuid))[:16:16],
		generation:    uint64(info.generation),
		ctransid:      uint64(info.ctransid),
		otransid:      uint64(info.otransid),
		stransid:      uint64(info.stransid),
		rtransid:      uint64(info.rtransid),
		ctime:         time.Unix(int64(info.ctime.tv_sec), int64(info.ctime.tv_nsec)),
		otime:         time.Unix(int64(info.otime.tv_sec), int64(info.otime.tv_nsec)),
		stime:         time.Unix(int64(info.stime.tv_sec), int64(info.stime.tv_nsec)),
		rtime:         time.Unix(int64(info.rtime.tv_sec), int64(info.rtime.tv_nsec)),
	}
	return subvol
}

func StrError(errInt uint32) error {
	if errInt != 0 {
		errStr := C.btrfs_util_strerror(errInt)
		return errors.New(C.GoString(errStr))
	}
	return nil
}

func Sync(path string) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := StrError(C.btrfs_util_sync(Cpath))
	return err
}

func SyncFd(fd int) error {
	err := StrError(C.btrfs_util_sync_fd(C.int(fd)))
	return err
}

func StartSync(path string) (uint64, error) {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	var transid C.uint64_t

	err := StrError(C.btrfs_util_start_sync(Cpath, &transid))
	return uint64(transid), err
}

func StratSyncFd(fd int) (uint64, error) {
	var transid C.uint64_t

	err := StrError(C.btrfs_util_start_sync_fd(C.int(fd), &transid))
	return uint64(transid), err
}

func WaitSync(path string, transid uint64) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	tid := C.uint64_t(transid)

	err := StrError(C.btrfs_util_wait_sync(Cpath, tid))
	return err
}

func WaitSyncFd(fd int, transid uint64) error {
	tid := C.uint64_t(transid)
	err := StrError(C.btrfs_util_wait_sync_fd(C.int(fd), tid))
	return err
}

func IsSubvolume(path string) (bool, error) {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))
	err := StrError(C.btrfs_util_is_subvolume(Cpath))
	if err == nil {
		return true, err
	}
	return false, err
}

func IsSubvolumeFd(fd int) (bool, error) {
	err := StrError(C.btrfs_util_is_subvolume_fd(C.int(fd)))
	if err == nil {
		return true, err
	}
	return false, err
}

func SubvolumeId(path string) (uint64, error) {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	var id_ret C.uint64_t
	err := StrError(C.btrfs_util_subvolume_id(Cpath, &id_ret))
	return uint64(id_ret), err
}

func SubvolumeIdFd(fd int) (uint64, error) {
	var id_ret C.uint64_t
	err := StrError(C.btrfs_util_subvolume_id_fd(C.int(fd), &id_ret))
	return uint64(id_ret), err
}

func SubvolumePath(path string, id uint64) (string, error) {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	var path_ret *C.char
	defer C.free(unsafe.Pointer(path_ret))

	err := StrError(C.btrfs_util_subvolume_path(Cpath, C.uint64_t(id), &path_ret))
	return C.GoString(path_ret), err
}

func SubvolumePathFd(fd int, id uint64) (string, error) {
	var path_ret *C.char
	defer C.free(unsafe.Pointer(path_ret))
	err := StrError(C.btrfs_util_subvolume_path_fd(C.int(fd), C.uint64_t(id), &path_ret))
	return C.GoString(path_ret), err
}

func SubvolumeInfo(path string, id uint64) (subvolumeInfo, error) {
	var info C.struct_btrfs_util_subvolume_info

	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := StrError(C.btrfs_util_subvolume_info(Cpath, C.uint64_t(id), &info))
	return newSubvolumeInfo(&info), err
}

func SubvolumeInfoFd(fd int, id uint64) (subvolumeInfo, error) {
	var info C.struct_btrfs_util_subvolume_info

	err := StrError(C.btrfs_util_subvolume_info_fd(C.int(fd), C.uint64_t(id), &info))
	return newSubvolumeInfo(&info), err
}

func GetSubvolumeReadOnly(path string) (bool, error) {
	var ret C.bool

	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := StrError(C.btrfs_util_get_subvolume_read_only(Cpath, &ret))
	return bool(ret), err
}

func GetSubvolumeReadOnlyFd(fd int) (bool, error) {
	var ret C.bool

	err := StrError(C.btrfs_util_get_subvolume_read_only_fd(C.int(fd), &ret))
	return bool(ret), err
}

func SetSubvolumeReadOnly(path string, read_only bool) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := StrError(C.btrfs_util_set_subvolume_read_only(Cpath, C.bool(read_only)))
	return err
}

func SetSubvolumeReadOnlyFd(fd int, read_only bool) error {
	err := StrError(C.btrfs_util_set_subvolume_read_only_fd(C.int(fd), C.bool(read_only)))
	return err
}
func GetDefaultSubvolume(path string) (uint64, error) {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	var id_ret C.uint64_t

	err := StrError(C.btrfs_util_get_default_subvolume(Cpath, &id_ret))
	return uint64(id_ret), err
}

func GetDefaultSubvolumeFd(fd int) (uint64, error) {
	var id_ret C.uint64_t
	err := StrError(C.btrfs_util_get_default_subvolume_fd(C.int(fd), &id_ret))
	return uint64(id_ret), err
}

func SetDefaultSubvolume(path string, id uint64) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := StrError(C.btrfs_util_set_default_subvolume(Cpath, C.uint64_t(id)))
	return err
}

func SetDefaultSubvolumeFd(fd int, id uint64) error {
	err := StrError(C.btrfs_util_set_default_subvolume_fd(C.int(fd), C.uint64_t(id)))
	return err
}

func CreateSubvolume(path string, qgroup_inherit *QgroupInherit) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := StrError(C.btrfs_util_create_subvolume(Cpath, 0, nil, qgroup_inherit.inherit))
	return err
}

func CreateSubvolumeFd(parent_fd int, name string, qgroup_inherit *QgroupInherit) error {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))

	err := StrError(C.btrfs_util_create_subvolume_fd(C.int(parent_fd), Cname, 0, nil, qgroup_inherit.inherit))
	return err
}

func CreateSnapshot(source string, path string, recursive bool, read_only bool, qgroup_inherit *QgroupInherit) error {
	Csource := C.CString(source)
	defer C.free(unsafe.Pointer(Csource))

	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	flags := 0

	if recursive {
		flags |= C.BTRFS_UTIL_CREATE_SNAPSHOT_RECURSIVE
	}

	if read_only {
		flags |= C.BTRFS_UTIL_CREATE_SNAPSHOT_READ_ONLY
	}

	err := StrError(C.btrfs_util_create_snapshot(Csource, Cpath, C.int(flags), nil, qgroup_inherit.inherit))
	return err
}

func CreateSnapshotFd(fd int, path string, recursive bool, read_only bool, qgroup_inherit *QgroupInherit) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	flags := 0

	if recursive {
		flags |= C.BTRFS_UTIL_CREATE_SNAPSHOT_RECURSIVE
	}

	if read_only {
		flags |= C.BTRFS_UTIL_CREATE_SNAPSHOT_READ_ONLY
	}

	err := StrError(C.btrfs_util_create_snapshot_fd(C.int(fd), Cpath, C.int(flags), nil, qgroup_inherit.inherit))
	return err
}

func CreateSnapshotFd2(fd int, parent_fd int, name string, recursive bool, read_only bool, qgroup_inherit *QgroupInherit) error {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))

	flags := 0

	if recursive {
		flags |= C.BTRFS_UTIL_CREATE_SNAPSHOT_RECURSIVE
	}

	if read_only {
		flags |= C.BTRFS_UTIL_CREATE_SNAPSHOT_READ_ONLY
	}

	err := StrError(C.btrfs_util_create_snapshot_fd2(C.int(fd), C.int(parent_fd), Cname, C.int(flags), nil, qgroup_inherit.inherit))
	return err
}

func DeleteSubvolume(path string, recursive bool) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	flags := 0

	if recursive {
		flags |= C.BTRFS_UTIL_DELETE_SUBVOLUME_RECURSIVE
	}

	err := StrError(C.btrfs_util_delete_subvolume(Cpath, C.int(flags)))
	return err
}

func DeleteSubvolumeFd(parent_fd int, name string, recursive bool) error {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))

	flags := 0

	if recursive {
		flags |= C.BTRFS_UTIL_DELETE_SUBVOLUME_RECURSIVE
	}

	err := StrError(C.btrfs_util_delete_subvolume_fd(C.int(parent_fd), Cname, C.int(flags)))
	return err
}

func DeleteSubvolumeByIdFd(parent_fd int, subvolid uint64) error {
	err := StrError(C.btrfs_util_delete_subvolume_by_id_fd(C.int(parent_fd), C.uint64_t(subvolid)))
	return err
}

func DeletedSubvolumes(path string) ([]uint64, error) {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	var n C.size_t
	var Cids *C.uint64_t
	defer C.free(unsafe.Pointer(Cids))

	err := StrError(C.btrfs_util_deleted_subvolumes(Cpath, &Cids, &n))

	var ids []uint64

	if n != 0 {
		ids = (*[1 << 31]uint64)(unsafe.Pointer(Cids))[:n:n]
	}
	return ids, err
}

func DeletedSubvolumesFd(fd int) ([]uint64, error) {
	var n C.size_t
	var Cids *C.uint64_t
	defer C.free(unsafe.Pointer(Cids))

	err := StrError(C.btrfs_util_deleted_subvolumes_fd(C.int(fd), &Cids, &n))

	var ids []uint64

	if n != 0 {
		ids = (*[1 << 31]uint64)(unsafe.Pointer(Cids))[:n:n]
	}
	return ids, err
}
