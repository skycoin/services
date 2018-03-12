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

import { checkStatus, getAddress, getConfig, } from '../../utils/distributionAPI';

const Wrapper = styled.div`
  background-color: ${COLORS.gray[1]};
  padding: ${rem(SPACE[5])} 0;

  ${media.md.css`
    padding: ${rem(SPACE[7])} 0;
  `}
`;

const Address = Heading.extend`
  display: flex;
  justify-content: space-between;
  word-break: break-all;
  background-color: ${COLORS.gray[0]};
  border-radius: ${BORDER_RADIUS.base};
  box-shadow: ${BOX_SHADOWS.base};
  padding: 1rem;
  margin-bottom: 1rem;
`;

const StatusModal = ({ status, statusIsOpen, closeModals, skyAddress, intl }) => (
  <Modal
    contentLabel="Status"
    style={styles}
    isOpen={statusIsOpen}
    onRequestClose={closeModals}
  >
    <Heading heavy color="black" fontSize={[2, 3]} my={[3, 5]}>
      <FormattedMessage
        id="distribution.statusFor"
        values={{
          skyAddress,
        }}
      />
    </Heading>

    <Text as="div" color="black" fontSize={[2, 3]} my={[3, 5]}>
      {status.map((status, i) => (
        <p key={i}>
          <FormattedMessage
            id={`distribution.statuses.${status.status}`}
            values={{
              updated: moment.unix(status.updated_at).locale(intl.locale).format('LL LTS'),
            }}
          />
        </p>
      ))}
    </Text>
  </Modal>
);

const getDisabledReasonMessage = reason => {
  switch (reason) {
    case 'coinsSoldOut':
      return <FormattedMessage id="distribution.errors.coinsSoldOut" />;
    case 'paused':
      return <FormattedMessage id="distribution.errors.paused" />;
    default:
      return <FormattedMessage id="distribution.headingEnded" />;
  }
};

const StatusErrorMessage = ({ disabledReason }) => (<Flex column>
  <Heading heavy as="h2" fontSize={[5, 6]} color="black" mb={[4, 6]}>
    {getDisabledReasonMessage(disabledReason)}
  </Heading>
  <Text heavy color="black" fontSize={[2, 3]} as="div">
    <FormattedHTMLMessage id="distribution.ended" />
  </Text>
</Flex>);

const btcToSatochi = 0.00000001
const roundTo = 100000000
const DistributionFormInfo = ({ sky_btc_exchange_rate, balance }) => (
  <div>
    <Heading heavy as="h2" fontSize={[5, 6]} color="black" mb={[4, 6]}>
      <FormattedMessage id="distribution.heading" />
    </Heading>
    {sky_btc_exchange_rate &&
      <Text heavy color="black" fontSize={[2, 3]} mb={[4, 6]} as="div">
        <FormattedHTMLMessage
          id="distribution.rate"
          values={{
            rate: +(Math.round(sky_btc_exchange_rate * btcToSatochi * roundTo) / roundTo),
          }}
        />
      </Text>}
    <Text heavy color="black" fontSize={[2, 3]} mb={[4, 6]} as="div">
      <FormattedMessage
        id="distribution.inventory"
        values={{
          coins: balance.toString(),
        }}
      />
    </Text>

    <Text heavy color="black" fontSize={[2, 3]} as="div">
      <FormattedHTMLMessage id="distribution.instructions" />
    </Text>
  </div>);

const DistributionForm = ({
  sky_btc_exchange_rate,
  balance,
  intl,

  address,
  handleAddressChange,

  drop_address,
  status_address,
  handleStatusAddressChange,
  getAddress,
  addressLoading,

  checkStatus,
  statusLoading,
}) => (
    <Flex justify="center">
      <Box width={[1 / 1, 1 / 1, 2 / 3]} py={[5, 7]}>
        <DistributionFormInfo sky_btc_exchange_rate={sky_btc_exchange_rate} balance={balance} />

        <Input
          placeholder={intl.formatMessage({ id: 'distribution.enterAddress' })}
          value={address}
          onChange={handleAddressChange}
        />

        {drop_address && <Address heavy color="black" fontSize={[2, 3]} as="div">
          <Box>
            <strong><FormattedHTMLMessage id="distribution.btcAddress" />: </strong>
            {drop_address}
          </Box>
          <Box px={5}>
            <QRCode value={drop_address} size={64} />
          </Box>
        </Address>}

        <div>
          <Button
            big
            onClick={getAddress}
            color="white"
            bg="base"
            mr={[2, 5]}
            fontSize={[1, 3]}
          >
            {addressLoading
              ? <FormattedMessage id="distribution.loading" />
              : <FormattedMessage id="distribution.getAddress" />}
          </Button>
        </div>
        <Input
          placeholder={intl.formatMessage({ id: 'distribution.enterAddressBTC' })}
          value={status_address}
          onChange={handleStatusAddressChange}
        />
        <div>
          <Button
            onClick={checkStatus}
            color="base"
            big
            outlined
            fontSize={[1, 3]}
          >
            {statusLoading
              ? <FormattedMessage id="distribution.loading" />
              : <FormattedMessage id="distribution.checkStatus" />}
          </Button>
        </div>
      </Box>
    </Flex>);

