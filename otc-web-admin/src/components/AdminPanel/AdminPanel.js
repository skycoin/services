/* eslint-disable no-alert */

import React from 'react';
import styled from 'styled-components';
import Helmet from 'react-helmet';
import { Flex, Box } from 'grid-styled';
import { rem } from 'polished';
import { COLORS, SPACE } from 'config';
import TimeAgo from 'react-timeago';
import 'react-widgets/dist/css/react-widgets.css';

import Button from 'components/Button';
import Container from 'components/Container';
import Footer from 'components/Footer';
import Header from 'components/Header';
import Input from 'components/Input';
import Text from 'components/Text';
import media from '../../utils/media';

import { getStatus, setPrice, setSource, setOctState } from './admin-api';

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

const UpdatedPriceContainer = styled(Text) `
  font-size: 10px;
  line-height: 1.2;
  display: block;
  color: ${COLORS.gray[5]}
`;

const PriceHeading = styled(Text) `
  line-height: 1.2;
  margin-bottom: ${rem(SPACE[1])};
`;


const DateTimeView = ({ dt }) =>
  (<div>{new Date(dt * 1000).toLocaleDateString()} {new Date(dt * 1000).toLocaleTimeString()}</div>);

const UpdatePriceLabel = ({ updated, as }) => (
  <UpdatedPriceContainer as={as || 'p'}>
    Updated <TimeAgo date={new Date(updated * 1000)} /> at <DateTimeView dt={updated} />
  </UpdatedPriceContainer>
);

const PriceSource = ({ prices, source }) => {
  return (<div>
    <Text as="div">
      <PriceHeading as="p">
        {source === sources.internal ? `Internal price ${prices.internal / 1e8} ` : `Exchange price ${prices.exchange / 1e8} `} BTC
      </PriceHeading>
      <UpdatePriceLabel updated={source === sources.internal ? prices.internal_updated : prices.exchange_updated} />
    </Text>
  </div>)
};

const invalidInputStyle = {
  borderColor: 'red'
};

const Radio = styled.input`
  opacity: 0;
  position: absolute;
`;

const RadioLabel = styled.label.attrs({ htmlFor: props => props.for }) `
  display: block;
  margin: ${rem(SPACE[2])} 0;
`;

const PriceType = styled.div`
  ${RadioLabel} {
    position: relative;
    padding-left: 30px;

    &::before {
      content: '';
      width: 10px;
      height: 10px;
      display: inline-block;
      border-radius: 100%;
      background: white;
      border: 5px solid white;
      box-shadow: 0 0 1px black;
      position: absolute;
      top: 0;
      left: 0;
    }
  }

  ${Radio}:checked + ${RadioLabel}::before {
    background: black;
  }
`;

const PriceSelector = ({
  selectedPrice,
  selectedSource,

  prices,
  source,

  setPrice,
  setSource,
  save
}) => {
  const isPriceValid = !isNaN(selectedPrice) && selectedPrice !== '';
  return (
    <TransparenContainer>
      <PriceType>
        <Radio
          type="radio"
          id="radio_exchange"
          value={sources.exchange}
          checked={selectedSource === sources.exchange}
          onChange={() => setSource(sources.exchange)} />
        <RadioLabel for="radio_exchange">
          Exchange ({prices.exchange / 1e8} BTC)
          <UpdatePriceLabel as="span" updated={prices.exchange_updated} />
        </RadioLabel>
      </PriceType>
      <PriceType>
        <Radio
          type="radio"
          id="radio_internal"
          value={sources.internal}
          checked={selectedSource === sources.internal}
          onChange={() => setSource(sources.internal)} />
        <RadioLabel for="radio_internal">
          Internal ({prices.internal / 1e8} BTC)
          <UpdatePriceLabel as="span" updated={prices.internal_updated} />
        </RadioLabel>
        <Input
          value={selectedPrice}
          style={isPriceValid ? {} : invalidInputStyle}
          onChange={e => setPrice(e.target.value)}
          placeholder="Price" />
      </PriceType>
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
      </TransparenContainer>
    </TransparenContainer>);
};

const OtcUnavailableMessage = () => (
  <Wrapper>
    <Container>
      <Text>Skycoin OTC is currently unavailable.</Text>
    </Container>
  </Wrapper>);

export default class extends React.Component {
  state = {
    otcAvailable: false,
    prices: {
      internal: 0,
      exchange: 0,
      exchange_updated: 1519131184,
      internal_updated: 1519131184,
    },
    source: sources.internal,
    paused: true,
    loaded: false,

    selectedSource: sources.internal,
    selectedPrice: '0',
  };
  refreshStatus = async () => {
    try {
      const status = await getStatus();
      this.setState({
        ...this.state,
        ...status,

        otcAvailable: true,

        selectedSource: status.source,
        selectedPrice: `${status.prices.internal / 1e8}`,
      });
    } catch (e) {
      console.error(e);
    }
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
    const p = Number.parseFloat(price) * 1e8;
    await setPrice(Math.round(p));
    await setSource(source);
    await this.refreshStatus();
  }
  render = () => {
    const {
      paused,
      source,
      prices,

      loaded,
      otcAvailable,

      selectedSource,
      selectedPrice, } = this.state;

    if (!loaded) return null;
    return (
      <div>
        <Helmet>
          <title>OTC Admin Panel</title>
        </Helmet>

        <Header external />
        {!otcAvailable && <OtcUnavailableMessage />}
        {otcAvailable &&
          <Wrapper>

            <Container>
              <Flex row wrap justify="center" align="flex-start">
                <Panel flex={1}>
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
                <Panel ml={[0, 5]} mt={[5, 0]} flex="3 0 auto">
                  <H3Styled>Price source:</H3Styled>
                  <PriceSource source={source} prices={prices} />
                  <PriceSelector
                    prices={prices}
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
        }

        <Footer external />
      </div>
    );
  }
}
