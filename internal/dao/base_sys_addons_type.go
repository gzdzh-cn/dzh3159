// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"dzhgo/internal/dao/internal"
)

// internalBaseSysAddonsTypesDao is internal type for wrapping internal DAO implements.
type internalBaseSysAddonsTypesDao = *internal.BaseSysAddonsTypesDao

// baseSysAddonsTypesDao is the data access object for table base_sys_addons_types.
// You can define custom methods on it to extend its functionality as you wish.
type baseSysAddonsTypesDao struct {
	internalBaseSysAddonsTypesDao
}

var (
	// BaseSysAddonsTypes is globally public accessible object for table base_sys_addons_types operations.
	BaseSysAddonsTypes = baseSysAddonsTypesDao{
		internal.NewBaseSysAddonsTypesDao(),
	}
)

// Fill with you ideas below.
