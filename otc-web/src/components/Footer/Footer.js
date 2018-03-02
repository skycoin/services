import React from 'react';
import { Flex, Box } from 'grid-styled';

import Container from 'components/Container';
import Text from 'components/Text';
import Logo from 'components/Logo';

import Languages from './components/Languages';
import List from './components/List';
import Email from './components/Email';
import content from './content';

export default () => (
  <div>
    <Container>
      <Flex wrap my={[4, 8]} mx={-4}>
        <Box width={[1 / 2, 1 / 4]} my={2} px={4}>
          <Logo />

          <Text fontSize={[1, 2, 3]} color="gray.8" heavy mt={2}>
            <Email />
          </Text>

          <Text as="div" fontSize={[0, 0, 1]} color="gray.8" heavy>
            <Languages />
          </Text>
        </Box>

        {content.map(({ heading, links }, sectionIndex) => (
          <Box width={[1 / 2, 1 / 4]} my={2} px={4} key={sectionIndex}>
            <List heading={heading} links={links} />
          </Box>
        ))}
      </Flex>
    </Container>
  </div>
);
