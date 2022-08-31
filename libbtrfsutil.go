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

package libbtrfsutil

// #cgo LDFLAGS: -lbtrfsutil
// #include <stdlib.h>
// #include <btrfsutil.h>
import "C"
import (
	"time"
	"unsafe"
)

// SubvolumeInfo is a representation of a Btrfs subvolume or snapshot.
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

// Sync forces a sync on a specific Btrfs filesystem.
func Sync(path string) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := getError(C.btrfs_util_sync(Cpath))
	return err
}

// See Sync.
func SyncFd(fd int) error {
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
func StratSyncFd(fd int) (uint64, error) {
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
func WaitSyncFd(fd int, transid uint64) error {
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
func IsSubvolumeFd(fd int) (bool, error) {
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
func SubvolumeIdFd(fd int) (uint64, error) {
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
func SubvolumePathFd(fd int, id uint64) (string, error) {
	var path_ret *C.char
	defer C.free(unsafe.Pointer(path_ret))
	err := getError(C.btrfs_util_subvolume_path_fd(C.int(fd), C.uint64_t(id), &path_ret))
	return C.GoString(path_ret), err
}

// SubvolumeInfo returns information about a subvolume with a given ID or path.
// The given path may be any path in the Btrfs filesystem; it dose not have to
// refer to a subvolume unless id is zero. If the given ID is zero,
// the subvolume ID of the subvolume containing path is used.
func SubvolumeInfo(path string, id uint64) (subvolumeInfo, error) {
	var info C.struct_btrfs_util_subvolume_info

	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := getError(C.btrfs_util_subvolume_info(Cpath, C.uint64_t(id), &info))
	return newSubvolumeInfo(&info), err
}

// See SubvolumeInfo.
func SubvolumeInfoFd(fd int, id uint64) (subvolumeInfo, error) {
	var info C.struct_btrfs_util_subvolume_info

	err := getError(C.btrfs_util_subvolume_info_fd(C.int(fd), C.uint64_t(id), &info))
	return newSubvolumeInfo(&info), err
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
func GetSubvolumeReadOnlyFd(fd int) (bool, error) {
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
func SetSubvolumeReadOnlyFd(fd int, read_only bool) error {
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
func GetDefaultSubvolumeFd(fd int) (uint64, error) {
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
func SetDefaultSubvolumeFd(fd int, id uint64) error {
	err := getError(C.btrfs_util_set_default_subvolume_fd(C.int(fd), C.uint64_t(id)))
	return err
}

// CreateSubvolume creates a new subvolume under a given path, with Qgroups to inherit from.
// qgroup_inherit can be nil if the new subvolume should not inherit any Qgroups.
func CreateSubvolume(path string, qgroup_inherit *QgroupInherit) error {
	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	err := getError(C.btrfs_util_create_subvolume(Cpath, 0, nil, qgroup_inherit.inherit))
	return err
}

// CreateSubvolumeFd creates a new subvolume given its parent, a name and Qgroups to inherit from.
// qgroup_inherit can be nil if the new subvolume should not inherit any Qgroups.
func CreateSubvolumeFd(parent_fd int, name string, qgroup_inherit *QgroupInherit) error {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))

	err := getError(C.btrfs_util_create_subvolume_fd(C.int(parent_fd), Cname, 0, nil, qgroup_inherit.inherit))
	return err
}

// CreateSnapshot creates a new snapshot from a source subvolume path.
// qgroup_inherit can be nil if the new subvolume should not inherit any Qgroups.
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

	err := getError(C.btrfs_util_create_snapshot(Csource, Cpath, C.int(flags), nil, qgroup_inherit.inherit))
	return err
}

// See CreateSnapshot.
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

	err := getError(C.btrfs_util_create_snapshot_fd(C.int(fd), Cpath, C.int(flags), nil, qgroup_inherit.inherit))
	return err
}

// CreateSnapshotFd2 creates a new snapshot form a source subvolume file descriptor and a target parent file descriptor and name.
// qgroup_inherit can be nil if the new subvolume should not inherit any Qgroups.
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
func DeleteSubvolumeFd(parent_fd int, name string, recursive bool) error {
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
func DeleteSubvolumeByIdFd(parent_fd int, subvolid uint64) error {
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
func DeletedSubvolumesFd(fd int) ([]uint64, error) {
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
