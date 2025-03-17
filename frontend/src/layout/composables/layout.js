import { computed, reactive, readonly } from 'vue';

const layoutConfig = reactive({
  preset: 'Aura',
  primary: 'blue',
  surface: null,
  darkTheme: window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches,
  menuMode: 'static'
});

// Listen for changes to the prefers-color-scheme media query
if (window.matchMedia) {
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    layoutConfig.darkTheme = e.matches;
    toggleDarkMode();
  });
}

const toggleDarkMode = () => {
  if (!document.startViewTransition) {
    executeDarkModeToggle();
    return;
  }

  document.startViewTransition(() => executeDarkModeToggle(event));
};

if (layoutConfig.darkTheme) {
  toggleDarkMode()
}

const executeDarkModeToggle = () => {
  layoutConfig.darkTheme = !layoutConfig.darkTheme;
  document.documentElement.classList.toggle('app-dark');
};

const layoutState = reactive({
  staticMenuDesktopInactive: false,
  overlayMenuActive: false,
  profileSidebarVisible: false,
  configSidebarVisible: false,
  staticMenuMobileActive: false,
  menuHoverActive: false,
  activeMenuItem: null
});

export function useLayout() {
  const setPrimary = (value) => {
    layoutConfig.primary = value;
  };

  const setSurface = (value) => {
    layoutConfig.surface = value;
  };

  const setPreset = (value) => {
    layoutConfig.preset = value;
  };

  const setActiveMenuItem = (item) => {
    layoutState.activeMenuItem = item.value || item;
  };

  const setMenuMode = (mode) => {
    layoutConfig.menuMode = mode;
  };

  const onMenuToggle = () => {
    if (layoutConfig.menuMode === 'overlay') {
      layoutState.overlayMenuActive = !layoutState.overlayMenuActive;
    }

    if (window.innerWidth > 991) {
      layoutState.staticMenuDesktopInactive = !layoutState.staticMenuDesktopInactive;
    } else {
      layoutState.staticMenuMobileActive = !layoutState.staticMenuMobileActive;
    }
  };

  const resetMenu = () => {
    layoutState.overlayMenuActive = false;
    layoutState.staticMenuMobileActive = false;
    layoutState.menuHoverActive = false;
  };

  const isSidebarActive = computed(() => layoutState.overlayMenuActive || layoutState.staticMenuMobileActive);

  const isDarkTheme = computed(() => layoutConfig.darkTheme);

  const getPrimary = computed(() => layoutConfig.primary);

  const getSurface = computed(() => layoutConfig.surface);

  return { layoutConfig: readonly(layoutConfig), layoutState: readonly(layoutState), onMenuToggle, isSidebarActive, isDarkTheme, getPrimary, getSurface, setActiveMenuItem, toggleDarkMode, setPrimary, setSurface, setPreset, resetMenu, setMenuMode };
}
