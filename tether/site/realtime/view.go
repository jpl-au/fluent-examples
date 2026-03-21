package realtime

import (
	"fmt"
	"html"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/tether/components/composite/monitor"
	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the real-time dashboard page with live go-echarts
// charts for CPU, heap, and goroutine metrics.
func Render(s State) node.Node {
	return cpage.New(
		panel.Card(
			"System Monitor",
			"Live Go runtime metrics pushed from the server every second. A Session.Go "+
				"goroutine reads runtime.MemStats and measures process CPU via "+
				"syscall.Getrusage, updates the session state, and the framework "+
				"re-renders and diffs the charts automatically. The charts are built "+
				"by go-echarts - a Go charting library that wraps Apache ECharts. "+
				"Open the page and watch the charts move.",
			"sess.Go · sess.Update · go-echarts", panel.WS|panel.SSE,
			monitor.Charts(
				chartDiv("chartcpu", "CPU (%)", "#ee6666", toLineData(s.CPUPercent)),
				chartDiv("chartheap", "Heap (MB)", "#5470c6", toLineData(s.HeapMB)),
				chartDiv("chartgoroutines", "Goroutines", "#91cc75", intsToLineData(s.Goroutines)),
			).Dynamic("monitor-charts"),
		),
	)
}

// chartDiv builds a div wired to the echarts JS hook. The chart
// option JSON is built by go-echarts and stored HTML-escaped in a
// data attribute. The hook reads it with getAttribute (which auto-
// unescapes entities) and calls echarts.setOption().
func chartDiv(id, titleText, colour string, data []opts.LineData) node.Node {
	option := buildChartOption(id, titleText, colour, data)
	el := monitor.Chart(id)
	el.SetAttribute("style", "width:100%;height:250px")
	el.SetData("tether-hook", "echarts")
	el.SetData("chart-option", html.EscapeString(option))
	return el
}

// buildChartOption uses go-echarts to produce the ECharts JSON config.
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
		xAxis[i] = fmt.Sprintf("%d", i)
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

	// RenderSnippet().Option may include a trailing semicolon from the
	// go-echarts base template - strip it so the value is valid JSON.
	opt := line.RenderSnippet().Option
	return strings.TrimSuffix(strings.TrimSpace(opt), ";")
}

// toLineData converts float64 metric samples into go-echarts line
// data points for chart rendering.
func toLineData(values []float64) []opts.LineData {
	data := make([]opts.LineData, len(values))
	for i, v := range values {
		data[i] = opts.LineData{Value: fmt.Sprintf("%.1f", v)}
	}
	return data
}

// intsToLineData is the int-typed variant of toLineData.
func intsToLineData(values []int) []opts.LineData {
	data := make([]opts.LineData, len(values))
	for i, v := range values {
		data[i] = opts.LineData{Value: v}
	}
	return data
}
