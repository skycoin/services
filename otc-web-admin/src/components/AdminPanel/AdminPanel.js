/* eslint-disable no-alert */

import React from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';
import moment from 'moment';
import Helmet from 'react-helmet';
import { Flex, Box } from 'grid-styled';
import { FormattedMessage, FormattedHTMLMessage, injectIntl } from 'react-intl';
import { rem } from 'polished';
import { COLORS, SPACE, BOX_SHADOWS, BORDER_RADIUS } from 'config';
import QRCode from 'qrcode.react';

import Button from 'components/Button';
import Container from 'components/Container';
import Footer from 'components/Footer';
import Header from 'components/Header';
import Heading from 'components/Heading';
import Input from 'components/Input';
import Modal, { styles } from 'components/Modal';
import Text from 'components/Text';
import media from '../../utils/media';

import { getStatus, setPrice, setOctState } from './admin-api';

const Wrapper = styled.div`
  background-color: ${COLORS.gray[1]};
  padding: ${rem(SPACE[5])} 0;

  ${media.md.css`
    padding: ${rem(SPACE[7])} 0;
  `}
`;

const PriceSource = ({ price, source }) => {
  return (<div>
    <strong>{source === 'internal' ? 'Internal price' : 'Exchange price'}: </strong>
    {price / 1e8} BTC
  </div>)
};

const PriceSelector = ({
  selectedPrice,
  selectedSource,

  price,
  source,

  setPrice,
  setSource,
}) => {
  return (<div />);
};

export default class extends React.Component {
  state = {
    price: 0,
    source: 'internal',
    paused: true,
    loaded: false,

    selectedSource: 'internal',
    selectedPrice: 0,
  };
  componentWillMount = async () => {
    const status = await getStatus();
    this.setState({
      ...this.state,
      loaded: true,
      ...status,

      selectedSource: status.source,
      selectedPrice: status.price,
    });
  }
  setOctState = async pause => {
    this.setState({ ...this.state, loaded: false });
    await setOctState(pause);
    this.setState({ ...this.state, loaded: true, paused: pause });
  }
  setSource = source => {
    this.setState({ ...this.state, selectedSource: source });
  }
  setPrice = price => {
    this.setState({ ...this.state, selectedPrice: price });
  }
  render = () => {
    const {
      paused,
      source,
      price,

      selectedSource,
      selectedPrice, } = this.state;
    return (
      <div>
        <Helmet>
          <title>OTC Admin Panel</title>
        </Helmet>

        <Header external />

        <Wrapper>

          <Container>
            <Flex flexDirection="row">
              <Box px={10}>
                <Text as="h3">OTC Status:</Text>
                <Text>{paused ? 'Paused' : 'Running'}</Text>
                {paused
                  ? (<Button onClick={() => this.setOctState(false)}>Start</Button>)
                  : (<Button onClick={() => this.setOctState(true)}>Pause</Button>)}
              </Box>
              <Box px={10}>
                <h3>Price source:</h3>
                <PriceSource source={source} price={price} />
                <PriceSelector
                  price={price}
                  source={source}

                  selectedSource={selectedSource}
                  selectedPrice={selectedPrice}

                  setSource={this.setSource}
                  setPrice={this.setPrice} />
              </Box>
            </Flex>
          </Container>
        </Wrapper>

        <Footer external />
      </div>
    );
  }
}
