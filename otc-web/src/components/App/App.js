import React from 'react';
import PropTypes from 'prop-types';
import flatten from 'flat';
import values from 'lodash/values';
import { Helmet } from 'react-helmet';
import { ThemeProvider } from 'styled-components';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import { FlagsProvider } from 'flag';
import { IntlProvider, addLocaleData } from 'react-intl';
import zh from 'react-intl/locale-data/zh';
import ru from 'react-intl/locale-data/ru';

import { COLORS, BREAKPOINTS, SPACE, FONT_SIZES, FLAGS } from '../../config';
import * as locales from '../../locales';

import Routes from '../Routes';

addLocaleData([...zh, ...ru]);

const theme = {
  colors: flatten(COLORS),
  breakpoints: values(BREAKPOINTS),
  space: SPACE,
  fontSizes: FONT_SIZES,
};

const Root = ({ locale, ...props }) => (
  <IntlProvider locale={locale} messages={flatten(locales[locale])}>
    <div>
      <Helmet titleTemplate="%s &middot; Skycoin">
        <html lang={locale} />
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
        <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png" />
        <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png" />
        <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png" />
        <link rel="manifest" href="/manifest.json" />
        <link rel="mask-icon" href="/safari-pinned-tab.svg" color="#8481eb" />
        <meta name="apple-mobile-web-app-title" content="Skycoin" />
        <meta name="application-name" content="Skycoin" />
        <meta name="theme-color" content="#ffffff" />
      </Helmet>

      <ThemeProvider theme={theme}>
        <Routes {...props} />
      </ThemeProvider>
    </div>
  </IntlProvider>
);

Root.propTypes = {
  locale: PropTypes.string.isRequired,
};

export default () => (
  <FlagsProvider flags={FLAGS}>
    <Router>
      <Switch>
        <Route path="/cn" render={props => <Root {...props} locale="zh" />} />
        <Route path="/ru" render={props => <Root {...props} locale="ru" />} />
        <Route path="/" render={props => <Root {...props} locale="en" />} />
      </Switch>
    </Router>
  </FlagsProvider>
);
