import React from 'react';
import styled from 'styled-components';
import Helmet from 'react-helmet';
import { rem } from 'polished';
import { Flex, Box } from 'grid-styled';
import { DropdownList } from 'react-widgets';

import { COLORS, SPACE } from 'config';

import Text from 'components/Text';
import Header from 'components/Header';
import Footer from 'components/Footer';
import Container from 'components/Container';
import media from '../../utils/media';

import { transactionFilters, getTransactions } from './transactionsApi';

const H3Styled = styled.h3`
  font-family: "Montreal", sans-serif;
  font-weight: 400;
  line-height: ${rem(1.75)};
`;

const Wrapper = styled.div`
  background-color: ${COLORS.gray[1]};
  padding: ${rem(SPACE[5])} 0;

  ${media.md.css`
    padding: ${rem(SPACE[7])} 0;
  `}
`;

const DateTimeView = ({ dt }) =>
  (<div>{new Date(dt * 1000).toLocaleDateString()} {new Date(dt * 1000).toLocaleTimeString()}</div>);

const transactionStatuses = {
  waiting_confirm: 'Pending',
  waiting_deposit: 'Pending',
};
const transactionStatusToStr = s => transactionStatuses[s] || s;
const transactionStateFilters = [
  { id: 'all', name: 'All' },
  { id: 'pending', name: 'Pending' },
  { id: 'completed', name: 'Completed' },
];

const TransactionsFilter = ({ changeFilterState }) => (
  <Flex column mb={5}>
    <Text as="h4">Filters</Text>
    <Flex row justify="flex-start" align="flex-end">
      <Box mr={5}>
        State:
      </Box>
      <Box>
        <DropdownList
          style={{ width: '150px' }}
          defaultValue="All"
          valueField="id"
          textField="name"
          onChange={({ id }) => changeFilterState(id)}
          data={transactionStateFilters} />
      </Box>
    </Flex>
  </Flex>
);

const TableHeadCell = styled.th`
  font-weight: bold;
  padding-top: ${rem(SPACE[3])};
  padding-bottom: ${rem(SPACE[3])};
  font-size: 14px;
  color: #fff;
  line-height: 1.4;
  background-color: ${COLORS.blue[5]};
`;

const TableCell = styled.td`
  font-size: 12px;
  color: #808080;
  line-height: 1.4;
  padding-left: ${rem(SPACE[6])};
  padding-right: ${rem(SPACE[3])};
`;

const TableRow = styled.tr`
  &:nth-child(even){
    background-color: #f8f6ff;
  }
`;

const Table = styled.table`
  ${TableHeadCell}, ${TableCell} {
    padding-top: ${rem(SPACE[4])};
    padding-bottom: ${rem(SPACE[4])};

    &:first-child {
      border-left: 0;
    }
    &:last-child {
      border-right: 0;
    }
  }

  ${TableRow}:last-child ${TableCell} {
    border-bottom: 0;
  }

  border-collapse: collapse;
`;

const TableHead = styled.thead`
`;

const TransactionsTable = ({ transactions }) => {
  return (
    <div>
      <H3Styled>Transactions:</H3Styled>
      <Table>
        <TableHead>
          <TableRow>
            <TableHeadCell>Created</TableHeadCell>
            <TableHeadCell>Deposited</TableHeadCell>
            <TableHeadCell>Sky Sent</TableHeadCell>

            <TableHeadCell>Status</TableHeadCell>
            <TableHeadCell>Amount</TableHeadCell>
            <TableHeadCell>Rate</TableHeadCell>
            <TableHeadCell>Source</TableHeadCell>

            <TableHeadCell>Sky Address</TableHeadCell>
            <TableHeadCell>BTC Address</TableHeadCell>
          </TableRow>
        </TableHead>
        <tbody>
          {transactions.map((t, i) =>
            (<TableRow key={i}>
              <TableCell><DateTimeView dt={t.timestamps.created_at} /></TableCell>
              <TableCell>
                {t.timestamps.deposited_at !== 0 && <DateTimeView dt={t.timestamps.deposited_at} />}
                {t.timestamps.deposited_at === 0 && 'N/A'}
              </TableCell>
              <TableCell>
                {t.timestamps.sent_at !== 0 && <DateTimeView dt={t.timestamps.sent_at} />}
                {t.timestamps.sent_at === 0 && 'N/A'}
              </TableCell>

              <TableCell>{transactionStatusToStr(t.status)}</TableCell>
              <TableCell>{t.drop.amount / 1e8}</TableCell>
              <TableCell>{t.rate && t.rate.value / 1e8} {!t.rate && 'N/A'} </TableCell>
              <TableCell>{t.rate && t.rate.source[0]} {!t.rate && 'N/A'} </TableCell>

              <TableCell>{t.address}</TableCell>
              <TableCell>{t.drop.address}</TableCell>
            </TableRow>))}
        </tbody>
      </Table>
    </div>
  );
};

const ActivityIndicator = () => (
  <div className="spinner">
    <div className="dot1"></div>
    <div className="dot2"></div>
  </div>
);

const EmptyDataSet = () => (
  <H3Styled>No transactions were found by specified filter.</H3Styled>
);

export default class extends React.Component {
  state = {
    transactions: [],
    loading: true,
    filter: {
      state: transactionFilters.byState.all.name
    }
  }
  loadData = async filter => {
    this.setState({
      ...this.state,
      loading: true,
    });
    const transactions = await getTransactions(filter);
    this.setState({
      ...this.state,
      transactions,
      loading: false,
    });
  }
  componentWillMount = async () => {
    await this.loadData(this.state.filter);
  }
  changeFilterState = async state => {
    const filter = { ...this.state.filter, state };
    this.setState({ ...this.state, filter });
    await this.loadData(filter);
  }
  render = () => {
    const { transactions, loading } = this.state;
    return (
      <div>
        <Helmet>
          <title>OTC Admin Panel</title>
        </Helmet>

        <Header external />

        <Wrapper>
          <Container style={{
            margin: '0 auto',
            maxWidth: '70rem',
          }}>
            <TransactionsFilter changeFilterState={this.changeFilterState} />
            {loading && <ActivityIndicator />}
            {!loading && transactions.length !== 0 && <TransactionsTable transactions={transactions} />}
            {!loading && transactions.length === 0 && <EmptyDataSet />}
          </Container>
        </Wrapper>

        <Footer external />
      </div >
    );
  }
};
