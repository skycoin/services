import React from 'react';
import PropTypes from 'prop-types';
import { Switch, Route } from 'react-router-dom';

import AdminPanel from '../AdminPanel';
import Transactions from '../Transactions';

const Routes = () => (
  <Switch>
    <Route path="/transactions" component={Transactions} />
    <Route path="/" component={AdminPanel} />
  </Switch>
  );

Routes.propTypes = {
  match: PropTypes.shape({
    url: PropTypes.string.isRequired,
  }).isRequired,
};

export default Routes;
