package types

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func JsonData(value any) datatypes.JSON {
	buf, _ := json.Marshal(value)
	return datatypes.JSON(buf)
}
