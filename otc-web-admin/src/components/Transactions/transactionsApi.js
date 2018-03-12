import axios from 'axios';

export const transactionFilters = {
  byState: {
    all: { name: 'all', url: '' },
    pending: { name: 'pending', url: '/pending' },
    completed: { name: 'completed', url: '/completed' },
  },
};

export const getTransactions = (filter = { state: transactionFilters.byState.all.name }) =>
  axios.get(`/api/transactions${transactionFilters.byState[filter.state].url}`)
    .then(response => response.data)
    .catch((error) => { throw new Error(error.response.data); });
