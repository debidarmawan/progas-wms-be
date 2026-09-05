package enum

type MasterItemType string

const (
	MasterItemTypeGas    MasterItemType = "gas"
	MasterItemTypeLiquid MasterItemType = "liquid"
	MasterItemTypeMix    MasterItemType = "mix"
)

func (t MasterItemType) IsValid() bool {
	switch t {
	case MasterItemTypeGas, MasterItemTypeLiquid, MasterItemTypeMix:
		return true
	default:
		return false
	}
}
