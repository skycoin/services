import styled from 'styled-components';
import { rem } from 'polished';

import {
  COLORS,
  SPACE,
  BORDER_RADIUS,
  BOX_SHADOWS,
  FONT_SIZES,
  FONT_FAMILIES,
} from 'config';

import media from '../../utils/media';

export default styled.input`
  display: inline-block;
  border-radius: ${BORDER_RADIUS.base};
  box-shadow: ${BOX_SHADOWS.base} inset;
  width: 100%;
  font-family: ${FONT_FAMILIES.mono};
  font-weight: 500;
  border: 1px solid ${COLORS.gray[4]};
  font-size: ${rem(FONT_SIZES[1])};
  box-sizing: border-box;
  padding: ${rem(SPACE[3])};
  color: ${COLORS.black};
  margin-bottom: ${rem(SPACE[4])};

  &:focus {
    outline: none;
    border-color: ${COLORS.base};
  }

  padding: ${rem(SPACE[3])};

  ${media.sm.css`
    font-size: ${rem(FONT_SIZES[3])};
    padding: ${rem(SPACE[4])};
  `}
`;
