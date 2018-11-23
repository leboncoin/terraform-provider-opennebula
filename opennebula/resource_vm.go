package opennebula

import (
	"encoding/xml"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/leboncoin/terraform-provider-opennebula/opennebula/encoding/template"
	m "github.com/leboncoin/terraform-provider-opennebula/opennebula/models"
)

type UserVm struct {
	Id             string          `xml:"ID"`
	Name           string          `xml:"NAME"`
	Uid            int             `xml:"UID"`
	Gid            int             `xml:"GID"`
	Uname          string          `xml:"UNAME"`
	Gname          string          `xml:"GNAME"`
	Permissions    *Permissions    `xml:"PERMISSIONS"`
	State          int             `xml:"STATE"`
	LcmState       int             `xml:"LCM_STATE"`
	VmTemplate     *VmTemplate     `xml:"TEMPLATE"`
	VmUserTemplate *VmUserTemplate `xml:"USER_TEMPLATE"`
}

type UserVms struct {
	UserVm []*UserVm `xml:"VM"`
}

type VmTemplate struct {
	XMLName xml.Name   `xml:"TEMPLATE"`
	Context *m.Context `xml:"CONTEXT"`
	Disk    []*m.Disk  `xml:"DISK"`
	Cpu     int        `xml:"CPU"`
	Vcpu    int        `xml:"VCPU"`
	Memory  int        `xml:"MEMORY"`
}

type VmUserTemplate struct {
	Labels string `xml:"LABELS"`
}

func (d *VmTemplate) MarshalTemplate() string {
	template, err := template.Marshal(d)
	if err != nil {
		fmt.Printf("failed to solve problem: %s\n", err)
	}
	return template
}

func resourceVm() *schema.Resource {
	return &schema.Resource{
		Create: resourceVmCreate,
		Read:   resourceVmRead,
		Exists: resourceVmExists,
		Update: resourceVmUpdate,
		Delete: resourceVmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the VM. If empty, defaults to 'templatename-<vmid>'",
			},
			"instance": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Final name of the VM instance",
			},
			"template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Id of the VM template to use. Either 'template_name' or 'template_id' is required",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "CPU count of the VM instance",
			},
			"vcpu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "VCPU count of the VM instance",
			},
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Memory in MB",
			},
			"labels": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Attatch labels on Vm instance",
			},
			"root_disk": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dev_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"image_id": {
							Type:     schema.TypeInt,
							ForceNew: true,
							Optional: true,
						},
						"image": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"driver": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"image_uid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"image_uname": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
				Set: m.ResourceDiskHash,
			},
			"disk": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dev_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"driver": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"format": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"target": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				Set: m.ResourceDiskHash,
			},
			"permissions": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Permissions for the template (in Unix format, owner-group-other, use-manage-admin)",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if len(value) != 3 {
						errors = append(errors, fmt.Errorf("%q has specify 3 permission sets: owner-group-other", k))
					}

					all := true
					for _, c := range strings.Split(value, "") {
						if c < "0" || c > "7" {
							all = false
						}
					}
					if !all {
						errors = append(errors, fmt.Errorf("Each character in %q should specify a Unix-like permission set with a number from 0 to 7", k))
					}

					return
				},
			},

			"uid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the user that will own the VM",
			},
			"gid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the group that will own the VM",
			},
			"uname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the user that will own the VM",
			},
			"gname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the group that will own the VM",
			},
			"state": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current state of the VM",
			},
			"lcmstate": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current LCM state of the VM",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current Ip of the VM",
			},
		},
	}
}

func resourceVmCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	resp, err := client.Call(
		"one.template.instantiate",
		d.Get("template_id"),
		d.Get("name"),
		false,
		buildUserTemplateAttributesString(d),
		false,
	)
	if err != nil {
		return err
	}

	d.SetId(resp)

	_, err = waitForVmState(d, meta, "running")
	if err != nil {
		return fmt.Errorf(
			"Error waiting for virtual machine (%s) to be in state RUNNING: %s", d.Id(), err)
	}

	if _, err = changePermissions(intId(d.Id()), permission(d.Get("permissions").(string)), client, "one.vm.chmod"); err != nil {
		return err
	}

	return resourceVmRead(d, meta)
}

