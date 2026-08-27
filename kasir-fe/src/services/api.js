import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

export const barangApi = {
  getAll: (params) => api.get('/barang', { params }),
  getById: (id) => api.get(`/barang/${id}`),
  create: (data) => api.post('/barang', data),
  update: (id, data) => api.put(`/barang/${id}`, data),
  delete: (id) => api.delete(`/barang/${id}`)
}

export const kontakApi = {
  getAll: (params) => api.get('/kontaks', { params }),
  create: (data) => api.post('/kontaks', data),
  update: (id, data) => api.put(`/kontaks/${id}`, data),
  delete: (id) => api.delete(`/kontaks/${id}`)
}

export const transaksiApi = {
  mutasi: (data) => api.post('/transaksi-stok', data)
}

export const reportApi = {
  getTopBarang: () => api.get('/top-barang'),
  getTrackingStok: (params) => api.get('/tracking-stok', { params }),
  getStokPerBarang: () => api.get('/stok-per-barang'),
  getLowStockAlerts: () => api.get('/alerts/low-stock')
}
