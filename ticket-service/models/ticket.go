package models

type tickets struct{

	ID          uint64 `gorm:"primaryKey" json:"id"`
	Priority string `json:" priority"`
	Description string `json:"description"`
	RaisedBy uint64  `json:"raisedby"`
	AssignedTo uint64 `json:"assignedto"`
	Status string `json:"status"`
	CreatedTime string `json:"createdtime"`
	RemovedTime string `json:"raisedtime"`
}	