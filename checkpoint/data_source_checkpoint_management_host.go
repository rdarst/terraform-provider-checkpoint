package checkpoint

import (
	"fmt"
	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"reflect"
	"strconv"
)

func dataSourceManagementHost() *schema.Resource {
	return &schema.Resource{

		Read: datareadManagementHost,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Object name. Should be unique in the domain.",
			},
			"ipv4_address": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv4 address.",
			},
			"ipv6_address": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv6 address.",
			},
			"interfaces": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Host interfaces.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
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
						"ignore_warnings": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Apply changes ignoring warnings.",
						},
						"ignore_errors": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Apply changes ignoring errors. You won't be able to publish such a changes. If ignore-warnings flag was omitted - warnings will also be ignored.",
						},
						"color": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "black",
							Description: "Color of the object. Should be one of existing colors.",
						},
						"comments": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Comments string.",
						},
					},
				},
			},
			"nat_settings": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "NAT settings.",
				//Default: map[string]interface{}{"auto_rule":false},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_rule": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
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
			"host_servers": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Servers Configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_server": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Gets True if this server is a DNS Server.",
						},
						"mail_server": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Gets True if this server is a Mail Server.",
						},
						"web_server": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Gets True if this server is a Web Server.",
						},
						"web_server_config": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Web Server configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"additional_ports": &schema.Schema{
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Server additional ports.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"application_engines": &schema.Schema{
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Application engines of this web server.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"listen_standard_port": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "Whether server listens to standard port.",
									},
									"operating_system": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "other",
										Description: "Operating System.",
									},
									"protected_by": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "97aeb368-9aea-11d5-bd16-0090272ccb30",
										Description: "Network object which protects this server identified by the name or UID.",
									},
								},
							},
						},
					},
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
			"color": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Color of the object. Should be one of existing colors.",
				Default:     "black",
			},
			"comments": &schema.Schema{
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
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Collection of tag identifiers.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func datareadManagementHost(d *schema.ResourceData, m interface{}) error {

	client := m.(*checkpoint.ApiClient)

	payload := map[string]interface{}{
		"name": d.Get("name").(string),
	}

	showHostRes, err := client.ApiCall("show-host", payload, client.GetSessionID(), true, false)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	if !showHostRes.Success {
		return fmt.Errorf(showHostRes.ErrorMsg)
	}
	
	host := showHostRes.GetData()

	log.Println("Read Host - Show JSON = ", host)
	
	if v := host["uid"]; v != nil {
		_ = d.Set("uid", v)
		d.SetId(host["uid"].(string))
	}

	if v := host["name"]; v != nil {
		_ = d.Set("name", v)
	}

	if v := host["ipv4-address"]; v != nil {
		_ = d.Set("ipv4_address", v)
	}

	if v := host["ipv6-address"]; v != nil {
		_ = d.Set("ipv6_address", v)
	}

	if v := host["comments"]; v != nil {
		_ = d.Set("comments", v)
	}

	if v := host["color"]; v != nil {
		_ = d.Set("color", v)
	}

	//we are compromising here since we cant represent map inside map
	//see also  https://github.com/hashicorp/terraform-plugin-sdk/issues/155
	if host["interfaces"] != nil {

		interfacesList := host["interfaces"].([]interface{})

		if len(interfacesList) > 0 {

			var interfacesListToReturn []map[string]interface{}

			for i := range interfacesList {

				interfaceMap := interfacesList[i].(map[string]interface{})

				interfaceMapToAdd := make(map[string]interface{})

				if v, _ := interfaceMap["name"]; v != nil {
					interfaceMapToAdd["name"] = v
				}
				if v, _ := interfaceMap["subnet4"]; v != nil {
					interfaceMapToAdd["subnet4"] = v
				}
				if v, _ := interfaceMap["subnet6"]; v != nil {
					interfaceMapToAdd["subnet6"] = v
				}
				if v, _ := interfaceMap["mask-length4"]; v != nil {
					interfaceMapToAdd["mask_length4"] = v
				}
				if v, _ := interfaceMap["mask-length6"]; v != nil {
					interfaceMapToAdd["mask_length6"] = v
				}
				if v, _ := interfaceMap["color"]; v != nil {
					interfaceMapToAdd["color"] = v
				}
				if v, _ := interfaceMap["comments"]; v != nil {
					interfaceMapToAdd["comments"] = v
				}
				interfaceMapToAdd["ignore_errors"] = false
				interfaceMapToAdd["ignore_warnings"] = false

				interfacesListToReturn = append(interfacesListToReturn, interfaceMapToAdd)
			}

			_ = d.Set("interfaces", interfacesListToReturn)
		} else {
			_ = d.Set("interfaces", interfacesList)
		}
	} else {
		_ = d.Set("interfaces", nil)
	}

	if host["nat-settings"] != nil {

		natSettingsMap := host["nat-settings"].(map[string]interface{})

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

	if host["host-servers"] != nil {

		hostServersMap := host["host-servers"].(map[string]interface{})

		hostServersMapToReturn := make(map[string]interface{})

		if v, _ := hostServersMap["dns-server"]; v != nil {
			hostServersMapToReturn["dns_server"] = v
		}
		if v, _ := hostServersMap["mail-server"]; v != nil {
			hostServersMapToReturn["mail_server"] = v
		}
		if v, _ := hostServersMap["web-server"]; v != nil {
			hostServersMapToReturn["web_server"] = v
		}
		if v, ok := hostServersMap["web-server-config"]; ok {

			webServerConfigMap := v.(map[string]interface{})
			webServerConfigMapToReturn := make(map[string]interface{})

			if v, _ := webServerConfigMap["additional-ports"]; v != nil {
				webServerConfigMapToReturn["additional_ports"] = v
			}
			if v, _ := webServerConfigMap["application-engines"]; v != nil {
				webServerConfigMapToReturn["application_engines"] = v
			}
			if v, _ := webServerConfigMap["listen-standard-port"]; v != nil {
				webServerConfigMapToReturn["listen_standard_port"] = v
			}
			if v, _ := webServerConfigMap["operating-system"]; v != nil {
				webServerConfigMapToReturn["operating_system"] = v
			}
			if v, _ := webServerConfigMap["protected-by"]; v != nil {

				//show returned the uid, we want to set the name.
				payload := map[string]interface{}{
					"uid": v,
				}
				showProtectedByRes, err := client.ApiCall("show-object", payload, client.GetSessionID(), true, false)
				if err != nil || !showProtectedByRes.Success {
					if showProtectedByRes.ErrorMsg != "" {
						return fmt.Errorf(showProtectedByRes.ErrorMsg)
					}
					return fmt.Errorf(err.Error())
				}

				webServerConfigMapToReturn["protected_by"] = showProtectedByRes.GetData()["object"].(map[string]interface{})["name"]
			}
			hostServersMapToReturn["web_server_config"] = []interface{}{webServerConfigMapToReturn}
		}

		_ = d.Set("host_servers", []interface{}{hostServersMapToReturn})

	} else {
		_ = d.Set("host_servers", nil)
	}

	if host["groups"] != nil {
		groupsJson := host["groups"].([]interface{})
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

	if host["tags"] != nil {
		tagsJson := host["tags"].([]interface{})
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
