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
import Switch from "react-switch";

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

const Panel = styled(Box) `
  background-color: #fff;
  box-shadow: 1px 2px 4px rgba(0, 0, 0, .5);
  padding: ${rem(SPACE[6])} ${rem(SPACE[4])};
`;

const H3Styled = styled.h3`
  font-family: "Montreal", sans-serif;
  font-weight: 400;
  line-height: ${rem(1.75)};
`;

const sources = {
  internal: 'internal',
  exchange: 'exchange',
};

const TransparenContainer = styled(Container) `
  background-color: 'transparent' !important;
  padding-left: 0;
`;

const Wrapper = styled.div`
  background-color: ${COLORS.gray[1]};
  padding: ${rem(SPACE[5])} 0;

  ${media.md.css`
    padding: ${rem(SPACE[7])} 0;
  `}
`;

const TransparenWrapper = styled(Wrapper) `
  background-color: 'transparent';
`;

const PriceSource = ({ price, source }) => {
  return (<div>
    <Text as="p">
      {source === 'internal' ? 'Internal price' : 'Exchange price'} {price / 1e8} BTC
    </Text>
  </div>)
};

const invalidInputStyle = {
  borderColor: 'red'
};

const PriceSelector = ({
  selectedPrice,
  selectedSource,

  price,
  source,

  setPrice,
  setSource,
  save
}) => {
  const isPriceValid = !isNaN(selectedPrice) && selectedPrice !== '';
  return (
    <TransparenContainer>
      <Text as="a" mr={5}>Internal</Text>
      <Switch
        onChange={checked => setSource(checked ? sources.exchange : sources.internal)}
        checked={selectedSource === sources.exchange}
        offColor={COLORS.blue[7]}
        onColor={COLORS.green[8]}
        uncheckedIcon={false}
        checkedIcon={false}
      />
      <Text as="a" ml={5}>Exchange</Text>
      {selectedSource === sources.internal &&
        <Input
          value={selectedPrice}
          style={isPriceValid ? {} : invalidInputStyle}
          onChange={e => setPrice(e.target.value)}
          placeholder="Price" />}
      <TransparenContainer mt={5}>
        <Button
          bg={COLORS.green[8]}
          color="white"
          onClick={() => {
            if (isPriceValid || selectedSource === sources.exchange) {
              save(selectedSource, selectedPrice);
            }
          }}
        >Save</Button>
        <Button
          m={3}
          bg={COLORS.blue[7]}
          color="white"
          onClick={() => {
            setPrice(price);
            setSource(source);
          }}>Reset</Button>
      </TransparenContainer>
    </TransparenContainer>);
};

export default class extends React.Component {
  state = {
    price: 0,
    source: sources.internal,
    paused: true,
    loaded: false,

    selectedSource: sources.internal,
    selectedPrice: '0',
  };
  refreshStatus = async () => {
    const status = await getStatus();
    this.setState({
      ...this.state,
      ...status,

      selectedSource: status.source,
      selectedPrice: `${status.price / 1e8}`,
    });
  }
  componentWillMount = async () => {
    await this.refreshStatus();
    this.setState({ ...this.state, loaded: true });
  }
  setOctState = async pause => {
    this.setState({ ...this.state });
    await setOctState(pause);
    this.setState({ ...this.state, paused: pause });
  }
  setSource = source => {
    this.setState({ ...this.state, selectedSource: source });
  }
  setPrice = price => {
    this.setState({ ...this.state, selectedPrice: price });
  }
  save = async (source, price) => {
    await setPrice(Number.parseFloat(price) * 1e8, source);
    await this.refreshStatus();
  }
  render = () => {
    const {
      paused,
      source,
      price,

      loaded,

      selectedSource,
      selectedPrice, } = this.state;

    if (!loaded) return null;

    return (
      <div>
        <Helmet>
          <title>OTC Admin Panel</title>
        </Helmet>

        <Header external />

        <Wrapper>

          <Container>
            <Flex column width="33.3333%">
              <Panel>
                <H3Styled>OTC Status:</H3Styled>
                <Text>{paused ? 'Paused' : 'Running'}</Text>
                {paused
                  ? (<Button
                        bg={COLORS.green[8]}
                        color="white"
                        onClick={() => this.setOctState(false)}
                        >Start</Button>)
                  : (<Button 
                        bg={COLORS.red[7]}
                        color="white"
                        onClick={() => this.setOctState(true)}>Pause</Button>)}
              </Panel>
              <Panel mt={5}>
                <H3Styled>Price source:</H3Styled>
                <PriceSource source={source} price={price} />
                <PriceSelector
                  price={price}
                  source={source}

                  selectedSource={selectedSource}
                  selectedPrice={selectedPrice}

                  setSource={this.setSource}
                  setPrice={this.setPrice}
                  save={this.save} />
              </Panel>
            </Flex>
          </Container>
        </Wrapper>

        <Footer external />
      </div>
    );
  }
}
