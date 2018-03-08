import React from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';
import { rem } from 'polished';
import { Flex, Box } from 'grid-styled';

import Link from 'components/Link';
import Text from 'components/Text';
import logo from './logo.svg';
import logoWhite from './logoWhite.svg';

const StyledLink = styled(Link) `
  display: block;
`;

const Img = styled.img.attrs({
  alt: 'Skycoin',
}) `
  display: block;
  height: ${rem(40)};
  max-width: 100%;
`;

const Logo = props => (
  <Flex>
    <Box>
      <StyledLink to="/">
        <Img {...props} src={props.white ? logoWhite : logo} />
      </StyledLink>
    </Box>
    <Box>
      <Text as="h3" style={{ marginTop: '10px', marginLeft: '5px' }}> OTC Admin</Text>
    </Box>
  </Flex>
);

Logo.propTypes = {
  white: PropTypes.bool,
};

Logo.defaultProps = {
  white: false,
};

export default Logo;