func resourceVmRead(d *schema.ResourceData, meta interface{}) error {
	var vm *UserVm
	var vms *UserVms

	client := meta.(*Client)
	found := false
	name := d.Get("name").(string)
	if name == "" {
		name = d.Get("instance").(string)
	}

	// Try to find the vm by ID, if specified
	if d.Id() != "" {
		resp, err := client.Call("one.vm.info", intId(d.Id()))
		if err == nil {
			found = true
			if err = xml.Unmarshal([]byte(resp), &vm); err != nil {
				return err
			}
		} else {
			log.Printf("Could not find VM by ID %s", d.Id())
		}
	}

	// Otherwise, try to find the vm by (user, name) as the de facto compound primary key
	if d.Id() == "" || !found {
		resp, err := client.Call("one.vmpool.info", -3, -1, -1)
		if err != nil {
			return err
		}
		if err = xml.Unmarshal([]byte(resp), &vms); err != nil {
			return err
		}

		for _, v := range vms.UserVm {
			if v.Name == name {
				vm = v
				found = true
				break
			}
		}

		if !found || vm == nil {
			d.SetId("")
			log.Printf("Could not find vm with name %s for user %s", name, client.Username)
			return nil
		}
	}

	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": vm.VmTemplate.Context.IP,
	})

	disks := m.DisksState(vm.VmTemplate.Disk)

	d.SetId(vm.Id)
	d.Set("instance", vm.Name)
	d.Set("uid", vm.Uid)
	d.Set("gid", vm.Gid)
	d.Set("uname", vm.Uname)
	d.Set("gname", vm.Gname)
	d.Set("state", vm.State)
	d.Set("lcmstate", vm.LcmState)
	d.Set("cpu", vm.VmTemplate.Cpu)
	d.Set("vcpu", vm.VmTemplate.Vcpu)
	d.Set("memory", vm.VmTemplate.Memory)
	d.Set("labels", vm.VmUserTemplate.Labels)
	d.Set("root_disk", m.DiskMapping(disks["root"].([]map[string]interface{})))
	d.Set("disk", m.DiskMapping(disks["disk"].([]map[string]interface{})))
	d.Set("ip", vm.VmTemplate.Context.IP)
	d.Set("permissions", permissionString(vm.Permissions))

	return nil
}

func resourceVmExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	err := resourceVmRead(d, meta)
	// a terminated VM is in state 6 (DONE)
	if err != nil || d.Id() == "" || d.Get("state").(int) == 6 {
		return false, err
	}

	return true, nil
}

func resourceVmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	if d.HasChange("permissions") {
		resp, err := changePermissions(intId(d.Id()), permission(d.Get("permissions").(string)), client, "one.vm.chmod")
		if err != nil {
			return err
		}
		log.Printf("[INFO] Successfully updated VM %s\n", resp)
	}

	if d.HasChange("cpu") || d.HasChange("vcpu") || d.HasChange("memory") {
		resp, err := resizeVm(d, meta)
		if err != nil {
			return err
		}
		log.Printf("[INFO] Successfully resized VM %s\n", resp)
	}

	if d.HasChange("name") {
		resp, err := client.Call(
			"one.vm.rename",
			intId(d.Id()),
			fmt.Sprintf("%s", d.Get("name").(string)),
		)
		if err != nil {
			return err
		}
		log.Printf("[INFO] Successfully updated VM %s\n", resp)
	}

	if d.HasChange("labels") {
		resp, err := client.Call(
			"one.vm.update",
			intId(d.Id()),
			fmt.Sprintf("LABELS=\"%s\"", d.Get("labels").(string)),
			1,
		)
		if err != nil {
			return err
		}
		log.Printf("[INFO] Successfully updated conf of VM %s\n", resp)
	}

	return nil
}

