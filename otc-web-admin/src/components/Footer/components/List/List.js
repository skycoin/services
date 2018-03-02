import React from 'react';
import PropTypes from 'prop-types';
import styled, { css } from 'styled-components';
import { FormattedMessage } from 'react-intl';

import Buy from 'components/Buy';
import Heading from 'components/Heading';
import Text from 'components/Text';
import Link from 'components/Link';

const LinkList = styled.ul`
  list-style: none;
  margin: 0;
`;

const styles = css`
  text-decoration: none;
  cursor: pointer;

  &:hover {
    text-decoration: underline;
  }
`;

const StyledLink = styled(Link)`
 ${styles}
`;

const StyledBuy = styled(Buy)`
 ${styles}
`;

const List = ({ heading, links }) => (
  <div>
    <Heading color="black" fontSize={3} heavy>
      <FormattedMessage id={heading} />
    </Heading>

    <Text as="div" fontSize={[1, 2, 3]} color="gray.8">
      <LinkList>
        {links.map(({ label, href, to, buy }, linkIndex) => (
          <li key={linkIndex}>
            {buy ? (
              <StyledBuy asAnchor>
                <FormattedMessage id={label} />
              </StyledBuy>
            ) : (
              <StyledLink href={href} to={to}>
                <FormattedMessage id={label} />
              </StyledLink>
            )}
          </li>
        ))}
      </LinkList>
    </Text>
  </div>
);

List.propTypes = {
  heading: PropTypes.string.isRequired,
  links: PropTypes.arrayOf(PropTypes.shape({
    buy: PropTypes.bool,
    to: PropTypes.string,
    label: PropTypes.string.isRequired,
  })).isRequired,
};

export default List;
