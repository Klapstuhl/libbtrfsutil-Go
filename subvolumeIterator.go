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

type SubvolumeIterator struct {
	iterator *C.struct_btrfs_util_subvolume_iterator
	fd       int
}

type SubvolumeIteratorData struct {
	err  error
	path string
	id   uint64
}

type SubvolumeIteratorInfo struct {
	err  error
	path string
	info subvolumeInfo
}

func CreateSubvolumeIterator(path string, top uint64, post_order bool) (*SubvolumeIterator, error) {
	it := new(SubvolumeIterator)

	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	flags := 0
	if post_order {
		flags |= C.BTRFS_UTIL_SUBVOLUME_ITERATOR_POST_ORDER
	}

	err := StrError(C.btrfs_util_create_subvolume_iterator(Cpath, C.uint64_t(top), C.int(flags), &it.iterator))
	it.fd = int(C.btrfs_util_subvolume_iterator_fd(it.iterator))
	return it, err
}

func CreateSubvolumeIteratorFd(fd int, top uint64, post_order bool) (*SubvolumeIterator, error) {
	it := new(SubvolumeIterator)

	flags := 0
	if post_order {
		flags |= C.BTRFS_UTIL_SUBVOLUME_ITERATOR_POST_ORDER
	}

	err := StrError(C.btrfs_util_create_subvolume_iterator_fd(C.int(fd), C.uint64_t(top), C.int(flags), &it.iterator))
	return it, err
}

func (it SubvolumeIterator) Destroy() {
	C.btrfs_util_destroy_subvolume_iterator(it.iterator)
}

func (it SubvolumeIterator) Next() <-chan SubvolumeIteratorData {
	ch := make(chan SubvolumeIteratorData)
	var Cpath *C.char
	defer C.free(unsafe.Pointer(Cpath))

	var id C.uint64_t

	go func() {
		for {
			err_id := C.btrfs_util_subvolume_iterator_next(it.iterator, &Cpath, &id)
			if err_id == 1 {
				break
			}
			ch <- SubvolumeIteratorData{
				err:  StrError(err_id),
				path: C.GoString(Cpath),
				id:   uint64(id),
			}
			C.free(unsafe.Pointer(Cpath))
		}
		close(ch)
	}()
	return ch
}

func (it SubvolumeIterator) NextInfo() <-chan SubvolumeIteratorInfo {
	ch := make(chan SubvolumeIteratorInfo)
	var Cpath *C.char
	defer C.free(unsafe.Pointer(Cpath))

	var info C.struct_btrfs_util_subvolume_info

	go func() {
		for {
			err_id := C.btrfs_util_subvolume_iterator_next_info(it.iterator, &Cpath, &info)
			if err_id == 1 {
				break
			}
			ch <- SubvolumeIteratorInfo{
				err:  StrError(err_id),
				path: C.GoString(Cpath),
				info: newSubvolumeInfo(&info),
			}
			C.free(unsafe.Pointer(Cpath))
		}
		close(ch)
	}()
	return ch
}
