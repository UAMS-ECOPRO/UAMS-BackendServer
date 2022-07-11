package models

type SwagCreateGateway struct {
	AreaID    uint   `json:"areaId"`
	GatewayID string `json:"gatewayId"`
	Name      string `json:"name"`
}

type SwagUpateGateway struct {
	GormModel
	SwagCreateGateway
}

type SwagCreateArea struct {
	Gateway Gateway `json:"gateway"`
	Name    string  `json:"name"`
	Manager string  `json:"manager"`
}

type SwagUpdateArea struct {
	GormModel
	SwagCreateArea
}
