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
    role: 'CUSTOMER', // Default Role untuk Registrasi
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
        // Alur Login
        const response = await api.auth.post('/auth/login', {
          email: formData.email,
          password: formData.password,
        });

        const { accessToken, user } = response.data;
        
        // Simpan sesi login ke LocalStorage
        localStorage.setItem('token', accessToken);
        localStorage.setItem('user', JSON.stringify(user));

        // Arahkan ke Dashboard sesuai Role masing-masing
        if (user.role === 'CUSTOMER') navigate('/customer');
        else if (user.role === 'COURIER') navigate('/courier');
        else if (user.role === 'WAREHOUSE_STAFF') navigate('/warehouse');
        else if (user.role === 'ADMIN') navigate('/admin');
      } else {
        // Alur Registrasi
        await api.auth.post('/auth/register', {
          name: formData.name,
          email: formData.email,
          password: formData.password,
          phone: formData.phone,
          role: formData.role,
        });

        setSuccess('Registrasi sukses! Silakan login.');
        setIsLogin(true); // Ganti tab otomatis ke login setelah sukses
        setFormData({ ...formData, password: '' }); // Bersihkan password
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Terjadi kesalahan sistem');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-950 px-4">
      <div className="max-w-md w-full bg-slate-900 border border-slate-800 rounded-2xl p-8 shadow-2xl">
        {/* Title */}
        <div className="text-center mb-8">
          <h2 className="text-3xl font-extrabold tracking-tight bg-gradient-to-r from-blue-500 to-indigo-400 bg-clip-text text-transparent">
            Papiton Express
          </h2>
          <p className="mt-2 text-sm text-slate-400">
            {isLogin ? 'Masuk ke akun logistik Anda' : 'Buat akun ekspedisi baru'}
          </p>
        </div>

        {/* Tab Toggle */}
        <div className="flex bg-slate-950 rounded-lg p-1 mb-6 border border-slate-800">
          <button
            onClick={() => { setIsLogin(true); setError(''); }}
            className={`flex-1 py-2 text-sm font-medium rounded-md transition-all ${isLogin ? 'bg-blue-600 text-white shadow' : 'text-slate-400 hover:text-white'}`}
          >
            Masuk
          </button>
          <button
            onClick={() => { setIsLogin(false); setError(''); }}
            className={`flex-1 py-2 text-sm font-medium rounded-md transition-all ${!isLogin ? 'bg-blue-600 text-white shadow' : 'text-slate-400 hover:text-white'}`}
          >
            Daftar
          </button>
        </div>

        {/* Alert Error / Success */}
        {error && <div className="mb-4 p-3 bg-red-950/50 border border-red-900 text-red-400 text-sm rounded-lg">{error}</div>}
        {success && <div className="mb-4 p-3 bg-green-950/50 border border-green-900 text-green-400 text-sm rounded-lg">{success}</div>}

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-4">
          {!isLogin && (
            <div>
              <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Nama Lengkap</label>
              <div className="relative">
                <User className="absolute left-3 top-3 h-5 w-5 text-slate-500" />
                <input
                  type="text"
                  name="name"
                  required
                  value={formData.name}
                  onChange={handleChange}
                  placeholder="Masukkan nama lengkap"
                  className="pl-10 w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                />
              </div>
            </div>
          )}

          <div>
            <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Alamat Email</label>
            <div className="relative">
              <Mail className="absolute left-3 top-3 h-5 w-5 text-slate-500" />
              <input
                type="email"
                name="email"
                required
                value={formData.email}
                onChange={handleChange}
                placeholder="email@domain.com"
                className="pl-10 w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
              />
            </div>
          </div>

          <div>
            <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Kata Sandi</label>
            <div className="relative">
              <Lock className="absolute left-3 top-3 h-5 w-5 text-slate-500" />
              <input
                type="password"
                name="password"
                required
                value={formData.password}
                onChange={handleChange}
                placeholder="••••••••"
                className="pl-10 w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
              />
            </div>
          </div>

          {!isLogin && (
            <>
              <div>
                <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Nomor Telepon</label>
                <div className="relative">
                  <Phone className="absolute left-3 top-3 h-5 w-5 text-slate-500" />
                  <input
                    type="text"
                    name="phone"
                    required
                    value={formData.phone}
                    onChange={handleChange}
                    placeholder="Contoh: 08123456789"
                    className="pl-10 w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                  />
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Peran Akun (Role)</label>
                <div className="relative">
                  <Shield className="absolute left-3 top-3 h-5 w-5 text-slate-500" />
                  <select
                    name="role"
                    value={formData.role}
                    onChange={handleChange}
                    className="pl-10 w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500 appearance-none"
                  >
                    <option value="CUSTOMER">Customer (Pengirim Paket)</option>
                    <option value="COURIER">Courier (Kurir Pengirim)</option>
                    <option value="WAREHOUSE_STAFF">Warehouse Staff (Gudang)</option>
                    <option value="ADMIN">System Administrator</option>
                  </select>
                </div>
              </div>
            </>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full mt-6 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white font-medium rounded-lg py-2.5 text-sm transition-all flex items-center justify-center gap-2 shadow-lg disabled:opacity-50"
          >
            {loading ? 'Memproses...' : isLogin ? 'Masuk' : 'Daftar Akun'}
            <ArrowRight className="h-4 w-4" />
          </button>
        </form>
      </div>
    </div>
  );
};

export default Login;