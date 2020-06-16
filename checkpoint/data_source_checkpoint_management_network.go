package checkpoint

import (
	"fmt"
	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"reflect"
	"strconv"
)

func dataSourceManagementNetwork() *schema.Resource {
	return &schema.Resource{
		
		Read:   datareadManagementNetwork,
		
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Object name. Should be unique in the domain.",
			},
			"subnet4": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv4 network address.",
			},
			"subnet6": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv6 network address.",
			},
			"mask_length4": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "IPv4 network mask length.",
			},
			"mask_length6": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "IPv6 network mask length.",
			},
			"nat_settings": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "NAT settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_rule": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to add automatic address translation rules.",
						},
						"ipv4_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IPv4 address.",
						},
						"ipv6_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IPv6 address.",
						},
						"hide_behind": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Hide behind method. This parameter is not required in case \"method\" parameter is \"static\".",
						},
						"install_on": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Which gateway should apply the NAT translation.",
						},
						"method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "NAT translation method.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Collection of tag identifiers.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"broadcast": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allow broadcast address inclusion.",
				Default:     "allow",
			},
			"color": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Color of the object. Should be one of existing colors.",
				Default:     "black",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comments string.",
			},
			"groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Collection of group identifiers.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ignore_warnings": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Apply changes ignoring warnings.",
				Default:     false,
			},
			"ignore_errors": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Apply changes ignoring errors. You won't be able to publish such a changes. If ignore-warnings flag was omitted - warnings will also be ignored.",
				Default:     false,
			},
		},
	}
}

func datareadManagementNetwork(d *schema.ResourceData, m interface{}) error {

	client := m.(*checkpoint.ApiClient)

	payload := map[string]interface{}{
		"name": d.Get("name").(string),
	}

	showNetworkRes, err := client.ApiCall("show-network", payload, client.GetSessionID(), true, false)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	if !showNetworkRes.Success {
		return fmt.Errorf(showNetworkRes.ErrorMsg)
	}

	network := showNetworkRes.GetData()

	if v := network["uid"]; v != nil {
		_ = d.Set("uid", v)
		d.SetId(network["uid"].(string))
	}

	if v := network["name"]; v != nil {
		_ = d.Set("name", v)
	}

	if v := network["subnet4"]; v != nil {
		_ = d.Set("subnet4", v)
	}

	if v := network["subnet6"]; v != nil {
		_ = d.Set("subnet6", v)
	}

	if v := network["mask-length4"]; v != nil {
		_ = d.Set("mask_length4", v)
	}

	if v := network["mask-length6"]; v != nil {
		_ = d.Set("mask_length6", v)
	}

	if v := network["subnet-mask"]; v != nil {
		_ = d.Set("subnet_mask", v)
	}

	if v := network["broadcast"]; v != nil {
		_ = d.Set("broadcast", v)
	}

	if v := network["comments"]; v != nil {
		_ = d.Set("comments", v)
	}

	if v := network["color"]; v != nil {
		_ = d.Set("color", v)
	}

	if network["nat-settings"] != nil {

		natSettingsMap := network["nat-settings"].(map[string]interface{})

		natSettingsMapToReturn := make(map[string]interface{})

		if v, _ := natSettingsMap["auto-rule"]; v != nil {
			natSettingsMapToReturn["auto_rule"] = strconv.FormatBool(v.(bool))
		}

		if v, _ := natSettingsMap["ipv4-address"]; v != "" && v != nil {
			natSettingsMapToReturn["ipv4_address"] = v
		}

		if v, _ := natSettingsMap["ipv6-address"]; v != "" && v != nil {
			natSettingsMapToReturn["ipv6_address"] = v
		}

		if v, _ := natSettingsMap["hide-behind"]; v != nil {
			natSettingsMapToReturn["hide_behind"] = v
		}

		if v, _ := natSettingsMap["install-on"]; v != nil {
			natSettingsMapToReturn["install_on"] = v
		}

		if v, _ := natSettingsMap["method"]; v != nil {
			natSettingsMapToReturn["method"] = v
		}

		_, natSettingInConf := d.GetOk("nat_settings")
		defaultNatSettings := map[string]interface{}{"auto_rule": "false"}

		if reflect.DeepEqual(defaultNatSettings, natSettingsMapToReturn) && !natSettingInConf {
			_ = d.Set("nat_settings", map[string]interface{}{})
		} else {
			_ = d.Set("nat_settings", natSettingsMapToReturn)
		}

	} else {
		_ = d.Set("nat_settings", nil)
	}

	if network["groups"] != nil {
		groupsJson := network["groups"].([]interface{})
		groupsIds := make([]string, 0)
		if len(groupsJson) > 0 {
			// Create slice of group names
			for _, group := range groupsJson {
				group := group.(map[string]interface{})
				groupsIds = append(groupsIds, group["name"].(string))
			}
		}
		_ = d.Set("groups", groupsIds)
	} else {
		_ = d.Set("groups", nil)
	}

	if network["tags"] != nil {
		tagsJson := network["tags"].([]interface{})
		var tagsIds = make([]string, 0)
		if len(tagsJson) > 0 {
			// Create slice of tag names
			for _, tag := range tagsJson {
				tag := tag.(map[string]interface{})
				tagsIds = append(tagsIds, tag["name"].(string))
			}
		}
		_ = d.Set("tags", tagsIds)
	} else {
		_ = d.Set("tags", nil)
	}

	return nil
}

