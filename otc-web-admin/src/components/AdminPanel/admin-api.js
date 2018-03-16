import axios from 'axios';

export const getStatus = () =>
  axios.get('/api/status')
    .then(response => response.data)
    .catch((error) => { throw new Error(error.response.data); });

export const setPrice = price =>
  axios.post('/api/price', { price }, {
    headers: {
      'Content-Type': 'application/json',
    },
  })
    .then(response => response.data)
    .catch((error) => {
      throw new Error(error.response.data || 'An unknown error occurred.');
    });

export const setSource = source =>
  axios.post('/api/source', { source }, {
    headers: {
      'Content-Type': 'application/json',
    },
  })
    .then(response => response.data)
    .catch((error) => {
      throw new Error(error.response.data || 'An unknown error occurred.');
    });

export const setOctState = pause =>
  axios.post('/api/pause', { pause }, {
    headers: {
      'Content-Type': 'application/json',
    },
  })
    .then(response => response.data)
    .catch((error) => { throw new Error(error.response.data); });

export const getHoldingBtc = () =>
  axios.get('/api/holding/btc', {
    headers: {
      'Content-Type': 'application/json',
    },
  })
  .then(response => response.data)
  .catch((error) => { throw new Error(error.response.data); });

export const getSkyAddresses = () =>
  axios.get('/api/addresses/sky', {
    headers: {
      'Content-Type': 'application/json',
    },
  })
  .then(response => response.data)
  .catch((error) => { throw new Error(error.response.data); });
