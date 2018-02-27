import axios from 'axios';

// Use Axios for HTTP requests. Refer to https://github.com/mzabriskie/axios
// for usage instructions. If the promises returned by #checkStatus or
// #getAddress reject, they should reject with an Error object containing
// a meaningful error message (will be shown to the user)
//
// export const checkStatus = skyAddress =>
//   axios.get(`https://fake.api/status?address=${skyAddress}`)
//     .catch(() => {
//       throw new Error(`Unable to check status for ${skyAddress}`)
//     });
//

export const getConfig = () =>
  axios.get('/api/config')
    .then(response => response.data);

export const checkStatus = skyAddress =>
  axios.get(`/api/status?skyaddr=${skyAddress}`)
    .then(response => response.data.statuses || [])
    .catch((error) => { throw new Error(error.response.data); });

export const getAddress = skyAddress =>
  axios.post('/api/bind', { skyaddr: skyAddress, coin_type: 'BTC' }, {
    headers: {
      'Content-Type': 'application/json',
    },
  })
    .then(response => response.data.deposit_address)
    .catch((error) => {
      throw new Error(error.response.data || 'An unknown error occurred.');
    });

export const checkExchangeStatus = () =>
  axios.get('/api/exchange-status')
    .then(response => response.data)
    .catch((error) => { throw new Error(error.response.data); });
