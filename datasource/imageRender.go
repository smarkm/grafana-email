package datasource

import (
	"log"

	"github.com/go-resty/resty/v2"
	"smark.freecoop.net/grafana-email/config"
)

func PanelImage(orgID string, dID string, panelID string, vars map[string]string) []byte {
	apiKey := config.Instance.OrgAPIKeys[orgID]
	c := resty.New()
	req := c.R().SetHeader("Authorization", "Bearer "+apiKey).
		SetQueryParams(map[string]string{
			"orgId":   orgID,
			"panelId": panelID,
		})
	for k, v := range vars {
		req.SetQueryParam(k, string(v))
	}
	rs, err := req.Get(config.Instance.GrafanaUrl + "/render/d-solo/" + dID)

	if err != nil {
		log.Println("Error: " + err.Error())
	} else {
		d := rs.Body()
		return d
	}

	if config.Instance.DebugModel {
		log.Println("Debug: ")
	}
	return nil
}
