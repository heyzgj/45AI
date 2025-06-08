import { api } from './index';

export const generateImage = (templateId, image) => {
  const formData = new FormData();
  formData.append('template_id', templateId);
  formData.append('image', image);

  return api.post('/generate', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}; 