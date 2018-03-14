import React from 'react';
import styled from 'styled-components';
import PropTypes from 'prop-types';
import { Flex, Box } from 'grid-styled';
import { rem } from 'polished';

import { SPACE, COLORS } from 'config';
import Container from '../Container';
import Logo from '../Logo';
import Navigation from './components/Navigation';

const Wrapper = styled.div`
  padding: ${rem(SPACE[6])} 0;
  width: 100%;
  border-bottom: ${props => (props.border ? `2px solid ${COLORS.gray[1]}` : 'none')}
`;

const Header = ({ white, border }) => (
  <Wrapper border={border}>
    <Container>
      <Flex align="center" wrap>
        <Box width={[1 / 1, 1 / 4]}>
          <Logo white={white} />
        </Box>

        <Box width={[1 / 1, 3 / 4]}>
          <Navigation white={white} />
        </Box>
      </Flex>
    </Container>
  </Wrapper>
);

Header.propTypes = {
  white: PropTypes.bool,
  border: PropTypes.bool,
};

Header.defaultProps = {
  white: false,
  border: false,
};

export default Header;
