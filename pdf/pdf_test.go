package pdf

import (
	"testing"

	"smark.freecoop.net/grafana-email/config"
	"smark.freecoop.net/grafana-email/datasource"
)

func TestPDF(t *testing.T) {
	config.Init("../config.json")
	pd := InitPDF()
	bytes := datasource.PanelImage("1", "Gfgpou3Vk", "4", nil)
	InsertImage("test", pd, bytes, 70)
	bytes = datasource.PanelImage("1", "Gfgpou3Vk", "2", nil)
	InsertImageInNewPage("test2", pd, bytes)

	pd.OutputFileAndClose("t.pdf")
}
