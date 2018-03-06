import axios from 'axios';

export const getStatus = () =>
  axios.get('/api/status')
    .then(response => response.data)
    .catch((error) => { throw new Error(error.response.data); });

export const setPrice = (price, source) =>
  axios.post('/api/price', { price, source }, {
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
