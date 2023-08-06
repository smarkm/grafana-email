package datasource

import (
	"log"

	"github.com/go-resty/resty/v2"
	"smark.freecoop.net/grafana-email/config"
)

type DashboardMeta struct {
	Dashboard Dashboard `json:"dashboard"`
}
type Dashboard struct {
	Panels []Panel `json:"panels"`
}
type Panel struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// /api/dashboards/uid/${uid} //uid=dashboardID
func DashboardPanels(orgID string, dID string) []int {
	apiKey := config.Instance.OrgAPIKeys[orgID]
	c := resty.New()
	rs, err := c.R().SetHeader("Authorization", "Bearer "+apiKey).
		SetResult(&DashboardMeta{}).
		Get(config.Instance.GrafanaUrl + "/api/dashboards/uid/" + dID)

	if err != nil {
		log.Println("Error: " + err.Error())
	} else {
		d := rs.Result().(*DashboardMeta)
		panels := make([]int, len(d.Dashboard.Panels))
		for i, p := range d.Dashboard.Panels {
			panels[i] = p.ID
		}
		return panels
	}

	if config.Instance.DebugModel {
		log.Println("Debug: ")
	}
	return nil
}
