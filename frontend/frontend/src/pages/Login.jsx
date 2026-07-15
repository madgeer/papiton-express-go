import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Mail, Lock, User, Phone, Shield, ArrowRight } from 'lucide-react';
import { api } from '../services/api';

const Login = () => {
  const navigate = useNavigate();
  const [isLogin, setIsLogin] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);

  // Form State
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    phone: '',
    role: 'CUSTOMER',
  });

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
    setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      if (isLogin) {
        const response = await api.auth.post('/auth/login', {
          email: formData.email,
          password: formData.password,
        });

        const { accessToken, user } = response.data;
        
        localStorage.setItem('token', accessToken);
        localStorage.setItem('user', JSON.stringify(user));

        if (user.role === 'CUSTOMER') navigate('/customer');
        else if (user.role === 'COURIER') navigate('/courier');
        else if (user.role === 'WAREHOUSE_STAFF') navigate('/warehouse');
        else if (user.role === 'ADMIN') navigate('/admin');
      } else {
        await api.auth.post('/auth/register', {
          name: formData.name,
          email: formData.email,
          password: formData.password,
          phone: formData.phone,
          role: formData.role,
        });

        setSuccess('Registrasi sukses! Silakan login melalui tab Masuk.');
        setIsLogin(true);
        setFormData({ ...formData, password: '' });
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Terjadi kesalahan sistem');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#F0F0F0] flex flex-col justify-between">
      {/* Top Brand Banner Header */}
      <header className="bg-[#4D148C] py-4 px-6 border-b-4 border-[#FF6600] flex justify-between items-center shadow-md">
        <div className="flex items-center gap-2">
          {/* FedEx 2003 Logo style: Bold Sans-Serif, Purple & Orange */}
          <h1 className="text-3xl font-black tracking-tighter">
            <span className="text-white">Papiton</span>
            <span className="text-[#FF6600] bg-white px-2 ml-1 rounded-sm">Express</span>
          </h1>
        </div>
        <span className="text-xs text-white font-mono hidden sm:inline">Logistics Management Suite v2003</span>
      </header>

      {/* Main Login Box */}
      <div className="flex-1 flex items-center justify-center p-4 my-8">
        <div className="max-w-md w-full bg-white border-2 border-[#CCCCCC] rounded-sm p-6 shadow-md">
          {/* Header text */}
          <div className="border-b border-[#CCCCCC] pb-3 mb-4 flex justify-between items-center">
            <h2 className="text-lg font-bold text-[#4D148C]">
              {isLogin ? 'Login Pengguna' : 'Pendaftaran Pengguna Baru'}
            </h2>
            <span className="text-xs text-slate-500 font-mono">Secure Access</span>
          </div>

          {/* Web 1.0 Style Tab Toggles */}
          <div className="flex mb-6 border-b-2 border-[#4D148C]">
            <button
              onClick={() => { setIsLogin(true); setError(''); }}
              className={`py-2 px-6 text-sm font-bold border-t border-x rounded-t-sm transition-all ${
                isLogin 
                  ? 'bg-white border-[#4D148C] text-[#4D148C] translate-y-[2px] z-10' 
                  : 'bg-[#EAEAEA] border-[#CCCCCC] text-[#666666] hover:bg-[#F4F4F4]'
              }`}
            >
              Masuk
            </button>
            <button
              onClick={() => { setIsLogin(false); setError(''); }}
              className={`py-2 px-6 text-sm font-bold border-t border-x rounded-t-sm transition-all ${
                !isLogin 
                  ? 'bg-white border-[#4D148C] text-[#4D148C] translate-y-[2px] z-10' 
                  : 'bg-[#EAEAEA] border-[#CCCCCC] text-[#666666] hover:bg-[#F4F4F4]'
              }`}
            >
              Daftar Baru
            </button>
          </div>

          {/* Messages */}
          {error && <div className="mb-4 p-2.5 bg-[#FFF2F2] border border-[#FF9999] text-[#990000] text-xs font-semibold rounded-sm">{error}</div>}
          {success && <div className="mb-4 p-2.5 bg-[#F2FFF2] border border-[#99FF99] text-[#006600] text-xs font-semibold rounded-sm">{success}</div>}

          {/* Forms */}
          <form onSubmit={handleSubmit} className="space-y-4">
            {!isLogin && (
              <div>
                <label className="block text-xs font-bold text-[#333333] mb-1">Nama Lengkap</label>
                <input
                  type="text"
                  name="name"
                  required
                  value={formData.name}
                  onChange={handleChange}
                  placeholder="Nama Lengkap"
                  className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-2 px-3 text-sm text-[#333333] focus:outline-none focus:border-[#4D148C]"
                />
              </div>
            )}

            <div>
              <label className="block text-xs font-bold text-[#333333] mb-1">Alamat Email</label>
              <input
                type="email"
                name="email"
                required
                value={formData.email}
                onChange={handleChange}
                placeholder="contoh@mail.com"
                className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-2 px-3 text-sm text-[#333333] focus:outline-none focus:border-[#4D148C]"
              />
            </div>

            <div>
              <label className="block text-xs font-bold text-[#333333] mb-1">Kata Sandi (Password)</label>
              <input
                type="password"
                name="password"
                required
                value={formData.password}
                onChange={handleChange}
                placeholder="Password"
                className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-2 px-3 text-sm text-[#333333] focus:outline-none focus:border-[#4D148C]"
              />
            </div>

            {!isLogin && (
              <>
                <div>
                  <label className="block text-xs font-bold text-[#333333] mb-1">Nomor Telepon</label>
                  <input
                    type="text"
                    name="phone"
                    required
                    value={formData.phone}
                    onChange={handleChange}
                    placeholder="Nomor Telepon"
                    className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-2 px-3 text-sm text-[#333333] focus:outline-none focus:border-[#4D148C]"
                  />
                </div>

                <div>
                  <label className="block text-xs font-bold text-[#333333] mb-1">Peran Akses (Role)</label>
                  <select
                    name="role"
                    value={formData.role}
                    onChange={handleChange}
                    className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-2 px-3 text-sm text-[#333333] focus:outline-none focus:border-[#4D148C]"
                  >
                    <option value="CUSTOMER">Customer (Pelanggan)</option>
                    <option value="COURIER">Courier (Kurir)</option>
                    <option value="WAREHOUSE_STAFF">Warehouse Staff (Petugas Gudang)</option>
                    <option value="ADMIN">System Administrator</option>
                  </select>
                </div>
              </>
            )}

            {/* Beveled 2003 Submit Button Style */}
            <button
              type="submit"
              disabled={loading}
              className="w-full mt-4 bg-[#FF6600] hover:bg-[#E05300] text-white font-bold rounded-sm py-2 px-4 text-sm transition-all border-b-2 border-[#B34700] hover:border-[#993D00] shadow-sm flex items-center justify-center gap-2 disabled:opacity-50"
            >
              {loading ? 'Sedang Memproses...' : isLogin ? 'Masuk Sekarang' : 'Daftarkan Akun'}
              <ArrowRight className="h-4 w-4" />
            </button>
          </form>
        </div>
      </div>

      {/* Footer */}
      <footer className="bg-[#EAEAEA] border-t border-[#CCCCCC] py-3 text-center text-xs text-[#666666] font-sans">
        © 2003 Papiton Express Inc. All rights reserved. Global Trade Services, Tracking Systems.
      </footer>
    </div>
  );
};

export default Login;