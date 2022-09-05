package btrfsutil

import "errors"

var (
	ErrStopIteration          = errors.New("stop iteration")
	ErrNoMemory               = errors.New("cannot allocate memory")
	ErrInvalidArgument        = errors.New("invalid argument")
	ErrNotBtrfs               = errors.New("not a Btrfs filesystem")
	ErrNotSubvolume           = errors.New("not a Btrfs subvolume")
	ErrSubvolumeNotFound      = errors.New("subvolume not found")
	ErrOpenFailed             = errors.New("could not open")
	ErrRmdirFailed            = errors.New("could not rmdir")
	ErrUnlinkFailed           = errors.New("could not unlink")
	ErrStatFailed             = errors.New("could not stat")
	ErrStatfsFailed           = errors.New("could not statfs")
	ErrSearchFailed           = errors.New("could not search B-tree")
	ErrInoLookupFailed        = errors.New("could not lookup inode")
	ErrSubvolGetflagsFailed   = errors.New("could not get subvolume flags")
	ErrSubvolSetflagsFailed   = errors.New("could not set subvolume flags")
	ErrSubvolCreateFailed     = errors.New("could not create subvolume")
	ErrSnapCreateFailed       = errors.New("could not create snapshot")
	ErrSnapDestroyFailed      = errors.New("could not destroy subvolume/snapshot")
	ErrDefaultSubvolFailed    = errors.New("could not set default subvolume")
	ErrSyncFailed             = errors.New("could not sync filesystem")
	ErrStartSyncFailed        = errors.New("could not start filesystem sync")
	ErrWaitSyncFailed         = errors.New("could not wait for filesystem sync")
	ErrGetSubvolInfoFailed    = errors.New("could not get subvolume information with BTRFS_IOC_GET_SUBVOL_INFO")
	ErrGetSubvolRootrefFailed = errors.New("could not get rootref information with BTRFS_IOC_GET_SUBVOL_ROOTREF")
	ErrInoLookupUserFailed    = errors.New("could not resolve subvolume path with BTRFS_IOC_INO_LOOKUP_USER")
	ErrFsInfoFailed           = errors.New("could not get filesystem information")
)

var errorMap = map[uint32]error{
	1:  ErrStopIteration,
	2:  ErrNoMemory,
	3:  ErrInvalidArgument,
	4:  ErrNotBtrfs,
	5:  ErrNotSubvolume,
	6:  ErrSubvolumeNotFound,
	7:  ErrOpenFailed,
	8:  ErrRmdirFailed,
	9:  ErrUnlinkFailed,
	10: ErrStatFailed,
	11: ErrStatfsFailed,
	12: ErrSearchFailed,
	13: ErrInoLookupFailed,
	14: ErrSubvolGetflagsFailed,
	15: ErrSubvolSetflagsFailed,
	16: ErrSubvolCreateFailed,
	17: ErrSnapCreateFailed,
	18: ErrSnapDestroyFailed,
	19: ErrDefaultSubvolFailed,
	20: ErrSyncFailed,
	21: ErrStartSyncFailed,
	22: ErrWaitSyncFailed,
	23: ErrGetSubvolInfoFailed,
	24: ErrGetSubvolRootrefFailed,
	25: ErrInoLookupUserFailed,
	26: ErrFsInfoFailed,
}

func getError(errInt uint32) error {
	if errInt != 0 {
		return errorMap[errInt]
	}
	return nil
}
