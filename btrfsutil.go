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

// #cgo LDFLAGS: -lbtrfsutil
// #include <stdlib.h>
// #include <btrfsutil.h>
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

// SubvolumeInfo is a representation of a Btrfs subvolume or snapshot.
type SubvolumeInfo struct {
	Id           uint64
	ParentId     uint64
	DirId        uint64
	Flags        uint64
	UUID         string
	ParentUUID   string
	ReceivedUUID string
	Generation   uint64
	Ctransid     uint64
	Otransid     uint64
	Stransid     uint64
	Rtransid     uint64
	Ctime        time.Time
	Otime        time.Time
	Stime        time.Time
	Rtime        time.Time
}

func newSubvolumeInfo(info *C.struct_btrfs_util_subvolume_info) *SubvolumeInfo {
	subvol := SubvolumeInfo{
		Id:           uint64(info.id),
		ParentId:     uint64(info.parent_id),
		DirId:        uint64(info.dir_id),
		Flags:        uint64(info.flags),
		UUID:         uuidString(info.uuid),
		ParentUUID:   uuidString(info.parent_uuid),
		ReceivedUUID: uuidString(info.received_uuid),
		Generation:   uint64(info.generation),
		Ctransid:     uint64(info.ctransid),
		Otransid:     uint64(info.otransid),
		Stransid:     uint64(info.stransid),
		Rtransid:     uint64(info.rtransid),
		Ctime:        time.Unix(int64(info.ctime.tv_sec), int64(info.ctime.tv_nsec)),
		Otime:        time.Unix(int64(info.otime.tv_sec), int64(info.otime.tv_nsec)),
		Stime:        time.Unix(int64(info.stime.tv_sec), int64(info.stime.tv_nsec)),
		Rtime:        time.Unix(int64(info.rtime.tv_sec), int64(info.rtime.tv_nsec)),
	}
	return &subvol
}

func uuidString(uuid [16]C.uchar) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

// Sync forces a sync on a specific Btrfs filesystem.
func Sync(path string) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := getError(C.btrfs_util_sync(Cpath))
	return err
}

// See Sync.
func SyncFd(fd uintptr) error {
	err := getError(C.btrfs_util_sync_fd(C.int(fd)))
	return err
}

// StartsSync starts a sync on a specific Btrfs filesystem but dose not wait for it.
func StartSync(path string) (uint64, error) {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	var transid C.uint64_t

	err := getError(C.btrfs_util_start_sync(Cpath, &transid))
	return uint64(transid), err
}

// See StartSync.
func StratSyncFd(fd uintptr) (uint64, error) {
	var transid C.uint64_t

	err := getError(C.btrfs_util_start_sync_fd(C.int(fd), &transid))
	return uint64(transid), err
}

// WaitSync waits for a transaction with a given ID to sync.
// If the given ID is zero, WaitSync waits for the current transaction.
func WaitSync(path string, transid uint64) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	tid := C.uint64_t(transid)

	err := getError(C.btrfs_util_wait_sync(Cpath, tid))
	return err
}

// See WaitSync.
func WaitSyncFd(fd uintptr, transid uint64) error {
	tid := C.uint64_t(transid)
	err := getError(C.btrfs_util_wait_sync_fd(C.int(fd), tid))
	return err
}

// IsSubvolume returns whether a given path is a Btrfs subvolume.
func IsSubvolume(path string) (bool, error) {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))
	err := getError(C.btrfs_util_is_subvolume(Cpath))
	if err == nil {
		return true, err
	}
	return false, err
}

// See IsSubvolume.
func IsSubvolumeFd(fd uintptr) (bool, error) {
	err := getError(C.btrfs_util_is_subvolume_fd(C.int(fd)))
	if err == nil {
		return true, err
	}
	return false, err
}

// SubvolumeId returns the ID of the subvolume containing a given path.
func SubvolumeId(path string) (uint64, error) {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	var id_ret C.uint64_t
	err := getError(C.btrfs_util_subvolume_id(Cpath, &id_ret))
	return uint64(id_ret), err
}

// See SubvolumeId.
func SubvolumeIdFd(fd uintptr) (uint64, error) {
	var id_ret C.uint64_t
	err := getError(C.btrfs_util_subvolume_id_fd(C.int(fd), &id_ret))
	return uint64(id_ret), err
}

// SubvolumePath returns the path of the subvolume with a given ID.
func SubvolumePath(path string, id uint64) (string, error) {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	var path_ret *C.char
	defer C.free(unsafe.Pointer(path_ret))

	err := getError(C.btrfs_util_subvolume_path(Cpath, C.uint64_t(id), &path_ret))
	return C.GoString(path_ret), err
}

// See SubvolumePath.
func SubvolumePathFd(fd uintptr, id uint64) (string, error) {
	var path_ret *C.char
	defer C.free(unsafe.Pointer(path_ret))
	err := getError(C.btrfs_util_subvolume_path_fd(C.int(fd), C.uint64_t(id), &path_ret))
	return C.GoString(path_ret), err
}

