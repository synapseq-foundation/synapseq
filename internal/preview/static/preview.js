(() => {
  const params = new URLSearchParams(window.location.search);
  const requestedTheme = params.get("theme");
  const urlTheme =
    requestedTheme === "light" || requestedTheme === "dark"
      ? requestedTheme
      : "";
  const storedTheme = localStorage.getItem("synapseq-theme");
  const preferredTheme = window.matchMedia("(prefers-color-scheme: dark)")
    .matches
    ? "dark"
    : "light";

  document.documentElement.setAttribute(
    "data-theme",
    urlTheme || storedTheme || preferredTheme,
  );
})();

const initializePreviewPage = () => {
  const params = new URLSearchParams(window.location.search);
  const requestedTheme = params.get("theme");
  const hasThemeParam = requestedTheme === "light" || requestedTheme === "dark";
  const requestedHome = params.get("home");
  const clampPercent = (value, minimum = 0) =>
    `${Math.max(Number(value) || 0, minimum)}%`;
  const charts = [];
  const pointStyles = ["circle", "rectRot", "triangle", "rectRounded"];

  const cssColor = (name) =>
    getComputedStyle(document.documentElement).getPropertyValue(name).trim();

  const readChartSeries = (element) => {
    try {
      return JSON.parse(element.dataset.chartSeries || "{}");
    } catch {
      return {};
    }
  };

  const chartTheme = () => ({
    text: cssColor("--text"),
    muted: cssColor("--muted"),
    line: cssColor("--line"),
    lineStrong: cssColor("--line-strong"),
    panel: cssColor("--panel-strong"),
  });

  const applyChartTheme = (chart) => {
    const theme = chartTheme();
    chart.options.color = theme.muted;
    chart.options.scales.x.grid.color = theme.line;
    chart.options.scales.x.border.color = theme.lineStrong;
    chart.options.scales.y.grid.color = theme.line;
    chart.options.scales.y.border.color = theme.lineStrong;
    chart.options.scales.y.ticks.color = theme.muted;
    chart.options.plugins.tooltip.backgroundColor = theme.panel;
    chart.options.plugins.tooltip.titleColor = theme.text;
    chart.options.plugins.tooltip.bodyColor = theme.muted;
    chart.options.plugins.tooltip.borderColor = theme.lineStrong;
    chart.update("none");
  };

  const resizeVisibleCharts = () => {
    requestAnimationFrame(() => {
      charts.forEach((chart) => {
        if (chart.canvas.closest(".graph-metric-view.is-active")) {
          chart.resize();
          chart.update("none");
        }
      });
    });
  };

  const initializeCharts = () => {
    if (!window.Chart) {
      return;
    }

    window.Chart.defaults.font.family = '"Avenir Next", "Segoe UI", sans-serif';

    document.querySelectorAll("[data-chart-series]").forEach((canvas) => {
      const series = readChartSeries(canvas);
      if (!series.curve || series.curve.length === 0) {
        return;
      }

      const chart = new window.Chart(canvas, {
        type: "line",
        data: {
          datasets: [
            {
              label: series.label,
              data: series.curve,
              parsing: false,
              borderColor: series.color,
              backgroundColor: series.color,
              borderWidth: 2.4,
              cubicInterpolationMode: "monotone",
              tension: 0.28,
              pointRadius: 0,
              pointHitRadius: 10,
            },
            {
              label: `${series.label} nodes`,
              data: series.markers || [],
              parsing: false,
              showLine: false,
              borderColor: series.color,
              backgroundColor: series.color,
              pointBorderColor: "#fffaf1",
              pointBorderWidth: 2,
              pointRadius: 4.5,
              pointHoverRadius: 7,
              pointStyle: (context) =>
                pointStyles[context.dataIndex % pointStyles.length],
            },
          ],
        },
        options: {
          animation: false,
          responsive: true,
          maintainAspectRatio: false,
          interaction: {
            intersect: false,
            mode: "nearest",
          },
          plugins: {
            legend: {
              display: false,
            },
            tooltip: {
              displayColors: false,
              borderWidth: 1,
              padding: 10,
              callbacks: {
                title: (items) => items[0]?.raw?.timeLabel || "",
                label: (item) => item.raw?.pointLabel || item.raw?.valueLabel || "",
                afterLabel: (item) =>
                  item.raw?.valueLabel ? `Value: ${item.raw.valueLabel}` : "",
              },
            },
          },
          scales: {
            x: {
              type: "linear",
              display: false,
              min: 0,
              max: series.durationMs || undefined,
              ticks: {
                maxTicksLimit: 6,
              },
              grid: {
                display: false,
                drawTicks: false,
              },
            },
            y: {
              ticks: {
                maxTicksLimit: 4,
              },
              grid: {
                drawTicks: false,
              },
            },
          },
        },
      });

      applyChartTheme(chart);
      charts.push(chart);
    });
  };

  document.querySelectorAll("[data-left]").forEach((element) => {
    element.style.left = clampPercent(element.dataset.left);
  });

  document.querySelectorAll("[data-width]").forEach((element) => {
    element.style.width = clampPercent(element.dataset.width, 0.85);
  });

  if (hasThemeParam) {
    document.querySelectorAll("[data-theme-switcher]").forEach((element) => {
      element.hidden = true;
    });
  }

  if (requestedHome === "false") {
    document.querySelectorAll("[data-home-link]").forEach((element) => {
      element.hidden = true;
    });
  }

  const graphTabs = document.querySelectorAll("[data-graph-target]");
  const graphViews = document.querySelectorAll("[data-graph-view]");

  const setGraphMetric = (key) => {
    graphTabs.forEach((tab) => {
      const isActive = tab.dataset.graphTarget === key;
      tab.classList.toggle("is-active", isActive);
      tab.setAttribute("aria-pressed", isActive ? "true" : "false");
    });

    graphViews.forEach((view) => {
      view.classList.toggle("is-active", view.dataset.graphView === key);
    });
    resizeVisibleCharts();
  };

  graphTabs.forEach((tab) => {
    tab.addEventListener("click", () =>
      setGraphMetric(tab.dataset.graphTarget),
    );
  });

  const defaultGraphTab = document.querySelector(
    "[data-graph-target].is-active",
  );
  if (defaultGraphTab) {
    setGraphMetric(defaultGraphTab.dataset.graphTarget);
  }

  const root = document.documentElement;
  const themeButtons = document.querySelectorAll("[data-theme-button]");

  const setTheme = (theme) => {
    root.setAttribute("data-theme", theme);
    localStorage.setItem("synapseq-theme", theme);
    themeButtons.forEach((button) => {
      const isActive = button.dataset.themeButton === theme;
      button.classList.toggle("is-active", isActive);
      button.setAttribute("aria-pressed", isActive ? "true" : "false");
    });
    charts.forEach(applyChartTheme);
  };

  themeButtons.forEach((button) => {
    button.addEventListener("click", () =>
      setTheme(button.dataset.themeButton),
    );
  });

  setTheme(root.getAttribute("data-theme") || "light");
  initializeCharts();
  charts.forEach(applyChartTheme);
  resizeVisibleCharts();

  if (window.lucide) {
    window.lucide.createIcons();
  }
};

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", initializePreviewPage, {
    once: true,
  });
} else {
  initializePreviewPage();
}
