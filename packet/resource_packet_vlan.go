package packet

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketVlan() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketVolumeCreate,
		Read:   resourcePacketVolumeRead,
		Update: resourcePacketVolumeUpdate,
		Delete: resourcePacketVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"facility": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vxlan": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourcePacketVlanCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*packngo.Client)
	createRequest := &packngo.VirtualNetworkCreateRequest{
		ProjectID:   d.Get("project_id").(string),
		Description: d.Get("description").(string),
		Facility:    d.Get("facility").(string),
	}
	vlan, _, err := c.ProjectVirtualNetworks.Create(createRequest)
	if err != nil {
		return friendlyError(err)
	}
	d.SetId(vlan.ID)
	return resourcePacketVlanRead(d, meta)
}

func resourcePacketVlanRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*packngo.Client)

	vlan, _, err := c.ProjectVirtualNetworks.Get(d.Id(),
		&packngo.GetOptions{Includes: []string{"assigned_to"}})
	if err != nil {
		return friendlyError(err)

		if isNotFound(err) {
			d.SetId("")
			return nil
		}

		return err
	}
	d.Set("description", vlan.Description)
	d.Set("project_id", vlan.Project.ID)
	d.Set("vxlan", vlan.VXLAN)
	return nil
}

func resourcePacketVlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	_, err := client.ProjectVirtualNetworks.Delete(d.Id())
	if err != nil {
		return friendlyError(err)
	}

	d.SetId("")
	return nil
}