// GetSubvolumeInfo returns information about a subvolume with a given ID or path.
// The given path may be any path in the Btrfs filesystem; it dose not have to
// refer to a subvolume unless id is zero. If the given ID is zero,
// the subvolume ID of the subvolume containing path is used.
func GetSubvolumeInfo(path string, id uint64) (*SubvolumeInfo, error) {
	var info C.struct_btrfs_util_subvolume_info

	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := getError(C.btrfs_util_subvolume_info(Cpath, C.uint64_t(id), &info))
	if err != nil {
		return nil, err
	}
	return newSubvolumeInfo(&info), nil
}

// See GetSubvolumeInfo.
func GetSubvolumeInfoFd(fd uintptr, id uint64) (*SubvolumeInfo, error) {
	var info C.struct_btrfs_util_subvolume_info

	err := getError(C.btrfs_util_subvolume_info_fd(C.int(fd), C.uint64_t(id), &info))
	if err != nil {
		return nil, err
	}
	return newSubvolumeInfo(&info), nil
}

// GetSubvolumeReadOnly returns whether a subvolume is read-only.
func GetSubvolumeReadOnly(path string) (bool, error) {
	var ret C.bool

	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := getError(C.btrfs_util_get_subvolume_read_only(Cpath, &ret))
	return bool(ret), err
}

// See GetSubvolumeReadOnly.
func GetSubvolumeReadOnlyFd(fd uintptr) (bool, error) {
	var ret C.bool

	err := getError(C.btrfs_util_get_subvolume_read_only_fd(C.int(fd), &ret))
	return bool(ret), err
}

// SetSubvolumeReadOnly sets whether a subvolume is read-only.
func SetSubvolumeReadOnly(path string, read_only bool) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := getError(C.btrfs_util_set_subvolume_read_only(Cpath, C.bool(read_only)))
	return err
}

// See SetSubvolumeReadOnly.
func SetSubvolumeReadOnlyFd(fd uintptr, read_only bool) error {
	err := getError(C.btrfs_util_set_subvolume_read_only_fd(C.int(fd), C.bool(read_only)))
	return err
}

// GetDefaultSubvolume returns the default subvolume ID for a filesystem.
func GetDefaultSubvolume(path string) (uint64, error) {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	var id_ret C.uint64_t

	err := getError(C.btrfs_util_get_default_subvolume(Cpath, &id_ret))
	return uint64(id_ret), err
}

// See GetDefaultSubvolume.
func GetDefaultSubvolumeFd(fd uintptr) (uint64, error) {
	var id_ret C.uint64_t
	err := getError(C.btrfs_util_get_default_subvolume_fd(C.int(fd), &id_ret))
	return uint64(id_ret), err
}

// SetDefaultSubvolume sets the default subvolume for a filesystem.
// The given path may be any path in the Btrfs filesystem; it dose not have to
// refer to a subvolume unless id is zero.
// If the given ID is zero, the subvolume ID of the subvolume containing path is used.
func SetDefaultSubvolume(path string, id uint64) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := getError(C.btrfs_util_set_default_subvolume(Cpath, C.uint64_t(id)))
	return err
}

// See SetDefaultSubvolume.
func SetDefaultSubvolumeFd(fd uintptr, id uint64) error {
	err := getError(C.btrfs_util_set_default_subvolume_fd(C.int(fd), C.uint64_t(id)))
	return err
}

// CreateSubvolume creates a new subvolume under a given path.
func CreateSubvolume(path string) error {
	return CreateSubvolumeWithQgroup(path, &QgroupInherit{})
}

// CreateSubvolumeWithQgroup creates a new subvolume under a given path, with Qgroups to inherit from.
func CreateSubvolumeWithQgroup(path string, qgroup_inherit *QgroupInherit) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := getError(C.btrfs_util_create_subvolume(Cpath, 0, nil, qgroup_inherit.inherit))
	return err
}

func CreateSubvolumeFd(parent_fd uintptr, name string) error {
	return CreateSubvolumeWithQgroupFd(parent_fd, name, &QgroupInherit{})
}

// CreateSubvolumeWithQgroupFd creates a new subvolume given its parent file descriptor, a name and Qgroups to inherit from.
func CreateSubvolumeWithQgroupFd(parent_fd uintptr, name string, qgroup_inherit *QgroupInherit) error {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))

	err := getError(C.btrfs_util_create_subvolume_fd(C.int(parent_fd), Cname, 0, nil, qgroup_inherit.inherit))
	return err
}

// CreateSnapshot creates a new snapshot from a source subvolume path.
// If source is not a subvolume the subvolume containing source will be snapshotted
func CreateSnapshot(source string, path string, recursive bool, read_only bool) error {
	return CreateSnapshotWithQgroup(source, path, recursive, read_only, &QgroupInherit{})
}

