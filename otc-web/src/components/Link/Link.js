import React from 'react';
import PropTypes from 'prop-types';
import { Link as RouterLink, withRouter } from 'react-router-dom';
import join from 'join-path';
import omit from 'lodash/omit';
import trimEnd from 'lodash/trimEnd';

const getURL = (match, url) => trimEnd(
  url.includes('://')
    ? url
    : join('/', match.params.locale, url, url.includes('#') ? '' : '/'),
  );

const filterProps = props =>
  omit(props, ['location', 'history', 'staticContext', 'pill', 'outlined',
    'white', 'bg', 'big', 'color', 'fontSize', 'm', 'ml', 'mr', 'mt', 'mb']);

const Link = ({ to, href, match, children, ...props }) => {
  if (to) {
    return (
      <RouterLink to={getURL(match, to)} {...filterProps(props)}>
        {children}
      </RouterLink>
    );
  }

  if (href) {
    return <a href={getURL(match, href)} {...filterProps(props)}>{children}</a>;
  }

  return <a {...filterProps(props)}>{children}</a>;
};

Link.propTypes = {
  to: PropTypes.string,
  href: PropTypes.string,
  match: PropTypes.shape({
    params: PropTypes.shape({
      locale: PropTypes.string,
    }),
  }).isRequired,
  children: PropTypes.element.isRequired,
};

Link.defaultProps = {
  to: undefined,
  href: undefined,
};

export default withRouter(Link);
