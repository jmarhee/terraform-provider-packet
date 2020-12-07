package packet

import (
	"fmt"
	"log"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketIPAttachment() *schema.Resource {
	ipAttachmentSchema := packetIPResourceComputedFields()
	ipAttachmentSchema["device_id"] = &schema.Schema{
		Type:     schema.TypeString,
		ForceNew: true,
		Required: true,
	}
	ipAttachmentSchema["cidr_notation"] = &schema.Schema{
		Type:     schema.TypeString,
		ForceNew: true,
		Required: true,
	}
	return &schema.Resource{
		Create: resourcePacketIPAttachmentCreate,
		Read:   resourcePacketIPAttachmentRead,
		Delete: resourcePacketIPAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: ipAttachmentSchema,
	}
}

func resourcePacketIPAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	deviceID := d.Get("device_id").(string)
	ipa := d.Get("cidr_notation").(string)

	req := packngo.AddressStruct{Address: ipa}

	assignment, _, err := client.DeviceIPs.Assign(deviceID, &req)
	if err != nil {
		return fmt.Errorf("error assigning address %s to device %s: %s", ipa, deviceID, err)
	}

	d.SetId(assignment.ID)

	return resourcePacketIPAttachmentRead(d, meta)
}

func resourcePacketIPAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	assignment, _, err := client.DeviceIPs.Get(d.Id(), nil)
	if err != nil {
		err = friendlyError(err)

		// If the IP attachment was already destroyed, mark as succesfully gone.
		if isNotFound(err) {
			log.Printf("[WARN] IP attachment (%q) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	g := false
	if assignment.Global != nil {
		g = *(assignment.Global)
	}

	d.SetId(assignment.ID)

	return setMap(d, map[string]interface{}{
		"address":        assignment.Address,
		"gateway":        assignment.Gateway,
		"network":        assignment.Network,
		"netmask":        assignment.Netmask,
		"address_family": assignment.AddressFamily,
		"cidr":           assignment.CIDR,
		"public":         assignment.Public,
		"management":     assignment.Management,
		"manageable":     assignment.Manageable,
		"global":         g,
		"device_id":      path.Base(assignment.AssignedTo.Href),
		"cidr_notation":  fmt.Sprintf("%s/%d", assignment.Network, assignment.CIDR),
	})
}

func resourcePacketIPAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	_, err := client.DeviceIPs.Unassign(d.Id())
	if err != nil {
		return friendlyError(err)
	}

	d.SetId("")
	return nil
}
