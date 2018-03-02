import React from 'react';
import { Flex, Box } from 'grid-styled';

import Container from 'components/Container';
import Logo from 'components/Logo';

export default () => (
  <div>
    <Container>
      <Flex wrap my={[4, 8]} mx={-4} justify="center">
        <Box width={[1 / 2, 1 / 4]} my={2} px={4}>
          <Logo />
        </Box>
      </Flex>
    </Container>
  </div>
);
