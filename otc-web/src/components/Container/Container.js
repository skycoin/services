import styled from 'styled-components';
import { rem } from 'polished';

import { SPACE } from 'config';

export default styled.div`
  margin: 0 auto;
  max-width: 60rem;
  padding: 0 ${rem(SPACE[4])};
`;
