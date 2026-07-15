import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Package, Search, PlusCircle, CreditCard, LogOut, Clock, MapPin } from 'lucide-react';
import { api } from '../services/api';

const CustomerDashboard = () => {
  const navigate = useNavigate();
  const user = JSON.parse(localStorage.getItem('user') || '{}');
  
  // Tab control
  const [activeTab, setActiveTab] = useState('orders'); // 'orders', 'create', 'track'

  // State lists
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // Tracking State
  const [searchQuery, setSearchQuery] = useState('');
  const [trackingData, setTrackingData] = useState(null);
  const [trackingError, setTrackingError] = useState('');

  // Create Order Form State
  const [orderForm, setOrderForm] = useState({
    serviceTypeId: '550e8400-e29b-41d4-a716-446655440001', // Regular
    weight: 2.0,
    length: 10,
    width: 10,
    height: 10,
    senderName: user.name || '',
    senderPhone: user.phone || '',
    senderAddress: '',
    senderCity: 'Jakarta',
    senderProvince: 'DKI Jakarta',
    senderPostalCode: '',
    receiverName: '',
    receiverPhone: '',
    receiverAddress: '',
    receiverCity: 'Bandung',
    receiverProvince: 'Jawa Barat',
    receiverPostalCode: '',
  });

  const fetchOrders = async () => {
    setLoading(true);
    try {
      const response = await api.order.get('/orders');
      const customerOrders = response.data.filter(o => o.customerId === user.id);
      setOrders(customerOrders);
    } catch (err) {
      setError('Gagal memuat daftar pesanan');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, []);

  const handleLogout = () => {
    localStorage.clear();
    navigate('/login');
  };

  const handleOrderChange = (e) => {
    setOrderForm({ ...orderForm, [e.target.name]: e.target.value });
  };

  const handleCreateOrder = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');

    try {
      const payload = {
        customerId: user.id,
        serviceTypeId: orderForm.serviceTypeId,
        weight: parseFloat(orderForm.weight),
        length: parseFloat(orderForm.length),
        width: parseFloat(orderForm.width),
        height: parseFloat(orderForm.height),
        sender: {
          name: orderForm.senderName,
          phone: orderForm.senderPhone,
          address: orderForm.senderAddress,
          city: orderForm.senderCity,
          province: orderForm.senderProvince,
          postalCode: orderForm.senderPostalCode,
        },
        receiver: {
          name: orderForm.receiverName,
          phone: orderForm.receiverPhone,
          address: orderForm.receiverAddress,
          city: orderForm.receiverCity,
          province: orderForm.receiverProvince,
          postalCode: orderForm.receiverPostalCode,
        },
      };

      const response = await api.order.post('/orders', payload);
      const newOrder = response.data;

      // Automatically initialize tracking in background
      await api.tracking.post('/trackings', {
        orderId: newOrder.orderId,
        trackingNumber: newOrder.trackingNumber,
        initialStatus: 'WAITING_PAYMENT',
      });

      setSuccess(`Order berhasil dibuat! Resi Anda: ${newOrder.trackingNumber}`);
      
      setOrderForm({
        ...orderForm,
        weight: 2.0,
        length: 10,
        width: 10,
        height: 10,
        receiverName: '',
        receiverPhone: '',
        receiverAddress: '',
        receiverPostalCode: '',
      });

      fetchOrders();
      setActiveTab('orders');
    } catch (err) {
      setError(err.response?.data?.error || 'Gagal membuat order pengiriman');
    } finally {
      setLoading(false);
    }
  };

  const handleSimulatePayment = async (order) => {
    setLoading(true);
    setError('');
    setSuccess('');

    try {
      const invoiceRes = await api.payment.post('/payments', {
        orderId: order.id,
        amount: order.totalPrice,
        paymentMethod: 'SHOPEEPAY',
      });

      const paymentID = invoiceRes.data.id;

      await api.payment.post('/payments/webhook', {
        paymentId: paymentID,
        transactionId: 'TX-' + Math.random().toString(36).substring(2, 9).toUpperCase(),
        status: 'SUCCESS',
      });

      await api.order.post(`/orders/${order.id}/complete`);

      await api.tracking.post(`/trackings/${order.trackingNumber}/events`, {
        status: 'PAID',
        location: 'System Gateway',
        description: 'Payment successfully verified. Preparing package for courier pick up.',
      });

      setSuccess(`Pembayaran untuk order ${order.trackingNumber} sukses diverifikasi!`);
      fetchOrders();
    } catch (err) {
      setError('Gagal memproses simulasi pembayaran');
    } finally {
      setLoading(false);
    }
  };

  const handleTrackSearch = async (e) => {
    e.preventDefault();
    setTrackingError('');
    setTrackingData(null);

    if (!searchQuery) return;

    try {
      const response = await api.tracking.get(`/trackings/search?q=${searchQuery}`);
      setTrackingData(response.data);
    } catch (err) {
      setTrackingError('Nomor resi tidak terdaftar atau tidak ditemukan');
    }
  };

  return (
    <div className="min-h-screen bg-[#F0F0F0] flex flex-col font-sans text-[#333333]">
      {/* 2003 Brand Header Banner */}
      <header className="bg-[#4D148C] py-4 px-6 border-b-4 border-[#FF6600] flex justify-between items-center shadow-md">
        <div className="flex items-center gap-2">
          <h1 className="text-3xl font-black tracking-tighter">
            <span className="text-white">Papiton</span>
            <span className="text-[#FF6600] bg-white px-2 ml-1 rounded-sm">Express</span>
          </h1>
        </div>
        <div className="flex items-center gap-4 text-xs text-white">
          <span className="font-mono">Global Customer Portal</span>
          <button
            onClick={handleLogout}
            className="bg-[#FF6600] hover:bg-[#E05300] text-white px-3 py-1 font-bold rounded-sm border border-[#B34700]"
          >
            Log Out
          </button>
        </div>
      </header>

      {/* Main Grid: Sidebar + Workspace */}
      <div className="flex-1 flex flex-col md:flex-row">
        {/* Sidebar Container */}
        <aside className="w-full md:w-64 bg-[#EAEAEA] border-r border-[#CCCCCC] p-4 flex flex-col justify-between">
          <div className="space-y-6">
            {/* User Identity Info */}
            <div className="bg-white border border-[#CCCCCC] p-3 rounded-sm">
              <span className="text-[10px] text-slate-500 font-bold block uppercase mb-1">User Account</span>
              <h4 className="font-bold text-[#4D148C] text-sm truncate">{user.name}</h4>
              <span className="text-xs text-slate-600 block mt-0.5 truncate">{user.email}</span>
              <span className="inline-block mt-2 text-[9px] bg-slate-200 text-[#4D148C] font-bold px-1.5 py-0.5 rounded-sm border border-slate-300">
                {user.role}
              </span>
            </div>

            {/* Navigation Tabs */}
            <div className="space-y-1">
              <button
                onClick={() => setActiveTab('orders')}
                className={`w-full text-left py-2 px-3 text-xs font-bold rounded-sm border transition-all flex items-center gap-2 ${
                  activeTab === 'orders' 
                    ? 'bg-white border-[#4D148C] text-[#4D148C] shadow-sm' 
                    : 'bg-[#DFDFDF] border-[#CCCCCC] text-[#333333] hover:bg-white hover:border-[#CCCCCC]'
                }`}
              >
                <Package className="h-4 w-4 text-[#4D148C]" />
                PESANAN SAYA
              </button>
              <button
                onClick={() => setActiveTab('create')}
                className={`w-full text-left py-2 px-3 text-xs font-bold rounded-sm border transition-all flex items-center gap-2 ${
                  activeTab === 'create' 
                    ? 'bg-white border-[#4D148C] text-[#4D148C] shadow-sm' 
                    : 'bg-[#DFDFDF] border-[#CCCCCC] text-[#333333] hover:bg-white hover:border-[#CCCCCC]'
                }`}
              >
                <PlusCircle className="h-4 w-4 text-[#4D148C]" />
                KIRIM PAKET BARU
              </button>
              <button
                onClick={() => setActiveTab('track')}
                className={`w-full text-left py-2 px-3 text-xs font-bold rounded-sm border transition-all flex items-center gap-2 ${
                  activeTab === 'track' 
                    ? 'bg-white border-[#4D148C] text-[#4D148C] shadow-sm' 
                    : 'bg-[#DFDFDF] border-[#CCCCCC] text-[#333333] hover:bg-white hover:border-[#CCCCCC]'
                }`}
              >
                <Search className="h-4 w-4 text-[#4D148C]" />
                LACAK RESI
              </button>
            </div>
          </div>

          <div className="mt-8 border-t border-[#CCCCCC] pt-4 text-[10px] text-[#666666] leading-relaxed">
            Need help? Contact Global Support Center at 1-800-PAPITON.
          </div>
        </aside>

        {/* Content Panel */}
        <main className="flex-1 p-6">
          <div className="border-b border-[#CCCCCC] pb-3 mb-6">
            <h2 className="text-xl font-black text-[#4D148C] tracking-tight">
              {activeTab === 'orders' && 'RIWAYAT PENGIRIMAN'}
              {activeTab === 'create' && 'MENU TRANSAKSI: KIRIM PAKET'}
              {activeTab === 'track' && 'PELACAKAN NOMOR RESI'}
            </h2>
            <p className="text-xs text-slate-500 font-sans mt-0.5">Sistem integrasi logistik real-time Papiton Express.</p>
          </div>

          {/* Messages */}
          {error && <div className="mb-4 p-2.5 bg-[#FFF2F2] border border-[#FF9999] text-[#990000] text-xs font-semibold rounded-sm">{error}</div>}
          {success && <div className="mb-4 p-2.5 bg-[#F2FFF2] border border-[#99FF99] text-[#006600] text-xs font-semibold rounded-sm">{success}</div>}

          {/* TAB 1: LIST ORDERS */}
          {activeTab === 'orders' && (
            <div className="bg-white border border-[#CCCCCC] p-4 rounded-sm shadow-sm">
              {loading && <p className="text-xs text-slate-500">Memproses data...</p>}
              {!loading && orders.length === 0 && (
                <div className="text-center py-8">
                  <Package className="h-8 w-8 text-slate-400 mx-auto mb-2" />
                  <p className="text-xs text-slate-500">Tidak ada pengiriman aktif.</p>
                  <button
                    onClick={() => setActiveTab('create')}
                    className="mt-3 bg-[#FF6600] hover:bg-[#E05300] text-white text-xs font-bold py-1.5 px-4 rounded-sm border-b-2 border-[#B34700]"
                  >
                    Kirim Paket Sekarang
                  </button>
                </div>
              )}

              {!loading && orders.length > 0 && (
                <div className="overflow-x-auto">
                  <table className="w-full text-left text-xs border-collapse">
                    <thead>
                      <tr className="bg-[#EAEAEA] border border-[#CCCCCC] text-[#4D148C]">
                        <th className="p-2 border border-[#CCCCCC]">No. Resi</th>
                        <th className="p-2 border border-[#CCCCCC]">Kota Tujuan</th>
                        <th className="p-2 border border-[#CCCCCC]">Total Biaya</th>
                        <th className="p-2 border border-[#CCCCCC]">Status</th>
                        <th className="p-2 border border-[#CCCCCC] text-center">Aksi</th>
                      </tr>
                    </thead>
                    <tbody>
                      {orders.map((o) => (
                        <tr key={o.id} className="hover:bg-[#F9F9F9] border-b border-[#CCCCCC]">
                          <td className="p-2 border border-[#CCCCCC] font-mono font-bold text-blue-700 select-all">{o.trackingNumber}</td>
                          <td className="p-2 border border-[#CCCCCC]">{o.addresses?.find(a => a.type === 'RECEIVER')?.city}</td>
                          <td className="p-2 border border-[#CCCCCC] font-semibold">Rp {o.totalPrice?.toLocaleString()}</td>
                          <td className="p-2 border border-[#CCCCCC]">
                            <span className={`inline-block px-2 py-0.5 rounded-sm text-[10px] font-bold ${
                              o.status === 'COMPLETED' ? 'bg-[#F2FFF2] text-[#006600] border border-[#99FF99]' : 'bg-[#FFF2F2] text-[#990000] border border-[#FF9999]'
                            }`}>
                              {o.status}
                            </span>
                          </td>
                          <td className="p-2 border border-[#CCCCCC] text-center">
                            {o.status === 'WAITING_PAYMENT' ? (
                              <button
                                onClick={() => handleSimulatePayment(o)}
                                className="bg-[#008000] hover:bg-[#006600] text-white text-[10px] font-bold py-1 px-2.5 rounded-sm border-b-2 border-[#004d00]"
                              >
                                Simulasi Bayar
                              </button>
                            ) : (
                              <button
                                onClick={() => {
                                  setSearchQuery(o.trackingNumber);
                                  setActiveTab('track');
                                  setTrackingData(null);
                                }}
                                className="text-blue-700 font-bold hover:underline"
                              >
                                Lacak Detail
                              </button>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          )}

          {/* TAB 2: CREATE ORDER FORM */}
          {activeTab === 'create' && (
            <form onSubmit={handleCreateOrder} className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Sender Details */}
              <div className="bg-white border border-[#CCCCCC] p-4 rounded-sm shadow-sm space-y-3">
                <h3 className="text-sm font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 flex items-center gap-1.5">
                  <MapPin className="h-4 w-4 text-[#FF6600]" />
                  INFORMASI PENGIRIM (ORIGIN)
                </h3>
                <div>
                  <label className="block text-[11px] font-bold text-[#333333] mb-1">Nama Pengirim</label>
                  <input
                    type="text"
                    name="senderName"
                    required
                    value={orderForm.senderName}
                    onChange={handleOrderChange}
                    className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-[11px] font-bold text-[#333333] mb-1">Telepon Pengirim</label>
                  <input
                    type="text"
                    name="senderPhone"
                    required
                    value={orderForm.senderPhone}
                    onChange={handleOrderChange}
                    className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-[11px] font-bold text-[#333333] mb-1">Kota Asal</label>
                  <select
                    name="senderCity"
                    value={orderForm.senderCity}
                    onChange={handleOrderChange}
                    className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                  >
                    <option value="Jakarta">Jakarta</option>
                    <option value="Bandung">Bandung</option>
                    <option value="Surabaya">Surabaya</option>
                  </select>
                </div>
                <div>
                  <label className="block text-[11px] font-bold text-[#333333] mb-1">Alamat Lengkap</label>
                  <textarea
                    name="senderAddress"
                    required
                    rows="2"
                    value={orderForm.senderAddress}
                    onChange={handleOrderChange}
                    className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Provinsi</label>
                    <input
                      type="text"
                      name="senderProvince"
                      required
                      value={orderForm.senderProvince}
                      onChange={handleOrderChange}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Kode Pos</label>
                    <input
                      type="text"
                      name="senderPostalCode"
                      required
                      value={orderForm.senderPostalCode}
                      onChange={handleOrderChange}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    />
                  </div>
                </div>
              </div>

              {/* Receiver & Package Details */}
              <div className="space-y-6">
                <div className="bg-white border border-[#CCCCCC] p-4 rounded-sm shadow-sm space-y-3">
                  <h3 className="text-sm font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 flex items-center gap-1.5">
                    <MapPin className="h-4 w-4 text-[#FF6600]" />
                    INFORMASI PENERIMA (DESTINATION)
                  </h3>
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Nama Penerima</label>
                    <input
                      type="text"
                      name="receiverName"
                      required
                      value={orderForm.receiverName}
                      onChange={handleOrderChange}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Telepon Penerima</label>
                    <input
                      type="text"
                      name="receiverPhone"
                      required
                      value={orderForm.receiverPhone}
                      onChange={handleOrderChange}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Kota Tujuan</label>
                    <select
                      name="receiverCity"
                      value={orderForm.receiverCity}
                      onChange={handleOrderChange}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    >
                      <option value="Jakarta">Jakarta</option>
                      <option value="Bandung">Bandung</option>
                      <option value="Surabaya">Surabaya</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Alamat Lengkap</label>
                    <textarea
                      name="receiverAddress"
                      required
                      rows="2"
                      value={orderForm.receiverAddress}
                      onChange={handleOrderChange}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    />
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-[11px] font-bold text-[#333333] mb-1">Provinsi</label>
                      <input
                        type="text"
                        name="receiverProvince"
                        required
                        value={orderForm.receiverProvince}
                        onChange={handleOrderChange}
                        className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                      />
                    </div>
                    <div>
                      <label className="block text-[11px] font-bold text-[#333333] mb-1">Kode Pos</label>
                      <input
                        type="text"
                        name="receiverPostalCode"
                        required
                        value={orderForm.receiverPostalCode}
                        onChange={handleOrderChange}
                        className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                      />
                    </div>
                  </div>
                </div>

                {/* Package Specifications */}
                <div className="bg-white border border-[#CCCCCC] p-4 rounded-sm shadow-sm space-y-3">
                  <h3 className="text-sm font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 flex items-center gap-1.5">
                    <Package className="h-4 w-4 text-[#FF6600]" />
                    SPESIFIKASI PAKET
                  </h3>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-[11px] font-bold text-[#333333] mb-1">Jenis Layanan</label>
                      <select
                        name="serviceTypeId"
                        value={orderForm.serviceTypeId}
                        onChange={handleOrderChange}
                        className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                      >
                        <option value="550e8400-e29b-41d4-a716-446655440001">Regular (3 Hari)</option>
                        <option value="550e8400-e29b-41d4-a716-446655440002">NextDay (1 Hari)</option>
                      </select>
                    </div>
                    <div>
                      <label className="block text-[11px] font-bold text-[#333333] mb-1">Berat (Kg)</label>
                      <input
                        type="number"
                        step="0.1"
                        name="weight"
                        required
                        value={orderForm.weight}
                        onChange={handleOrderChange}
                        className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-2">
                    <div>
                      <label className="block text-[10px] font-bold text-slate-500 mb-1">Panjang (cm)</label>
                      <input type="number" name="length" value={orderForm.length} onChange={handleOrderChange} className="w-full bg-white border border-[#999999] rounded-sm p-1.5 text-xs text-[#333333]" />
                    </div>
                    <div>
                      <label className="block text-[10px] font-bold text-slate-500 mb-1">Lebar (cm)</label>
                      <input type="number" name="width" value={orderForm.width} onChange={handleOrderChange} className="w-full bg-white border border-[#999999] rounded-sm p-1.5 text-xs text-[#333333]" />
                    </div>
                    <div>
                      <label className="block text-[10px] font-bold text-slate-500 mb-1">Tinggi (cm)</label>
                      <input type="number" name="height" value={orderForm.height} onChange={handleOrderChange} className="w-full bg-white border border-[#999999] rounded-sm p-1.5 text-xs text-[#333333]" />
                    </div>
                  </div>

                  <button
                    type="submit"
                    disabled={loading}
                    className="w-full py-2 bg-[#FF6600] hover:bg-[#E05300] text-white text-xs font-bold rounded-sm border-b-2 border-[#B34700]"
                  >
                    {loading ? 'Menyimpan...' : 'KIRIM DAN HITUNG TARIF'}
                  </button>
                </div>
              </div>
            </form>
          )}

          {/* TAB 3: TRACKING TIMELINE */}
          {activeTab === 'track' && (
            <div className="bg-white border border-[#CCCCCC] p-4 rounded-sm shadow-sm space-y-6">
              {/* Search Form */}
              <form onSubmit={handleTrackSearch} className="flex gap-4">
                <input
                  type="text"
                  required
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Masukkan Nomor Resi (contoh: PPN-1234...)"
                  className="flex-1 bg-[#FFFFFF] border border-[#999999] rounded-sm py-2 px-3 text-xs focus:outline-none focus:border-[#4D148C]"
                />
                <button
                  type="submit"
                  className="bg-[#4D148C] hover:bg-[#390F66] text-white text-xs font-bold px-6 py-2 rounded-sm border-b-2 border-[#330D5C]"
                >
                  Cari Resi
                </button>
              </form>

              {trackingError && <div className="p-2.5 bg-[#FFF2F2] border border-[#FF9999] text-[#990000] text-xs font-semibold rounded-sm">{trackingError}</div>}

              {/* Stepper timeline FedEx 2003 style */}
              {trackingData && (
                <div className="border border-[#CCCCCC] rounded-sm p-4 bg-[#F9F9F9]">
                  <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center border-b border-[#CCCCCC] pb-3 mb-4">
                    <div>
                      <span className="text-[10px] text-slate-500 font-bold block uppercase">Tracking Number</span>
                      <h3 className="text-sm font-mono font-bold text-blue-700 select-all">{trackingData.trackingNumber}</h3>
                    </div>
                    <div className="mt-2 sm:mt-0">
                      <span className="text-[10px] text-slate-500 font-bold block uppercase">Current Location/Status</span>
                      <span className="inline-block px-2.5 py-0.5 rounded-sm text-[10px] font-bold bg-[#E9E1F5] text-[#4D148C] border border-[#D5C2EB]">
                        {trackingData.currentStatus}
                      </span>
                    </div>
                  </div>

                  {/* Web 1.0 Timeline: clean solid purple line, beveled boxes */}
                  <div className="relative border-l-2 border-[#4D148C] ml-3 space-y-6 py-2">
                    {trackingData.events?.map((evt, idx) => (
                      <div key={evt.id} className="relative pl-6">
                        {/* Dot indicator */}
                        <span className={`absolute -left-[7px] top-1.5 h-3 w-3 rounded-full flex items-center justify-center border ${
                          idx === trackingData.events.length - 1 ? 'bg-[#FF6600] border-[#FF6600]' : 'bg-[#4D148C] border-[#4D148C]'
                        }`}>
                          <span className="h-1 w-1 bg-white rounded-full"></span>
                        </span>

                        {/* Event details */}
                        <div className="bg-white border border-[#CCCCCC] p-3 rounded-sm">
                          <div className="flex justify-between items-start mb-1 text-xs">
                            <span className="font-bold text-[#4D148C]">{evt.status}</span>
                            <span className="text-[10px] text-slate-500 font-mono">
                              {new Date(evt.createdAt).toLocaleString('id-ID')}
                            </span>
                          </div>
                          <p className="text-xs text-slate-600">{evt.description}</p>
                          {evt.location && (
                            <span className="inline-block mt-1 text-[9px] bg-[#EAEAEA] text-slate-600 px-1.5 py-0.5 rounded-sm border border-slate-300">
                              Loc: {evt.location}
                            </span>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          )}
        </main>
      </div>

      {/* Footer */}
      <footer className="bg-[#EAEAEA] border-t border-[#CCCCCC] py-3 text-center text-xs text-[#666666] font-sans">
        © 2003 Papiton Express Inc. All rights reserved. Global Customer Portal.
      </footer>
    </div>
  );
};

export default CustomerDashboard;
