import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import CustomerDashboard from './pages/CustomerDashboard';
import CourierDashboard from './pages/CourierDashboard';
import WarehouseDashboard from './pages/WarehouseDashboard';
import AdminDashboard from './pages/AdminDashboard';

// Helper untuk memproteksi halaman agar tidak bisa diakses tanpa login
const ProtectedRoute = ({ children, allowedRoles }) => {
  const token = localStorage.getItem('token');
  const userJson = localStorage.getItem('user');
  
  if (!token || !userJson) {
    return <Navigate to="/login" replace />;
  }

  const user = JSON.parse(userJson);
  if (allowedRoles && !allowedRoles.includes(user.role)) {
    return <Navigate to="/login" replace />;
  }

  return children;
};

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        
        {/* Halaman Proteksi Berdasarkan Role */}
        <Route path="/customer" element={
          <ProtectedRoute allowedRoles={['CUSTOMER']}><CustomerDashboard /></ProtectedRoute>
        } />
        <Route path="/courier" element={
          <ProtectedRoute allowedRoles={['COURIER']}><CourierDashboard /></ProtectedRoute>
        } />
        <Route path="/warehouse" element={
          <ProtectedRoute allowedRoles={['WAREHOUSE_STAFF']}><WarehouseDashboard /></ProtectedRoute>
        } />
        <Route path="/admin" element={
          <ProtectedRoute allowedRoles={['ADMIN']}><AdminDashboard /></ProtectedRoute>
        } />

        {/* Redirect default ke halaman login */}
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </Router>
  );
}

export default App;