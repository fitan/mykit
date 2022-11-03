/**
 * @Time : 2022/9/29 9:55 AM
 * @Author : soupzhb@gmail.com
 * @File : template.go
 * @Software: GoLand
 */

package types

import (
	"gorm.io/gorm"
)

type Template struct {
	gorm.Model
	Name   string `gorm:"column:name;notnull;comment:模板标识" json:"name"`
	Alias  string `gorm:"column:alias;notnull;comment:模板名称" json:"alias"`
	Type   string `gorm:"column:type;notnull;comment:模板类型" json:"type"`
	Remark string `gorm:"column:remark;null;comment:备注" json:"remark"`
	Value  string `gorm:"column:value;null;type:text;comment:模板内容" json:"value"`
}

func (Template) TableName() string {
	return "template"
}
