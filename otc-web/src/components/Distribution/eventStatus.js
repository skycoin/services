import moment from 'moment';
import { DISTRIBUTION_START, DISTRIBUTION_END } from 'config';

export const preEvent = moment().isBefore(DISTRIBUTION_START);
export const postEvent = moment().isAfter(DISTRIBUTION_END);
export const eventInProgress = !preEvent && !postEvent;
