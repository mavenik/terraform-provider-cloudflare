package cloudflare

import (
	"fmt"
	"log"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceCloudflareDNSRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudflareDNSRecordsRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"zone_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"modified_on": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudflareDNSRecordsRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading DnsRecords")
	client := meta.(*cloudflare.API)

	zoneId := d.Get("zone_id").(string)

	recordFilter := cloudflare.DNSRecord{
		Type:    d.Get("type").(string),
		Name:    d.Get("name").(string),
		Content: d.Get("content").(string),
	}

	dns_records, err := client.DNSRecords(zoneId, recordFilter)
	if err != nil {
		return fmt.Errorf("error listing DNS Records: %s", err)
	}

	dnsRecordDetails := make([]interface{}, 0)
	for _, v := range dns_records {
		dnsRecordDetails = append(dnsRecordDetails, map[string]interface{}{
			"id":          v.ID,
			"name":        v.Name,
			"type":        v.Type,
			"content":     v.Content,
			"zone_id":     v.ZoneID,
			"zone_name":   v.ZoneName,
			"ttl":         v.TTL,
			"created_on":  v.CreatedOn.Format(time.RFC1123),
			"modified_on": v.ModifiedOn.Format(time.RFC1123),
		})
	}

	err = d.Set("dns_records", dnsRecordDetails)
	if err != nil {
		return fmt.Errorf("Error setting dns_records: %s", err)
	}

	d.SetId(time.Now().UTC().String())

	return nil
}
