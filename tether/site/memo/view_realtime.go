package memo

import (
	"html"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/tether/components/composite/monitor"
	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// RenderRealtime builds the memoised real-time dashboard. Each chart
// is wrapped in node.Memo with its Versioned key so only the chart
// whose data changed re-renders on each tick.
func RenderRealtime(s RealtimeState) node.Node {
	return cpage.New(
		panel.Card(
			"Memoised System Monitor",
			"Same live Go runtime metrics as the Real-time Dashboard, "+
				"but each chart is wrapped in node.Memo with a Versioned "+
				"key. On each tick, all three metrics change so all charts "+
				"re-render. The memoisation benefit shows when only some "+
				"metrics change - unchanged charts are skipped entirely.",
			"node.Memo · tether.Versioned · Memo: true · go-echarts", panel.WS|panel.SSE,
			monitor.Charts(
				div.New(
					node.Memo(s.CPUPercent.Version(), func() node.Node {
						return chartDiv("memocpu", "CPU (%)", "#ee6666", toLineData(s.CPUPercent.Val))
					}),
				).Dynamic("chart-cpu"),
				div.New(
					node.Memo(s.HeapMB.Version(), func() node.Node {
						return chartDiv("memoheap", "Heap (MB)", "#5470c6", toLineData(s.HeapMB.Val))
					}),
				).Dynamic("chart-heap"),
				div.New(
					node.Memo(s.Goroutines.Version(), func() node.Node {
						return chartDiv("memogoroutines", "Goroutines", "#91cc75", intsToLineData(s.Goroutines.Val))
					}),
				).Dynamic("chart-goroutines"),
			),
		),
	)
}

func chartDiv(id, titleText, colour string, data []opts.LineData) node.Node {
	option := buildChartOption(id, titleText, colour, data)
	el := monitor.Chart(id)
	el.SetAttribute("style", "width:100%;height:250px")
	el.SetData("tether-hook", "echarts")
	el.SetData("chart-option", html.EscapeString(option))
	return el
}

func buildChartOption(id, titleText, colour string, data []opts.LineData) string {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			ChartID: id,
			Width:   "100%",
			Height:  "250px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: titleText,
			Left:  "center",
			TitleStyle: &opts.TextStyle{
				FontSize: 13,
				Color:    "#888",
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{Show: opts.Bool(false)}),
		charts.WithYAxisOpts(opts.YAxis{
			SplitLine: &opts.SplitLine{
				LineStyle: &opts.LineStyle{Color: "#333"},
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(false)}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}),
	)

	xAxis := make([]string, len(data))
	for i := range xAxis {
		xAxis[i] = strconv.Itoa(i)
	}

	line.SetXAxis(xAxis).AddSeries("", data,
		charts.WithLineChartOpts(opts.LineChart{
			Smooth:     opts.Bool(true),
			ShowSymbol: opts.Bool(false),
		}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: colour, Width: 2}),
		charts.WithAreaStyleOpts(opts.AreaStyle{
			Color: colour + "30",
		}),
	)

	opt := line.RenderSnippet().Option
	return strings.TrimSuffix(strings.TrimSpace(opt), ";")
}

func toLineData(values []float64) []opts.LineData {
	data := make([]opts.LineData, len(values))
	for i, v := range values {
		data[i] = opts.LineData{Value: v}
	}
	return data
}

func intsToLineData(values []int) []opts.LineData {
	data := make([]opts.LineData, len(values))
	for i, v := range values {
		data[i] = opts.LineData{Value: v}
	}
	return data
}
