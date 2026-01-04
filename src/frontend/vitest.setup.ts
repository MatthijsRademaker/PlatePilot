import { config } from '@vue/test-utils';

// Stub Quasar components that render slots (so content is visible in tests)
const slotStub = {
  template: '<div><slot /></div>',
};

config.global.stubs = {
  // Card components - render slots for content visibility
  QCard: slotStub,
  QCardSection: slotStub,
  QCardActions: slotStub,

  // Display components - render slots
  QChip: slotStub,
  QBadge: slotStub,
  QBanner: slotStub,

  // Empty stubs for components without content
  QIcon: true,
  QImg: true,
  QAvatar: true,
  QSeparator: true,
  QSpace: true,
  QSkeleton: true,
  QSpinner: true,
  QSpinnerDots: true,

  // Form components
  QInput: true,
  QSelect: true,
  QBtn: slotStub,
  QCheckbox: true,
  QRadio: true,
  QToggle: true,
  QSlider: true,
  QForm: slotStub,

  // List components
  QList: slotStub,
  QItem: slotStub,
  QItemSection: slotStub,
  QItemLabel: slotStub,
  QExpansionItem: slotStub,

  // Layout components
  QLayout: slotStub,
  QPage: slotStub,
  QPageContainer: slotStub,
  QHeader: slotStub,
  QFooter: slotStub,
  QDrawer: slotStub,
  QToolbar: slotStub,
  QToolbarTitle: slotStub,

  // Navigation
  QTabs: slotStub,
  QTab: slotStub,
  QRouteTab: slotStub,
  QTabPanel: slotStub,
  QTabPanels: slotStub,

  // Scroll components
  QInfiniteScroll: slotStub,
  QPullToRefresh: slotStub,
  QScrollArea: slotStub,

  // Popup/Dialog
  QDialog: slotStub,
  QMenu: slotStub,
  QTooltip: true,
  QPopupProxy: true,

  // Router stubs
  RouterLink: {
    template: '<a><slot /></a>',
  },
  RouterView: {
    template: '<div><slot /></div>',
  },
};

// Mock ResizeObserver for Quasar components
global.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
};

// Mock IntersectionObserver for QInfiniteScroll etc.
global.IntersectionObserver = class IntersectionObserver {
  constructor() {}
  observe() {}
  unobserve() {}
  disconnect() {}
} as unknown as typeof IntersectionObserver;