func resourceVmDelete(d *schema.ResourceData, meta interface{}) error {
	err := resourceVmRead(d, meta)
	if err != nil || d.Id() == "" {
		return err
	}

	client := meta.(*Client)
	resp, err := client.Call("one.vm.action", "terminate-hard", intId(d.Id()))
	if err != nil {
		return err
	}

	_, err = waitForVmState(d, meta, "done")
	if err != nil {
		return fmt.Errorf(
			"Error waiting for virtual machine (%s) to be in state DONE: %s", d.Id(), err)
	}

	log.Printf("[INFO] Successfully terminated VM %s\n", resp)
	return nil
}

func resizeVm(d *schema.ResourceData, meta interface{}) (string, error) {

	// Shutdown VM
	client := meta.(*Client)
	// Todo: shutdown only if necessary
	client.Call("one.vm.action", "poweroff", intId(d.Id()))
	_, err := waitForVmState(d, meta, "poweroff")
	if err != nil {
		return "", fmt.Errorf(
			"Error waiting for virtual machine (%s) to be in state POWEROFF: %s", d.Id(), err)
	}

	// Resize VM
	resp, err := client.Call(
		"one.vm.resize",
		intId(d.Id()),
		strings.Join([]string{
			buildBaseAttribute(d, "cpu"),
			buildBaseAttribute(d, "vcpu"),
			buildBaseAttribute(d, "memory"),
		}, "\n"),
		true)
	if err != nil {
		return "", fmt.Errorf(
			"Error during resize of virtual machine (%s): %s", d.Id(), err)
	}

	// Resume VM
	_, err = client.Call("one.vm.action", "resume", intId(d.Id()))
	if err != nil {
		return "", err
	}
	_, err = waitForVmState(d, meta, "running")
	if err != nil {
		return "", fmt.Errorf(
			"Error waiting for virtual machine (%s) to be in state RUNNING: %s", d.Id(), err)
	}

	return resp, nil
}

func waitForVmState(d *schema.ResourceData, meta interface{}, state string) (interface{}, error) {
	var vm *UserVm
	client := meta.(*Client)

	log.Printf("Waiting for VM (%s) to be in state Done", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending: []string{"anythingelse"},
		Target:  []string{state},
		Refresh: func() (interface{}, string, error) {
			log.Println("Refreshing VM state...")
			if d.Id() != "" {
				resp, err := client.Call("one.vm.info", intId(d.Id()))
				if err == nil {
					if err = xml.Unmarshal([]byte(resp), &vm); err != nil {
						return nil, "", fmt.Errorf("Couldn't fetch VM state: %s", err)
					}
				} else {
					return nil, "", fmt.Errorf("Could not find VM by ID %s", d.Id())
				}
			}
			log.Printf("VM is currently in state %v and in LCM state %v", vm.State, vm.LcmState)
			if vm.State == 3 && vm.LcmState == 3 {
				return vm, "running", nil
			} else if vm.State == 6 {
				return vm, "done", nil
			} else if vm.State == 8 {
				return vm, "poweroff", nil
			} else {
				return nil, "anythingelse", nil
			}
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	return stateConf.WaitForState()
}

func buildUserTemplateAttributesString(d *schema.ResourceData) string {
	disks, _ := m.ReadDiskFromConfig(d)
	diskTpl := m.DisksTemplate(disks)
	data := strings.Join([]string{
		buildBaseAttribute(d, "labels"),
		buildBaseAttribute(d, "cpu"),
		buildBaseAttribute(d, "vcpu"),
		buildBaseAttribute(d, "memory"),
		diskTpl}, "\n")
	log.Printf(data)
	return data
}

func buildBaseAttribute(d *schema.ResourceData, key string) string {
	if d.HasChange(key) {
		value := d.Get(key)
		return getBaseAttributeFromKeyValue(key, value)
	}
	return ""
}

func getBaseAttributeFromKeyValue(key string, value interface{}) string {
	format := "%s = \"%d\""
	if reflect.TypeOf(value).String() == "string" {
		format = "%s = \"%s\""
	}
	return fmt.Sprintf(format, strings.ToUpper(key), value)
}
