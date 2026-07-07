import axios from 'axios';

const PORTS = {
  auth: 'http://localhost:8082',
  order: 'http://localhost:8081',
  payment: 'http://localhost:8083',
  shipping: 'http://localhost:8084',
  warehouse: 'http://localhost:8085',
  tracking: 'http://localhost:8086',
  notification: 'http://localhost:8087',
};

const createClient = (baseURL) => {
  const client = axios.create({ baseURL });
  client.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });
  return client;
};

export const api = {
  auth: createClient(PORTS.auth),
  order: createClient(PORTS.order),
  payment: createClient(PORTS.payment),
  shipping: createClient(PORTS.shipping),
  warehouse: createClient(PORTS.warehouse),
  tracking: createClient(PORTS.tracking),
  notification: createClient(PORTS.notification),
};