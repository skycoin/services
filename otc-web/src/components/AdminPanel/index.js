import React from 'react';
import PropTypes from 'prop-types';
import Helmet from 'react-helmet';
import { FormattedMessage, FormattedHTMLMessage, injectIntl } from 'react-intl';
import { COLORS, SPACE, BOX_SHADOWS, BORDER_RADIUS } from 'config';
import styled from 'styled-components';
import { rem } from 'polished';

import Header from 'components/Header';
import Footer from 'components/Footer';
import media from '../../utils/media';

const Wrapper = styled.div`
  background-color: ${COLORS.gray[1]};
  padding: ${rem(SPACE[5])} 0;

  ${media.md.css`
    padding: ${rem(SPACE[7])} 0;
  `}
`;

export default injectIntl(class extends React.Component {
  static propTypes = {
    intl: PropTypes.shape({
      formatMessage: PropTypes.func.isRequired,
    }).isRequired,
  };
  render = () => {
    const { intl } = this.props;
    return (
      <div>
        <Helmet>
          <title>{intl.formatMessage({ id: 'distribution.title' })}</title>
        </Helmet>

        <Header external />

        <Wrapper>
          Admin Panel
        </Wrapper>

        <Footer external suffix="/admin" />
      </div>
    );
  }
});
