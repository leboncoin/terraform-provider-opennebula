package opennebula

import (
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

type UserVm struct {
	Id          string       `xml:"ID"`
	Name        string       `xml:"NAME"`
	Uid         int          `xml:"UID"`
	Gid         int          `xml:"GID"`
	Uname       string       `xml:"UNAME"`
	Gname       string       `xml:"GNAME"`
	Permissions *Permissions `xml:"PERMISSIONS"`
	State       int          `xml:"STATE"`
	LcmState    int          `xml:"LCM_STATE"`
	VmTemplate  *VmTemplate  `xml:"TEMPLATE"`
}

type UserVms struct {
	UserVm []*UserVm `xml:"VM"`
}

type VmTemplate struct {
	Context *Context `xml:"CONTEXT"`
	Nic     *Nic     `xml:"NIC"`
	Disk    *Disk    `xml:"DISK"`
	Cpu     int      `xml:"CPU"`
	Vcpu    int      `xml:"VCPU"`
	Memory  int      `xml:"MEMORY"`
}

type Context struct {
	IP string `xml:"ETH0_IP"`
}

type Nic struct {
	Network             string `xml:"NETWORK"`
	NetworkUname        string `xml:"NETWORK_UNAME"`
	NetworkSearchDomain string `xml:"SEARCH_DOMAIN"`
	SecurityGroupId     int    `xml:"SECURITY_GROUPS"`
}

type Disk struct {
	Image       string `xml:"IMAGE"`
	Size        int    `xml:"SIZE"`
	ImageDriver string `xml:"DRIVER"`
	ImageUname  string `xml:"IMAGE_UNAME"`
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
				Required:    true,
				ForceNew:    true,
				Description: "CPU count of the VM instance",
			},
			"vcpu": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "VCPU count of the VM instance",
			},
			"memory": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Memory in MB",
			},
			"image": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Image Name",
			},
			"image_uname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Image Owner",
			},
			"image_driver": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Image Driver",
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "VM Size in MB",
			},
			"network": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network Name",
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Optional IP Addr. for Network",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					// todo: maybe better error msgs

					parts := strings.Split(value, ".")
					if len(parts) < 4 {
						errors = append(errors, fmt.Errorf("%q doesn't consists of four octets", k))
					}

					for _, x := range parts {
						if i, err := strconv.Atoi(x); err == nil {
							if i < 0 || i > 255 {
								errors = append(errors, fmt.Errorf("%q octets are not in a valid range ", k))
							}
						} else {
							errors = append(errors, fmt.Errorf("%q not an valid ip format", k)) //todo: error msg
						}
					}
					return
				},
			},
			"network_uname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network Owner",
			},
			"network_search_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network Search Domain",
			},
			"security_group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Security Group ID",
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
		//todo: maybe use backticks
		fmt.Sprintf(
			"CPU = \"%d\"\n"+
				"VCPU = \"%d\"\n"+
				"MEMORY = \"%d\"\n "+
				"DISK=[\n"+
				"  IMAGE=\"%s\",\n"+
				"  SIZE=\"%d\",\n"+
				"  IMAGE_UNAME=\"%s\",\n"+
				"  DRIVER=\"%s\"]\n"+
				"NIC=[\n"+
				"  NETWORK=\"%s\",\n"+
				"  NETWORK_UNAME=\"%s\",\n"+
				"  SEARCH_DOMAIN=\"%s\",\n"+
				"  SECURITY_GROUP=\"%d\",\n"+
				"  IP=\"%s\"]",
			d.Get("cpu"),
			d.Get("vcpu"),
			d.Get("memory"),
			d.Get("image"),
			d.Get("size"),
			d.Get("image_uname"),
			d.Get("image_driver"),
			d.Get("network"),
			d.Get("network_uname"),
			d.Get("network_search_domain"),
			d.Get("security_group_id"),
			d.Get("ip")),
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
	d.Set("image", vm.VmTemplate.Disk.Image)
	d.Set("size", vm.VmTemplate.Disk.Size)
	d.Set("image_driver", vm.VmTemplate.Disk.ImageDriver)
	d.Set("image_uname", vm.VmTemplate.Disk.ImageUname)
	d.Set("network_uname", vm.VmTemplate.Nic.NetworkUname)
	d.Set("network_search_domain", vm.VmTemplate.Nic.NetworkSearchDomain)
	d.Set("security_group_id", vm.VmTemplate.Nic.SecurityGroupId)
	d.Set("network", vm.VmTemplate.Nic.Network)
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

	if d.HasChange("size") {
		resp, err := client.Call(
			"one.vm.diskresize",
			intId(d.Id()),
			0,
			fmt.Sprintf("%d", d.Get("size").(int)),
		)
		if err != nil {
			return err
		}
		log.Printf("[INFO] Successfully updated VM %s\n", resp)
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