class Distribution extends React.Component {
  state = {
    status: [],
    skyAddress: null,
    drop_address: '',
    status_address: '',
    statusIsOpen: false,
    addressLoading: false,
    statusLoading: false,
    enabled: true,
    balance: 0,
    sky_btc_exchange_rate: null,
  };
  componentWillMount = async () => {
    try {
      const config = await getConfig();
      const stateMutation = { sky_btc_exchange_rate: config.price };
      switch (config.otcStatus) {
        case 'SOLD_OUT':
          stateMutation.disabledReason = 'coinsSoldOut';
          stateMutation.enabled = false;
          break;
        case 'PAUSED':
          stateMutation.disabledReason = 'paused';
          stateMutation.enabled = false;
          break;
        case 'WORKING':
          stateMutation.balance = config.balance;
          stateMutation.enabled = true;
          break;
      }
      this.setState({ ...this.state, ...stateMutation });
    } catch (_) {
      this.setState({ ...this.state, enabled: false, disabledReason: 'closed' });
    }
  }

  getAddress = () => {
    if (!this.state.skyAddress) {
      return alert(
        this.props.intl.formatMessage({
          id: 'distribution.errors.noSkyAddress',
        }),
      );
    }

    this.setState({
      addressLoading: true,
    });

    return getAddress(this.state.skyAddress)
      .then((res) => {
        this.setState({
          drop_address: res.drop_address,
          status_address: res.drop_address,
        });
      })
      .catch((err) => {
        alert(err.message);
      })
      .then(() => {
        this.setState({
          addressLoading: false,
        });
      });
  }

  handleAddressChange = (event) => {
    this.setState({
      skyAddress: event.target.value.trim(),
    });
  }

  handleStatusAddressChange = (event) => {
    this.setState({
      status_address: event.target.value.trim(),
    });
  }

  closeModals = () => {
    this.setState({
      statusIsOpen: false,
    });
  }

  checkStatus = () => {
    if (!this.state.status_address) {
      return alert(
        this.props.intl.formatMessage({
          id: 'distribution.errors.noDropAddress',
        }),
      );
    }

    this.setState({
      statusLoading: true,
    });

    return checkStatus({ drop_address: this.state.status_address, drop_currency: 'BTC' })
      .then((res) => {
        this.setState({
          statusIsOpen: true,
          status: res,
          statusLoading: false,
        });
      })
      .catch((err) => {
        alert(err.message);
      });
  }

  render = () => {
    const { intl } = this.props;
    const {
      statusIsOpen,
      skyAddress,
      status,
      disabledReason,
      enabled,
      sky_btc_exchange_rate,
      balance,
      drop_address,
      status_address,
      addressLoading,
      statusLoading } = this.state;
    return (
      <div>
        <Helmet>
          <title>{intl.formatMessage({ id: 'distribution.title' })}</title>
        </Helmet>

        <Header external />

        <Wrapper>
          <StatusModal
            statusIsOpen={statusIsOpen}
            closeModals={this.closeModals}
            skyAddress={status_address}
            intl={intl}
            status={status}
          />

          <Container>
            {!enabled
              ? <StatusErrorMessage disabledReason={disabledReason} />
              : <DistributionForm
                sky_btc_exchange_rate={sky_btc_exchange_rate}
                balance={balance}
                intl={intl}

                address={skyAddress}
                handleAddressChange={this.handleAddressChange}

                drop_address={drop_address}
                status_address={status_address}
                handleStatusAddressChange={this.handleStatusAddressChange}

                getAddress={this.getAddress}
                addressLoading={addressLoading}

                checkStatus={this.checkStatus}
                statusLoading={statusLoading}
              />}
          </Container>
        </Wrapper>

        <Footer external />
      </div>
    );
  }
}

Distribution.propTypes = {
  intl: PropTypes.shape({
    formatMessage: PropTypes.func.isRequired,
  }).isRequired,
};

export default injectIntl(Distribution);
