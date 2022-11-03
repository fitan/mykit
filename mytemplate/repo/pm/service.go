package pm

import (
	"context"
	"github.com/fitan/mykit/mygorm"
	"github.com/fitan/mykit/mytemplate/repo/types"
)

type Middleware func(Service) Service

type Service interface {
	Create(ctx context.Context, v types.PhysicalMachine) (err error)
	List(ctx context.Context, page int, pageSize int, order string, wheres ...interface{}) (list []types.PhysicalMachine, total int64, err error)
	Delete(ctx context.Context, uuid string) (err error)
	Update(ctx context.Context, uuid string, v types.PhysicalMachine) (err error)
	Get(ctx context.Context, uuid string) (v types.PhysicalMachine, err error)
	Filter(ctx context.Context, fieldName string, value string, filter *types.PhysicalMachine) (list []map[string]interface{}, err error)
	// 是否存在
	Exist(ctx context.Context, wheres ...interface{}) (exist bool, err error)
}

type service struct {
	myDB *mygorm.DB
}

func (s *service) Exist(ctx context.Context, wheres ...interface{}) (exist bool, err error) {
	db := s.myDB.GetDb(ctx).Model(&types.PhysicalMachine{})
	for _, where := range wheres {
		db = db.Where(where)
	}
	var count int64
	err = db.Count(&count).Error
	if err != nil {
		return
	}

	return count > 0, nil
}

func (s *service) Create(ctx context.Context, v types.PhysicalMachine) (err error) {
	return s.myDB.GetDb(ctx).Create(&v).Error
}

func (s *service) List(
	ctx context.Context, page int, pageSize int, order string, wheres ...interface{},
) (list []types.PhysicalMachine, total int64, err error) {
	tx := s.myDB.GetDb(ctx).Model(&types.PhysicalMachine{})

	if order != "" {
		tx.Order(order)
	} else {
		tx.Order("id desc")
	}

	for where := range wheres {
		tx.Where(where)
	}

	err = tx.Count(&total).Error
	if err != nil {
		return
	}

	if page != 0 && pageSize != 0 {
		tx.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	return list, total, tx.Find(&list).Error
}

func (s *service) Delete(ctx context.Context, uuid string) (err error) {
	return s.myDB.GetDb(ctx).Delete(&types.PhysicalMachine{}, "uuid = ?", uuid).Error
}

func (s *service) Update(ctx context.Context, uuid string, v types.PhysicalMachine) (err error) {
	return s.myDB.GetDb(ctx).Model(&types.PhysicalMachine{}).Where("uuid = ?", uuid).Save(&v).Error
}

func (s *service) Get(ctx context.Context, uuid string) (v types.PhysicalMachine, err error) {
	return v, s.myDB.GetDb(ctx).First(&v, "uuid = ?", uuid).Error
}

func (s *service) Filter(ctx context.Context, fieldName string, value string, filter *types.PhysicalMachine) (list []map[string]interface{}, err error) {
	tx := s.myDB.GetDb(ctx).Model(&types.PhysicalMachine{})
	if filter != nil {
		tx = tx.Where(filter)
	}
	return list, tx.Select("uuid", fieldName).Group(fieldName).Where(fieldName+" like ?", "%"+value+"%").Find(&list).Error
}

func New(db *mygorm.DB) Service {
	return &service{myDB: db}
}
