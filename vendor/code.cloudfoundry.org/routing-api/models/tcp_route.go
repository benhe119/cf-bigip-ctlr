package models

import (
	"fmt"
	"time"

	"github.com/nu7hatch/gouuid"
)

type TcpRouteMapping struct {
	Model
	ExpiresAt time.Time `json:"-"`
	TcpMappingEntity
}

type TcpMappingEntity struct {
	RouterGroupGuid string `json:"router_group_guid"`
	HostPort        uint16 `gorm:"not null; unique_index:idx_tcp_route; type:int" json:"backend_port"`
	HostIP          string `gorm:"not null; unique_index:idx_tcp_route" json:"backend_ip"`
	ExternalPort    uint16 `gorm:"not null; unique_index:idx_tcp_route; type: int" json:"port"`
	ModificationTag `json:"modification_tag"`
	TTL             *int `json:"ttl,omitempty"`
}

func (TcpRouteMapping) TableName() string {
	return "tcp_routes"
}

func NewTcpRouteMappingWithModel(tcpMapping TcpRouteMapping) (TcpRouteMapping, error) {
	guid, err := uuid.NewV4()
	if err != nil {
		return TcpRouteMapping{}, err
	}

	m := Model{Guid: guid.String()}
	return TcpRouteMapping{
		ExpiresAt:        time.Now().Add(time.Duration(*tcpMapping.TTL) * time.Second),
		Model:            m,
		TcpMappingEntity: tcpMapping.TcpMappingEntity,
	}, nil
}

func NewTcpRouteMapping(routerGroupGuid string, externalPort uint16, hostIP string, hostPort uint16, ttl int) TcpRouteMapping {
	mapping := TcpMappingEntity{
		RouterGroupGuid: routerGroupGuid,
		ExternalPort:    externalPort,
		HostPort:        hostPort,
		HostIP:          hostIP,
		TTL:             &ttl,
	}
	return TcpRouteMapping{
		TcpMappingEntity: mapping,
	}
}

func NewTcpRouteMappingWithModificationTag(
	routerGroupGuid string,
	externalPort uint16,
	hostIP string,
	hostPort uint16,
	ttl int,
	modTag ModificationTag,
) TcpRouteMapping {
	mapping := NewTcpRouteMapping(routerGroupGuid, externalPort, hostIP, hostPort, ttl)
	mapping.ModificationTag = modTag

	return mapping
}

func (m TcpRouteMapping) String() string {
	return fmt.Sprintf("%s:%d<->%s:%d", m.RouterGroupGuid, m.ExternalPort, m.HostIP, m.HostPort)
}

func (m TcpRouteMapping) Matches(other TcpRouteMapping) bool {
	return m.RouterGroupGuid == other.RouterGroupGuid &&
		m.ExternalPort == other.ExternalPort &&
		m.HostIP == other.HostIP &&
		m.HostPort == other.HostPort &&
		*m.TTL == *other.TTL
}

func (t *TcpRouteMapping) SetDefaults(maxTTL int) {
	// default ttl if not present
	// TTL is a pointer to a uint16 so that we can
	// detect if it's present or not (i.e. nil or 0)
	if t.TTL == nil {
		t.TTL = &maxTTL
	}
}