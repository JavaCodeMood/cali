package services

import (
	"github.com/google/uuid"
	"github.com/jiangmitiao/cali/app/models"
	"sync"
	"time"
)

type SysStatusService struct {
	lock *sync.Mutex
}

func (service SysStatusService) Get(key string) models.SysStatus {
	sysStatus := models.SysStatus{}
	engine.Where("ikey = ?", key).Get(&sysStatus)
	return sysStatus
}

func (service SysStatusService) QuerySysStatus(limit, start int) []models.SysStatus {
	sysStatus := make([]models.SysStatus, 0)
	engine.Desc("updated").Limit(limit, start).Find(&sysStatus)
	return sysStatus
}

func (service SysStatusService) UpdateStatus(sysStatus models.SysStatus) bool {
	sysStatus.UpdatedAt = time.Now().Unix()
	_, err := engine.Id(sysStatus.Id).Cols("ikey", "value", "updated").Update(sysStatus)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (service SysStatusService) AddSysStatus(sysStatus models.SysStatus) bool {
	sysStatus.Id = uuid.New().String()
	sysStatus.UpdatedAt = time.Now().Unix()
	sysStatus.CreatedAt = time.Now().Unix()
	if _, err := engine.InsertOne(sysStatus); err == nil {
		return true
	} else {
		return false
	}
}

func (service SysStatusService) DeleteSysStatus(sysStatus models.SysStatus) bool {
	if _, err := engine.Where("id = ?", sysStatus.Id).Delete(models.SysStatus{}); err == nil {
		return true
	} else {
		return false
	}
}
