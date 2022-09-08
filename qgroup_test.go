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
	"reflect"
	"testing"
)

func TestQgroupInherit(t *testing.T) {
	tests := []struct {
		name string
		want []uint64
	}{
		{"1,2", []uint64{1, 2}},
		{"77,934,3", []uint64{77, 934, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inherit, err := CreateQgroupInherit()
			if err != nil {
				t.Errorf("CreateQgroupInherit() error = %v", err)
			}
			defer inherit.Destroy()

			for _, i := range tt.want {
				if err := inherit.AddGroup(i); err != nil {
					t.Errorf("QgroupInherit.AddGroup(%v) error = %v", i, err)
				}
			}

			if got := inherit.GetGroups(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QgroupInherit.GetGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}
