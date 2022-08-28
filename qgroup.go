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
	"unsafe"
)

type QgroupInherit struct {
	inherit *C.struct_btrfs_util_qgroup_inherit
}

func NewQgroupInherit() (*QgroupInherit, error) {
	q := new(QgroupInherit)
	err := StrError(C.btrfs_util_create_qgroup_inherit(0, &q.inherit))
	return q, err
}

func (q QgroupInherit) Destroy() {
	C.btrfs_util_destroy_qgroup_inherit(q.inherit)
}

func (q QgroupInherit) AddGroup(groupid uint64) error {
	err := StrError(C.btrfs_util_qgroup_inherit_add_group(&q.inherit, C.uint64_t(groupid)))
	return err
}

func (q QgroupInherit) GetGroups() []uint64 {
	var n C.size_t
	var Cgroups *C.uint64_t
	defer C.free(unsafe.Pointer(Cgroups))

	C.btrfs_util_qgroup_inherit_get_groups(q.inherit, &Cgroups, &n)

	groups := (*[1 << 31]uint64)(unsafe.Pointer(Cgroups))[:n:n]
	return groups
}
