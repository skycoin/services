import styled, { css } from 'styled-components';
import { space, width, fontSize, color } from 'styled-system';
import { rem, darken } from 'polished';
import get from 'lodash/get';

import Link from 'components/Link';
import media from 'utils/media';
import { COLORS, FONT_FAMILIES, FONT_SIZES, SPACE, BORDER_RADIUS, BOX_SHADOWS } from 'config';

export default styled(Link)`
  display: inline-block;
  border: 1px solid transparent;
  cursor: pointer;
  line-height: 1;
  font-weight: 700;
  box-sizing: border-box;
  text-decoration: none;
  text-align: center;

  border-radius: ${props => (props.pill ? BORDER_RADIUS.pill : BORDER_RADIUS.base)};
  box-shadow: ${BOX_SHADOWS.base};
  font-family: ${FONT_FAMILIES.mono};
  font-size: ${rem(FONT_SIZES[3])};
  padding: ${rem(SPACE[3])} ${rem(SPACE[6])};
  transition: 150ms ease-in-out;

  ${fontSize}
  ${color}
  ${space}
  ${width}

  &:hover {
    background-color: ${props =>
      props.bg && darken(0.05, get(COLORS, props.bg) || props.bg)
    };

    box-shadow: ${BOX_SHADOWS.hover};
  }

  &:active {
    background-color: ${props =>
      props.bg && darken(0.1, get(COLORS, props.bg) || props.bg)
    };
  }

  ${props => props.outlined && css`
    border: 1px solid currentColor;

    &, &:hover, &:active {
      box-shadow: none;
    }

    &:hover {
      background-color: rgba(0, 0, 0, 0.05);
    }

    &:active {
      background-color: rgba(0, 0, 0, 0.1);
    }
  `}

  ${props => props.big && media.sm.css`
    padding: ${rem(SPACE[4])} ${rem(SPACE[7])};
  `}
`;
