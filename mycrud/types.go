package mycrud

type CrudService struct {
	GetOne     GetOneImpl
	GetMany    GetManyImpl
	CreateOne  CreateOneImpl
	CreateMany CreateManyImpl
	UpdateOne  UpdateOneImpl
	UpdateMany UpdateManyImpl
	DeleteOne  DeleteOneImpl
	DeleteMany DeleteManyImpl

	GetRelationOne  GetRelationOneImpl
	GetRelationMany GetRelationManyImpl

	CreateRelationOne  CreateRelationOneImpl
	CreateRelationMany CreateRelationManyImpl
}

func (c *CrudService) Handler() {
	if c.GetOne != nil {
		c.GetOne.GetOneHandler()
	}
	if c.GetMany != nil {
		c.GetMany.GetManyHandler()
	}

	if c.CreateOne != nil {
		c.CreateOne.CreateOneHandler()
	}

	if c.CreateMany != nil {
		c.CreateMany.CreateManyHandler()
	}

	if c.UpdateOne != nil {
		c.UpdateOne.UpdateOneHandler()
	}

	if c.UpdateMany != nil {
		c.UpdateMany.UpdateManyHandler()
	}

	if c.DeleteOne != nil {
		c.DeleteOne.DeleteOneHandler()
	}

	if c.DeleteMany != nil {
		c.DeleteMany.DeleteManyHandler()
	}

	if c.GetRelationOne != nil {
		c.GetRelationOne.GetRelationOneHandler()
	}

	if c.GetRelationMany != nil {
		c.GetRelationMany.GetRelationManyHandler()
	}

	if c.CreateRelationOne != nil {
		c.CreateRelationOne.CreateRelationOneHandler()
	}

	if c.CreateRelationMany != nil {
		c.CreateRelationMany.CreateRelationManyHandler()
	}

}
