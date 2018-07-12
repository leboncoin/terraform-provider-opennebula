package opennebula

import (
	"encoding/xml"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOpennebulaTemplateId() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceOpennebulaTemplateIdRead,

		Schema: map[string]*schema.Schema{
			"template_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceOpennebulaTemplateIdRead(d *schema.ResourceData, meta interface{}) error {
	var tmpl *UserTemplate
	var tmpls *UserTemplates

	client := meta.(*Client)
	found := false

	resp, err := client.Call("one.templatepool.info", -2, -1, -1)
	if err != nil {
		return err
	}

	if err = xml.Unmarshal([]byte(resp), &tmpls); err != nil {
		return err
	}

	for _, t := range tmpls.UserTemplate {
		if t.Name == d.Get("template_name").(string) {
			tmpl = t
			found = true
			break
		}
	}

	if !found || tmpl == nil {
		d.SetId("")
		log.Printf("Could not find template with template_name %s for user %s", d.Get("template_name").(string), client.Username)
		return fmt.Errorf("Could not find template with name: %s for user %s", d.Get("template_name").(string), client.Username)
	}

	d.SetId(strconv.Itoa(tmpl.Id))

	return nil
}
