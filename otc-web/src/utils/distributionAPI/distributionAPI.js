import axios from 'axios';

export const checkStatus = ({ drop_address, drop_currency }) =>
  axios.post('/api/status', { drop_address, drop_currency })
    .then(response => [response.data])
    .catch((error) => { throw new Error(error.response.data); });

export const getAddress = skyAddress =>
  axios.post('/api/bind', { address: skyAddress, drop_currency: 'BTC' }, {
    headers: {
      'Content-Type': 'application/json',
    },
  })
    .then(response => response.data)
    .catch((error) => {
      throw new Error(error.response.data || 'An unknown error occurred.');
    });

export const checkExchangeStatus = () =>
  axios.get('/api/exchange-status')
    .then(response => response.data)
    .catch((error) => { throw new Error(error.response.data); });
