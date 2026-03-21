// JS lifecycle hooks registered with the tether runtime.
//
// Each hook receives the DOM element and can set up or tear down
// client-side behaviour that cannot be expressed server-side.

Tether.hooks = Tether.hooks || {};

// tooltip: adds a title attribute from data-tooltip on mount,
// demonstrating the mounted/updated/destroyed lifecycle.
Tether.hooks.tooltip = {
  mounted: function (el) {
    var text = el.getAttribute("data-tooltip");
    if (text) {
      el.setAttribute("title", text);
      el.style.cursor = "help";
    }
  },
  updated: function (el) {
    var text = el.getAttribute("data-tooltip");
    if (text) {
      el.setAttribute("title", text);
    }
  },
  destroyed: function () {
    // Nothing to clean up for a tooltip.
  },
};

// echarts: renders and updates an Apache ECharts chart from a JSON
// option stored in data-chart-option. The option is built server-side
// by go-echarts and HTML-escaped; getAttribute auto-unescapes it.
//
// On mount: initialise the chart and call setOption.
// On update: re-read the option and update the chart in place.
// On destroy: dispose the chart and clean up the resize listener.
Tether.hooks.echarts = {
  mounted: function (el) {
    var optionStr = el.getAttribute("data-chart-option");
    if (!optionStr) return;

    var chart = echarts.init(el);
    try {
      var option = JSON.parse(optionStr);
      option.animation = false;
      chart.setOption(option);
    } catch (e) {
      console.error("echarts hook: parse error on mount", e);
    }

    el._echartsInstance = chart;

    var onResize = function () {
      chart.resize();
    };
    window.addEventListener("resize", onResize);
    el._echartsResize = onResize;
  },

  updated: function (el) {
    var optionStr = el.getAttribute("data-chart-option");
    if (!optionStr) return;

    // Idiomorph removes ECharts' internal DOM nodes (canvas, etc.)
    // during the morph, so the existing instance is broken. Dispose
    // it safely and re-initialise from scratch. With animation
    // disabled the chart appears instantly - no visible flicker.
    if (el._echartsInstance) {
      try {
        el._echartsInstance.dispose();
      } catch (_) {
        /* DOM already gone - safe to ignore */
      }
      el._echartsInstance = null;
    }
    if (el._echartsResize) {
      window.removeEventListener("resize", el._echartsResize);
      el._echartsResize = null;
    }

    Tether.hooks.echarts.mounted(el);
  },

  destroyed: function (el) {
    if (el._echartsInstance) {
      el._echartsInstance.dispose();
      el._echartsInstance = null;
    }
    if (el._echartsResize) {
      window.removeEventListener("resize", el._echartsResize);
    }
  },
};
