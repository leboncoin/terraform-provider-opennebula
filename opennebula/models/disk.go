package models

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/leboncoin/terraform-provider-opennebula/opennebula/encoding/template"
)

type Disk struct {
	DevPrefix   string `template:"DEV_PREFIX" xml:"DEV_PREFIX"`
	DiskID      int    `template:"DISK_ID" xml:"DISK_ID"`
	DiskType    string `template:"DISK_TYPE" xml:"DISK_TYPE"`
	ImageID     int    `template:"IMAGE_ID" xml:"IMAGE_ID"`
	Image       string `template:"IMAGE" xml:"IMAGE"`
	Size        int    `template:"SIZE" xml:"SIZE"`
	ImageDriver string `template:"DRIVER" xml:"DRIVER"`
	ImageFormat string `template:"FORMAT" xml:"FORMAT"`
	ImageUID    string `template:"IMAGE_UID" xml:"IMAGE_UID"`
	ImageUname  string `template:"IMAGE_UNAME" xml:"IMAGE_UNAME"`
	Target      string `template:"TARGET" xml:"TARGET"`
	Type        string `template:"TYPE" xml:"TYPE"`
}

func (d *Disk) marshalTemplate() string {
	template, err := template.MarshalWrap(d)
	if err != nil {
		fmt.Printf("failed to solve problem: %s\n", err)
	}
	return template
}

func (d *Disk) State() map[string]interface{} {
	val := reflect.ValueOf(d).Elem()
	attributes := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if (valueField.Kind().String() == "int") && valueField.Interface().(int) != 0 ||
			(valueField.Kind().String() == "string" && valueField.Interface().(string) != "") {
			attributes[strings.ToLower(tag.Get("xml"))] = valueField.Interface()
		}
	}
	return attributes

}

func ReadDiskFromConfig(
	d *schema.ResourceData) ([]*Disk, error) {
	disks := make([]*Disk, 0)
	if v, ok := d.GetOk("root_disk"); ok {
		vL := v.(*schema.Set).List()
		for _, v := range vL {
			bd := v.(map[string]interface{})
			disk := &Disk{}

			if v, ok := bd["image"].(string); ok && v != "" {
				disk.Image = v
			} else if v, ok := bd["image_id"].(int); ok && v != 0 {
				disk.ImageID = v
			} else {
				return nil, fmt.Errorf("image or image_id is mandatory")
			}

			if v, ok := bd["dev_prefix"].(string); ok && v != "" {
				disk.DevPrefix = v
			}
			if v, ok := bd["disk_id"].(int); ok && v != 0 {
				disk.DiskID = v
			}
			if v, ok := bd["disk_type"].(string); ok && v != "" {
				disk.DiskType = v
			}
			if v, ok := bd["size"].(int); ok && v != 0 {
				disk.Size = v
			}
			if v, ok := bd["driver"].(string); ok && v != "" {
				disk.ImageDriver = v
			}
			if v, ok := bd["format"].(string); ok && v != "" {
				disk.ImageFormat = v
			}
			if v, ok := bd["image_uid"].(string); ok && v != "" {
				disk.ImageUID = v
			}
			if v, ok := bd["image_uname"].(string); ok && v != "" {
				disk.ImageUname = v
			}

			disks = append(disks, disk)
		}
	}

	if v, ok := d.GetOk("disk"); ok {
		vL := v.(*schema.Set).List()
		for _, v := range vL {
			bd := v.(map[string]interface{})
			disk := &Disk{
				ImageFormat: bd["format"].(string),
			}

			if v, ok := bd["dev_prefix"].(string); ok && v != "" {
				disk.DevPrefix = v
			}
			if v, ok := bd["disk_id"].(int); ok && v != 0 {
				disk.DiskID = v
			}
			if v, ok := bd["size"].(int); ok && v != 0 {
				disk.Size = v
			}
			if v, ok := bd["driver"].(string); ok && v != "" {
				disk.ImageDriver = v
			}
			if v, ok := bd["format"].(string); ok && v != "" {
				disk.ImageFormat = v
			}
			if v, ok := bd["target"].(string); ok && v != "" {
				disk.Target = v
			}
			if v, ok := bd["type"].(string); ok && v != "" {
				disk.Type = v
			}

			disks = append(disks, disk)
		}
	}
	return disks, nil
}

func DisksTemplate(d []*Disk) string {
	attributes := make([]string, 0)
	for _, disk := range d {
		attributes = append(attributes, disk.marshalTemplate())
	}
	return strings.Join(attributes, "\n")
}

func DisksState(d []*Disk) map[string]interface{} {
	disks := make(map[string]interface{})
	disks["disk"] = make([]map[string]interface{}, 0)
	disks["root"] = make([]map[string]interface{}, 0)
	for _, disk := range d {
		diskState := disk.State()
		if disk.Image != "" || disk.ImageID != 0 {
			disks["root"] = append(disks["root"].([]map[string]interface{}), diskState)
		} else {
			disks["disk"] = append(disks["disk"].([]map[string]interface{}), diskState)
		}
	}
	return disks
}

func DiskMapping(m []map[string]interface{}) *schema.Set {
	s := &schema.Set{
		F: ResourceDiskHash,
	}
	for _, v := range m {
		s.Add(v)
	}
	return s
}

func ResourceDiskHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if val, ok := m["image"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", val.(string)))
	}
	if val, ok := m["image_id"]; ok {
		if val != 0 {
			buf.WriteString(fmt.Sprintf("%d-", val.(int)))
		}
	}
	if val, ok := m["target"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", val.(string)))
	}
	return hashcode.String(buf.String())
}
