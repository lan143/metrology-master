package entity

import (
	"github.com/lan143/metrology-master/internal/ha/enum"
)

type Base struct {
	Name           string              `json:"name,omitempty"`
	Device         Device              `json:"device"`
	EntityCategory enum.EntityCategory `json:"entity_category,omitempty"`
	ObjectID       string              `json:"object_id,omitempty"`
	UniqueID       string              `json:"unique_id"`
	ForceUpdate    bool                `json:"force_update,omitempty"`
}
