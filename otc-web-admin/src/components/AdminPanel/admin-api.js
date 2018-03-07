import axios from 'axios';

export const getStatus = () =>
  axios.get('/api/status')
    // .then(response => response.data)
    .then(() => ({
      prices: {
        internal: 150000,
        exchange: 119833,
        exchange_updated: 1519131184,
        internal_updated: 1519131184,
      },
      source: 'internal',
      paused: false,
    }))
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
  axios.post('/api/price', { source }, {
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
