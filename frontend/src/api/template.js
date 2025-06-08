import { api } from './index';

export const getTemplates = () => {
  return api.get('/templates');
};

export const getTemplateByID = (id) => {
  return api.get(`/templates/${id}`);
}; 