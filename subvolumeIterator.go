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

// #cgo LDFLAGS: -lbtrfsutil
// #include <stdlib.h>
// #include <btrfsutil.h>
import "C"
import (
	"unsafe"
)

type SubvolumeIteratorResult struct {
	Path string
	Id   uint64
}

type SubvolumeInfoIteratorResult struct {
	Path string
	Info *SubvolumeInfo
}

type SubvolumeIterator struct {
	lastResult *SubvolumeIteratorResult
	lastErr    error

	iterator *C.struct_btrfs_util_subvolume_iterator
}

// CreateSubvolumeIterator creates an iterator over subvolumes in a Btrfs filesystem.
// Lists all subvolumes beneath (but not including) the subvolume with the ID top.
// The given path may be any path in the Btrfs filesystem; it dose not have to
// refer to a subvolume unless top is zero. If the as top given ID is zero,
// the subvolume ID of the subvolume containing path is used.
// By default subvolumes are listed pre-order e.g., foo will be yielded before foo/bar.
// This behavior can be reversed by setting post_order.
// The returned SubvolumeIterator struct must be freed with Destroy().
func CreateSubvolumeIterator(path string, top uint64, post_order bool) (*SubvolumeIterator, error) {
	it := new(SubvolumeIterator)

	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	flags := 0
	if post_order {
		flags |= C.BTRFS_UTIL_SUBVOLUME_ITERATOR_POST_ORDER
	}

	err := getError(C.btrfs_util_create_subvolume_iterator(Cpath, C.uint64_t(top), C.int(flags), &it.iterator))
	return it, err
}

// See CreateSubvolumeIterator.
func CreateSubvolumeIteratorFd(fd uintptr, top uint64, post_order bool) (*SubvolumeIterator, error) {
	it := new(SubvolumeIterator)

	flags := 0
	if post_order {
		flags |= C.BTRFS_UTIL_SUBVOLUME_ITERATOR_POST_ORDER
	}

	err := getError(C.btrfs_util_create_subvolume_iterator_fd(C.int(fd), C.uint64_t(top), C.int(flags), &it.iterator))
	return it, err
}

// Fd returns the file descriptor referencing the SubvolumeIterator
func (it *SubvolumeIterator) Fd() uintptr {
	return uintptr(C.btrfs_util_subvolume_iterator_fd(it.iterator))
}

// Destroy destroys the SubvolumeIterator.
func (it *SubvolumeIterator) Destroy() {
	C.btrfs_util_destroy_subvolume_iterator(it.iterator)
	it.iterator = nil
}

// HasNext returns true if the SubvolumeIterator has a next value.
func (it *SubvolumeIterator) HasNext() bool {
	var Cpath *C.char
	defer C.free(unsafe.Pointer(Cpath))

	var id C.uint64_t
	it.lastErr = getError(C.btrfs_util_subvolume_iterator_next(it.iterator, &Cpath, &id))
	if it.lastErr == ErrStopIteration {
		it.lastResult = nil
		return false
	}

	it.lastResult = &SubvolumeIteratorResult{C.GoString(Cpath), uint64(id)}
	return true
}

// GetNext gets the Path and Id of the next subvolume from a SubvolumeIterator.
func (it *SubvolumeIterator) GetNext() (*SubvolumeIteratorResult, error) {
	if it.lastErr != nil {
		return nil, it.lastErr
	}
	return it.lastResult, it.lastErr
}

type SubvolumeInfoIterator struct {
	lastResult *SubvolumeInfoIteratorResult
	lastErr    error

	iterator *C.struct_btrfs_util_subvolume_iterator
}

// Identical to CreateSubvolumeIterator but GetNext() returns a SubvolumeInfo instead of a subvolume Id.
// The returned SubvolumeInfoIterator struct must be freed with Destroy().
func CreateSubvolumeInfoIterator(path string, top uint64, post_order bool) (*SubvolumeInfoIterator, error) {
	it := new(SubvolumeInfoIterator)

	Cpath := C.CString(path)
	defer C.free(unsafe.Pointer(Cpath))

	flags := 0
	if post_order {
		flags |= C.BTRFS_UTIL_SUBVOLUME_ITERATOR_POST_ORDER
	}

	err := getError(C.btrfs_util_create_subvolume_iterator(Cpath, C.uint64_t(top), C.int(flags), &it.iterator))

	return it, err
}

// See CreateSubvolumeInfoIterator.
func CreateSubvolumeInfoIteratorFd(fd uintptr, top uint64, post_order bool) (*SubvolumeInfoIterator, error) {
	it := new(SubvolumeInfoIterator)

	flags := 0
	if post_order {
		flags |= C.BTRFS_UTIL_SUBVOLUME_ITERATOR_POST_ORDER
	}

	err := getError(C.btrfs_util_create_subvolume_iterator_fd(C.int(fd), C.uint64_t(top), C.int(flags), &it.iterator))
	return it, err
}

// Fd returns the file descriptor referencing the SubvolumeInfoIterator
func (it *SubvolumeInfoIterator) Fd() uintptr {
	return uintptr(C.btrfs_util_subvolume_iterator_fd(it.iterator))
}

// Destroy destroys the SubvolumeInfoIterator.
func (it *SubvolumeInfoIterator) Destroy() {
	C.btrfs_util_destroy_subvolume_iterator(it.iterator)
	it.iterator = nil
}

// HasNext returns true if the SubvolumeInfoIterator has a next value.
func (it *SubvolumeInfoIterator) HasNext() bool {
	var Cpath *C.char
	defer C.free(unsafe.Pointer(Cpath))

	var info C.struct_btrfs_util_subvolume_info
	it.lastErr = getError(C.btrfs_util_subvolume_iterator_next_info(it.iterator, &Cpath, &info))
	if it.lastErr == ErrStopIteration {
		it.lastResult = nil
		return false
	}

	it.lastResult = &SubvolumeInfoIteratorResult{C.GoString(Cpath), newSubvolumeInfo(&info)}
	return true
}

// GetNext gets the Path and SubvolumeInfo of the next subvolume from a SubvolumeInfoIterator.
func (it *SubvolumeInfoIterator) GetNext() (*SubvolumeInfoIteratorResult, error) {
	if it.lastErr != nil {
		return nil, it.lastErr
	}
	return it.lastResult, it.lastErr
}
