// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BaseSysRoleMenuDao is the data access object for the table base_sys_role_menu.
type BaseSysRoleMenuDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  BaseSysRoleMenuColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// BaseSysRoleMenuColumns defines and stores column names for the table base_sys_role_menu.
type BaseSysRoleMenuColumns struct {
	Id         string //
	CreateTime string //
	UpdateTime string //
	DeletedAt  string //
	RoleId     string //
	MenuId     string //
}

// baseSysRoleMenuColumns holds the columns for the table base_sys_role_menu.
var baseSysRoleMenuColumns = BaseSysRoleMenuColumns{
	Id:         "id",
	CreateTime: "createTime",
	UpdateTime: "updateTime",
	DeletedAt:  "deleted_at",
	RoleId:     "roleId",
	MenuId:     "menuId",
}

// NewBaseSysRoleMenuDao creates and returns a new DAO object for table data access.
func NewBaseSysRoleMenuDao(handlers ...gdb.ModelHandler) *BaseSysRoleMenuDao {
	return &BaseSysRoleMenuDao{
		group:    "default",
		table:    "base_sys_role_menu",
		columns:  baseSysRoleMenuColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BaseSysRoleMenuDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BaseSysRoleMenuDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BaseSysRoleMenuDao) Columns() BaseSysRoleMenuColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BaseSysRoleMenuDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BaseSysRoleMenuDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *BaseSysRoleMenuDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
