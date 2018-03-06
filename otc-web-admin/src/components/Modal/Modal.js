import { BOX_SHADOWS, SPACE, BORDER_RADIUS } from 'config';

import styled, { keyframes } from 'styled-components';
import { rem } from 'polished';
import Modal from 'react-modal';

const animateIn = keyframes`
  0% {
    opacity: 0;
    transform: translate(-50%, -50%) scale(0.9);
  }

  100% {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1);
  }
`;

export const styles = { overlay: { backgroundColor: 'rgba(0, 0, 0, 0.75)' } };

export default styled(Modal)`
  animation: 150ms ${animateIn} ease-in-out;
  padding: ${rem(SPACE[6])};
  background: white;
  outline: none;
  border-radius: ${BORDER_RADIUS.base};
  box-shadow: ${BOX_SHADOWS.base};
  max-width: calc(100vw - 5rem);
  width: ${rem(600)};
  text-align: center;
  transform: translate(-50%, -50%);
  position: absolute;
  top: 50%;
  left: 50%;
  right: auto;
  bottom: auto;
  margin-right: -50%;
`;
