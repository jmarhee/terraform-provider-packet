package packet

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var packetMutexKV = mutexkv.NewMutexKV()

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PACKET_AUTH_TOKEN", nil),
				Description: "The API auth key for API operations.",
			},
			"max_simultaneous_devices_create": {
				Type:        schema.TypeInt,
				Required:    false,
				DefaultFunc: schema.EnvDefaultFunc("PACKET_MAX_SIMULTANEOUS_DEVICES_CREATE", 6),
				Description: "Maximum number of devices to create simultaneously",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"packet_ip_block_ranges":     dataSourcePacketIPBlockRanges(),
			"packet_precreated_ip_block": dataSourcePacketPreCreatedIPBlock(),
			"packet_operating_system":    dataSourceOperatingSystem(),
			"packet_organization":        dataSourcePacketOrganization(),
			"packet_spot_market_price":   dataSourceSpotMarketPrice(),
			"packet_device":              dataSourcePacketDevice(),
			"packet_project":             dataSourcePacketProject(),
			"packet_spot_market_request": dataSourcePacketSpotMarketRequest(),
			"packet_volume":              dataSourcePacketVolume(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"packet_device":               resourcePacketDevice(),
			"packet_ssh_key":              resourcePacketSSHKey(),
			"packet_project_ssh_key":      resourcePacketProjectSSHKey(),
			"packet_project":              resourcePacketProject(),
			"packet_organization":         resourcePacketOrganization(),
			"packet_volume":               resourcePacketVolume(),
			"packet_volume_attachment":    resourcePacketVolumeAttachment(),
			"packet_reserved_ip_block":    resourcePacketReservedIPBlock(),
			"packet_ip_attachment":        resourcePacketIPAttachment(),
			"packet_spot_market_request":  resourcePacketSpotMarketRequest(),
			"packet_vlan":                 resourcePacketVlan(),
			"packet_bgp_session":          resourcePacketBGPSession(),
			"packet_port_vlan_attachment": resourcePacketPortVlanAttachment(),
			"packet_connect":              resourcePacketConnect(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	authToken := d.Get("auth_token").(string)
	maxDevicesCreate := d.Get("max_simultaneous_devices_create").(int)
	return GetProviderConfig(authToken, maxDevicesCreate), nil
}

var resourceDefaultTimeouts = &schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(60 * time.Minute),
	Update:  schema.DefaultTimeout(60 * time.Minute),
	Delete:  schema.DefaultTimeout(60 * time.Minute),
	Default: schema.DefaultTimeout(60 * time.Minute),
}
