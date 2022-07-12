package models

type SwagCreateGateway struct {
	AreaID    string `json:"areaId"`
	GatewayID string `json:"gatewayId"`
	Name      string `json:"name"`
}

type SwagUpateGateway struct {
	SwagCreateGateway
}

type SwagCreateArea struct {
	Name    string `json:"name"`
	Manager string `json:"manager"`
}

type SwagUpdateArea struct {
	GormModel
	SwagCreateArea
}
