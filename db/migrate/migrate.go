package migrate

import "gorm.io/gorm"

type Model struct {
	Id        int
	Migration string
	Batch     int
	*gorm.DB  `gorm:"-"`
}

func (*Model) TableName() string {
	return "migrations"
}

func (model *Model) GetRan() ([]*Model, error) {
	var migrations []*Model
	if err := model.Model(model).
		Order("batch asc").
		Order("migration asc").
		Find(&migrations).Error; err != nil {
		return nil, err
	}
	return migrations, nil
}

func (model *Model) GetMigrations(step int) ([]*Model, error) {
	var migrations []*Model
	if err := model.Model(model).
		Where("batch >= ?", 1).
		Where("batch desc").
		Where("migration desc").
		Limit(step).Find(&migrations).Error; err != nil {
		return migrations, err
	}
	return migrations, nil
}

func (model *Model) GetLast() ([]*Model, error) {
	var migrations []*Model
	maxBatch, err := model.GetLastBatchNumber()
	if err != nil {
		return migrations, err
	}
	if err := model.Model(model).
		Where("batch = ?", maxBatch).
		Order("migration desc").
		Find(&migrations).
		Error; err != nil {
		return migrations, err
	}
	return migrations, nil
}
func (model *Model) GetMigrationBatches() (map[string]int, error) {
	var migrations []*Model
	batchesMap := make(map[string]int)
	if err := model.Model(model).
		Order("batch asc").
		Order("migration asc").
		Find(&migrations).Error; err != nil {
		return nil, err
	}
	for _, migration := range migrations {
		batchesMap[migration.Migration] = migration.Batch
	}
	return batchesMap, nil
}

func (model *Model) Log(migrationName string, batch int) error {
	var migration Model
	migration.Migration = migrationName
	migration.Batch = batch
	return model.Create(&migration).Error
}

func (model *Model) Delete(target *Model) error {
	return model.DB.Delete(&target).Error
}

func (model *Model) GetNextBatchNumber() (int, error) {
	last, err := model.GetLastBatchNumber()
	if err != nil {
		return 0, err
	}
	return last + 1, nil
}

func (model *Model) GetLastBatchNumber() (int, error) {
	var maxBatch int
	if err := model.Model(model).
		Select("IFNULL(MAX(batch), 0) max_match").
		Pluck("max_match", &maxBatch).
		Error; err != nil {
		return maxBatch, err
	}
	return maxBatch, nil
}

func NewDBMigration(db *gorm.DB) *Model {
	var migration Model
	if !db.Migrator().HasTable(&migration) {
		db.Migrator().CreateTable(&migration)
	}
	migration.DB = db
	return &migration
}
