import axios from "axios";

const PORTS = {
    auth: 'http://localhost:8002',
    order: 'http://localhost:8001',
    payment: 'http://localhost:8003',
    shipping: 'http://localhost:8004',
    warehouse: 'http://localhost:8005',
    tracking: 'http://localhost:8006',
    notification: 'http://localhost:8007',
};

const createClient = (baseURL) => {
    const client = axios.crate({ baseURL });
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