import palx from 'palx';
import moment from 'moment';

export const DISTRIBUTION_START = moment('2017-08-08T12:00:00.000Z');
export const DISTRIBUTION_END = moment('2017-08-31T12:00:00.000Z');
export const COLORS = palx('#0072FF');
export const SPACE = [0, 4, 8, 12, 16, 20, 24, 32, 40, 48, 56, 64, 72];
export const FONT_SIZES = [11, 13, 14, 15, 17, 20, 24, 28, 36, 40];

export const BREAKPOINTS = {
  sm: 40,
  md: 52,
  lg: 64,
};

export const FONT_FAMILIES = {
  mono: '"Inconsolata", monospace, sans-serif',
  sans: '"Montreal", sans-serif',
};

export const BORDER_RADIUS = {
  base: '4px',
  pill: '1000px',
};

export const BOX_SHADOWS = {
  base: '0 1px 2px rgba(0, 0, 0, 0.25)',
  hover: '0 1px 4px rgba(0, 0, 0, 0.25)',
};

export const FLAGS = {
  chinese: true,
  russian: true,
  timeline: false,
  distribution: true,
  network: false,
};
