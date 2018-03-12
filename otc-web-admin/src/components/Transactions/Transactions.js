import React from 'react';
import styled from 'styled-components';
import Helmet from 'react-helmet';
import { rem } from 'polished';
import { Flex, Box } from 'grid-styled';
import { DropdownList } from 'react-widgets';

import { COLORS, SPACE, BOX_SHADOWS, BORDER_RADIUS } from 'config';

import Text from 'components/Text';
import Header from 'components/Header';
import Footer from 'components/Footer';
import Container from 'components/Container';
import media from '../../utils/media';

import transactions from './transactions.json';

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
  waiting_confirm: 'Pending'
};
const transactionStatusToStr = s => transactionStatuses[s] || s;

const TransactionsFilter = () => (
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
          data={['All', 'Pending', 'Completed']} />
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
              <TableCell><DateTimeView dt={t.timestamps.deposited_at} /></TableCell>
              <TableCell><DateTimeView dt={t.timestamps.sent_at} /></TableCell>

              <TableCell>{transactionStatusToStr(t.status)}</TableCell>
              <TableCell>{t.drop.amount / 1e8}</TableCell>
              <TableCell>{t.rate.value / 1e8} </TableCell>
              <TableCell>{t.rate.source[0]} </TableCell>

              <TableCell>{t.address}</TableCell>
              <TableCell>{t.drop.address}</TableCell>
            </TableRow>))}
        </tbody>
      </Table>
    </div>
  );
};

export default class extends React.Component {
  render = () => {
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
            <TransactionsFilter />
            <TransactionsTable transactions={transactions} />
          </Container>
        </Wrapper>

        <Footer external />
      </div >
    );
  }
};