// CreateSnapshotWithQgroup creates a new snapshot from a source subvolume path with Qgroups to inherit from.
func CreateSnapshotWithQgroup(source string, path string, recursive bool, read_only bool, qgroup_inherit *QgroupInherit) error {
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

	err := getError(C.btrfs_util_create_snapshot(Csource, Cpath, C.int(flags), nil, qgroup_inherit.inherit))
	return err
}

// See CreateSnapshot
func CreateSnapshotFd(fd uintptr, path string, recursive bool, read_only bool) error {
	return CreateSnapshotWithQgroupFd(fd, path, recursive, read_only, &QgroupInherit{})
}

// See CreateSnapshotWithQgroup.
func CreateSnapshotWithQgroupFd(fd uintptr, path string, recursive bool, read_only bool, qgroup_inherit *QgroupInherit) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	flags := 0

	if recursive {
		flags |= C.BTRFS_UTIL_CREATE_SNAPSHOT_RECURSIVE
	}

	if read_only {
		flags |= C.BTRFS_UTIL_CREATE_SNAPSHOT_READ_ONLY
	}

	err := getError(C.btrfs_util_create_snapshot_fd(C.int(fd), Cpath, C.int(flags), nil, qgroup_inherit.inherit))
	return err
}

// CreateSnapshotFd2 creates a new snapshot form a source subvolume file descriptor, a target parent file descriptor and name.
func CreateSnapshotFd2(fd uintptr, parent_fd uintptr, name string, recursive bool, read_only bool) error {
	return CreateSnapshotWithQgroupFd2(fd, parent_fd, name, recursive, read_only, &QgroupInherit{})
}

// CreateSnapshotWithQgroupFd2 creates a new snapshot form a source subvolume file descriptor, a target parent file descriptor and name,
// with Qgroups to inherit from.
func CreateSnapshotWithQgroupFd2(fd uintptr, parent_fd uintptr, name string, recursive bool, read_only bool, qgroup_inherit *QgroupInherit) error {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))

	flags := 0

	if recursive {
		flags |= C.BTRFS_UTIL_CREATE_SNAPSHOT_RECURSIVE
	}

	if read_only {
		flags |= C.BTRFS_UTIL_CREATE_SNAPSHOT_READ_ONLY
	}

	err := getError(C.btrfs_util_create_snapshot_fd2(C.int(fd), C.int(parent_fd), Cname, C.int(flags), nil, qgroup_inherit.inherit))
	return err
}

// DeleteSubvolume deletes a subvolume or snapshot.
// If recursive is set subvolumes beneath the given subvolume will be deleted befor
// attempting to delete the given subvolume.
// Unless the filesystem is mounted with 'user_subvol_rm_allow', appropriate privileges are required (CAP_SYS_ADMIN).
func DeleteSubvolume(path string, recursive bool) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	flags := 0

	if recursive {
		flags |= C.BTRFS_UTIL_DELETE_SUBVOLUME_RECURSIVE
	}

	err := getError(C.btrfs_util_delete_subvolume(Cpath, C.int(flags)))
	return err
}

// DeleteSubvolumeFd deletes a subvolume or snapshot by its parent file descriptor and name.
// See DeleteSubvolume.
func DeleteSubvolumeFd(parent_fd uintptr, name string, recursive bool) error {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))

	flags := 0

	if recursive {
		flags |= C.BTRFS_UTIL_DELETE_SUBVOLUME_RECURSIVE
	}

	err := getError(C.btrfs_util_delete_subvolume_fd(C.int(parent_fd), Cname, C.int(flags)))
	return err
}

// DeleteSubvolumeByIdFd deletes a subvolume or snapshot by its parent file descriptor and id.
// See DeleteSubvolume
func DeleteSubvolumeByIdFd(parent_fd uintptr, subvolid uint64) error {
	err := getError(C.btrfs_util_delete_subvolume_by_id_fd(C.int(parent_fd), C.uint64_t(subvolid)))
	return err
}

// DeletedSubvolumes returns a list of subvolume IDs which have been deleted but not yet cleaned up.
func DeletedSubvolumes(path string) ([]uint64, error) {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	var n C.size_t
	var Cids *C.uint64_t
	defer C.free(unsafe.Pointer(Cids))

	err := getError(C.btrfs_util_deleted_subvolumes(Cpath, &Cids, &n))

	var ids []uint64

	if n != 0 {
		ids = (*[1 << 31]uint64)(unsafe.Pointer(Cids))[:n:n]
	}
	return ids, err
}

// See DeletedSubvolumesFd.
func DeletedSubvolumesFd(fd uintptr) ([]uint64, error) {
	var n C.size_t
	var Cids *C.uint64_t
	defer C.free(unsafe.Pointer(Cids))

	err := getError(C.btrfs_util_deleted_subvolumes_fd(C.int(fd), &Cids, &n))

	var ids []uint64

	if n != 0 {
		ids = (*[1 << 31]uint64)(unsafe.Pointer(Cids))[:n:n]
	}
	return ids, err
}
