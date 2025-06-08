import { api } from './index';

export const getProfile = () => {
  return api.get('/me');
};

export const getTransactions = () => {
  return api.get('/me/transactions');
}; 