(() => {
  const storedTheme = localStorage.getItem("synapseq-theme");
  const preferredTheme = window.matchMedia("(prefers-color-scheme: dark)")
    .matches
    ? "dark"
    : "light";

  document.documentElement.setAttribute(
    "data-theme",
    storedTheme || preferredTheme,
  );
})();

const initializePreviewPage = () => {
  const clampPercent = (value, minimum = 0) =>
    `${Math.max(Number(value) || 0, minimum)}%`;

  document.querySelectorAll("[data-left]").forEach((element) => {
    element.style.left = clampPercent(element.dataset.left);
  });

  document.querySelectorAll("[data-width]").forEach((element) => {
    element.style.width = clampPercent(element.dataset.width, 0.85);
  });

  document.querySelectorAll("[data-swatch-color]").forEach((element) => {
    element.style.background = element.dataset.swatchColor;
  });

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
  };

  themeButtons.forEach((button) => {
    button.addEventListener("click", () =>
      setTheme(button.dataset.themeButton),
    );
  });

  setTheme(root.getAttribute("data-theme") || "light");
};

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", initializePreviewPage, {
    once: true,
  });
} else {
  initializePreviewPage();
}
