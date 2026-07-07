import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Package, Search, PlusCircle, CreditCard, LogOut, CheckCircle, Clock, Truck, MapPin } from 'lucide-react';
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
      // Filter orders by current customer ID
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
      
      // Reset form (except sender info)
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
      // 1. Create Invoice in Payment Service
      const invoiceRes = await api.payment.post('/payments', {
        orderId: order.id,
        amount: order.totalPrice,
        paymentMethod: 'SHOPEEPAY',
      });

      const paymentID = invoiceRes.data.id;

      // 2. Trigger Webhook Gateway callback simulating SUCCESS payment
      await api.payment.post('/payments/webhook', {
        paymentId: paymentID,
        transactionId: 'TX-' + Math.random().toString(36).substring(2, 9).toUpperCase(),
        status: 'SUCCESS',
      });

      // 3. Mark Order as complete (in standalone mode we mark it directly)
      await api.order.post(`/orders/${order.id}/complete`);

      // 4. Update tracking status
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
    <div className="min-h-screen bg-slate-950 flex flex-col md:flex-row">
      {/* Sidebar */}
      <aside className="w-full md:w-64 bg-slate-900 border-r border-slate-800 p-6 flex flex-col justify-between">
        <div>
          <div className="flex items-center gap-2 mb-8">
            <Package className="h-6 w-6 text-blue-500" />
            <span className="font-extrabold text-xl bg-gradient-to-r from-blue-500 to-indigo-400 bg-clip-text text-transparent">
              Papiton Express
            </span>
          </div>

          <div className="mb-6 p-4 bg-slate-950 rounded-xl border border-slate-800">
            <p className="text-xs text-slate-500 font-medium">Masuk Sebagai:</p>
            <h4 className="font-semibold text-white truncate">{user.name}</h4>
            <span className="inline-block mt-1 text-[10px] bg-blue-950 text-blue-400 border border-blue-900 px-2 py-0.5 rounded-full font-bold uppercase">
              {user.role}
            </span>
          </div>

          <nav className="space-y-1">
            <button
              onClick={() => setActiveTab('orders')}
              className={`w-full text-left py-2.5 px-4 rounded-lg text-sm font-medium transition-all flex items-center gap-2 ${activeTab === 'orders' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'}`}
            >
              <Package className="h-4 w-4" />
              Pesanan Saya
            </button>
            <button
              onClick={() => setActiveTab('create')}
              className={`w-full text-left py-2.5 px-4 rounded-lg text-sm font-medium transition-all flex items-center gap-2 ${activeTab === 'create' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'}`}
            >
              <PlusCircle className="h-4 w-4" />
              Kirim Paket baru
            </button>
            <button
              onClick={() => setActiveTab('track')}
              className={`w-full text-left py-2.5 px-4 rounded-lg text-sm font-medium transition-all flex items-center gap-2 ${activeTab === 'track' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'}`}
            >
              <Search className="h-4 w-4" />
              Lacak Resi
            </button>
          </nav>
        </div>

        <button
          onClick={handleLogout}
          className="mt-8 w-full py-2 px-4 border border-slate-800 hover:border-red-900 rounded-lg text-sm font-medium text-slate-400 hover:text-red-400 transition-all flex items-center justify-center gap-2"
        >
          <LogOut className="h-4 w-4" />
          Keluar
        </button>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 p-6 md:p-10 overflow-y-auto">
        {/* Header */}
        <header className="mb-8">
          <h1 className="text-2xl md:text-3xl font-extrabold text-white">
            {activeTab === 'orders' && 'Daftar Pengiriman'}
            {activeTab === 'create' && 'Kirim Paket Baru'}
            {activeTab === 'track' && 'Pelacakan Posisi Resi'}
          </h1>
          <p className="text-slate-400 text-sm">Kelola dan pantau barang kiriman Anda dalam satu dashboard.</p>
        </header>

        {/* Global Alert messages */}
        {error && <div className="mb-6 p-4 bg-red-950/50 border border-red-900 text-red-400 text-sm rounded-xl">{error}</div>}
        {success && <div className="mb-6 p-4 bg-green-950/50 border border-green-900 text-green-400 text-sm rounded-xl">{success}</div>}

        {/* 1. ORDERS TAB */}
        {activeTab === 'orders' && (
          <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md">
            {loading && <p className="text-slate-400 text-sm">Memuat data...</p>}
            {!loading && orders.length === 0 && (
              <div className="text-center py-10">
                <Package className="h-12 w-12 text-slate-600 mx-auto mb-4" />
                <h3 className="text-lg font-semibold text-slate-300">Belum Ada Pengiriman</h3>
                <p className="text-slate-500 text-sm mt-1">Anda belum membuat order pengiriman.</p>
                <button
                  onClick={() => setActiveTab('create')}
                  className="mt-4 inline-flex items-center gap-2 bg-blue-600 hover:bg-blue-500 text-white font-medium rounded-lg text-sm px-4 py-2"
                >
                  Buat Order Pertama
                </button>
              </div>
            )}

            {!loading && orders.length > 0 && (
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead>
                    <tr className="border-b border-slate-800 text-slate-400 font-semibold">
                      <th className="pb-4">No. Resi</th>
                      <th className="pb-4">Tujuan</th>
                      <th className="pb-4">Biaya Total</th>
                      <th className="pb-4">Status</th>
                      <th className="pb-4 text-center">Aksi</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-800">
                    {orders.map((o) => (
                      <tr key={o.id} className="text-slate-300">
                        <td className="py-4 font-mono font-bold text-blue-400">{o.trackingNumber}</td>
                        <td className="py-4 truncate max-w-xs">{o.addresses?.find(a => a.type === 'RECEIVER')?.city}</td>
                        <td className="py-4">Rp {o.totalPrice?.toLocaleString()}</td>
                        <td className="py-4">
                          <span className={`inline-flex px-2.5 py-0.5 rounded-full text-xs font-semibold ${
                            o.status === 'COMPLETED' ? 'bg-green-950 text-green-400 border border-green-900' : 'bg-yellow-950 text-yellow-400 border border-yellow-900'
                          }`}>
                            {o.status}
                          </span>
                        </td>
                        <td className="py-4 text-center">
                          {o.status === 'WAITING_PAYMENT' ? (
                            <button
                              onClick={() => handleSimulatePayment(o)}
                              className="inline-flex items-center gap-1.5 bg-green-600 hover:bg-green-500 text-white px-3 py-1 rounded-lg text-xs font-semibold"
                            >
                              <CreditCard className="h-3 w-3" />
                              Simulasi Bayar
                            </button>
                          ) : (
                            <button
                              onClick={() => {
                                setSearchQuery(o.trackingNumber);
                                setActiveTab('track');
                                setTrackingData(null);
                              }}
                              className="text-blue-400 hover:underline text-xs"
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

        {/* 2. CREATE ORDER TAB */}
        {activeTab === 'create' && (
          <form onSubmit={handleCreateOrder} className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* Sender Address */}
            <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 space-y-4">
              <h3 className="text-lg font-bold text-white mb-2 flex items-center gap-2 border-b border-slate-800 pb-3">
                <MapPin className="h-5 w-5 text-blue-500" />
                Alamat Pengirim
              </h3>
              <div>
                <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Nama Pengirim</label>
                <input
                  type="text"
                  name="senderName"
                  required
                  value={orderForm.senderName}
                  onChange={handleOrderChange}
                  className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Telepon Pengirim</label>
                <input
                  type="text"
                  name="senderPhone"
                  required
                  value={orderForm.senderPhone}
                  onChange={handleOrderChange}
                  className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Kota Asal</label>
                <select
                  name="senderCity"
                  value={orderForm.senderCity}
                  onChange={handleOrderChange}
                  className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                >
                  <option value="Jakarta">Jakarta</option>
                  <option value="Bandung">Bandung</option>
                  <option value="Surabaya">Surabaya</option>
                </select>
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Alamat Detail</label>
                <textarea
                  name="senderAddress"
                  required
                  rows="2"
                  value={orderForm.senderAddress}
                  onChange={handleOrderChange}
                  placeholder="Nama jalan, nomor rumah, RT/RW..."
                  className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Provinsi</label>
                  <input
                    type="text"
                    name="senderProvince"
                    required
                    value={orderForm.senderProvince}
                    onChange={handleOrderChange}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Kode Pos</label>
                  <input
                    type="text"
                    name="senderPostalCode"
                    required
                    value={orderForm.senderPostalCode}
                    onChange={handleOrderChange}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                  />
                </div>
              </div>
            </div>

            {/* Receiver Address & Package Info */}
            <div className="space-y-8">
              <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 space-y-4">
                <h3 className="text-lg font-bold text-white mb-2 flex items-center gap-2 border-b border-slate-800 pb-3">
                  <MapPin className="h-5 w-5 text-indigo-500" />
                  Alamat Penerima
                </h3>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Nama Penerima</label>
                  <input
                    type="text"
                    name="receiverName"
                    required
                    value={orderForm.receiverName}
                    onChange={handleOrderChange}
                    placeholder="Nama penerima paket"
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Telepon Penerima</label>
                  <input
                    type="text"
                    name="receiverPhone"
                    required
                    value={orderForm.receiverPhone}
                    onChange={handleOrderChange}
                    placeholder="Contoh: 0819876543"
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Kota Tujuan</label>
                  <select
                    name="receiverCity"
                    value={orderForm.receiverCity}
                    onChange={handleOrderChange}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                  >
                    <option value="Jakarta">Jakarta</option>
                    <option value="Bandung">Bandung</option>
                    <option value="Surabaya">Surabaya</option>
                  </select>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Alamat Detail</label>
                  <textarea
                    name="receiverAddress"
                    required
                    rows="2"
                    value={orderForm.receiverAddress}
                    onChange={handleOrderChange}
                    placeholder="Nama jalan, nomor rumah, RT/RW..."
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Provinsi</label>
                    <input
                      type="text"
                      name="receiverProvince"
                      required
                      value={orderForm.receiverProvince}
                      onChange={handleOrderChange}
                      className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Kode Pos</label>
                    <input
                      type="text"
                      name="receiverPostalCode"
                      required
                      value={orderForm.receiverPostalCode}
                      onChange={handleOrderChange}
                      className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                    />
                  </div>
                </div>
              </div>

              {/* Package Detail */}
              <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 space-y-4">
                <h3 className="text-lg font-bold text-white mb-2 flex items-center gap-2 border-b border-slate-800 pb-3">
                  <Package className="h-5 w-5 text-green-500" />
                  Detail Barang
                </h3>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Jenis Layanan</label>
                    <select
                      name="serviceTypeId"
                      value={orderForm.serviceTypeId}
                      onChange={handleOrderChange}
                      className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                    >
                      <option value="550e8400-e29b-41d4-a716-446655440001">Regular (3 Hari)</option>
                      <option value="550e8400-e29b-41d4-a716-446655440002">NextDay (1 Hari)</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Berat (Kg)</label>
                    <input
                      type="number"
                      step="0.1"
                      name="weight"
                      required
                      value={orderForm.weight}
                      onChange={handleOrderChange}
                      className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-3 gap-2">
                  <div>
                    <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1">Panjang (cm)</label>
                    <input type="number" name="length" value={orderForm.length} onChange={handleOrderChange} className="w-full bg-slate-950 border border-slate-800 rounded-lg p-2 text-sm text-white" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1">Lebar (cm)</label>
                    <input type="number" name="width" value={orderForm.width} onChange={handleOrderChange} className="w-full bg-slate-950 border border-slate-800 rounded-lg p-2 text-sm text-white" />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1">Tinggi (cm)</label>
                    <input type="number" name="height" value={orderForm.height} onChange={handleOrderChange} className="w-full bg-slate-950 border border-slate-800 rounded-lg p-2 text-sm text-white" />
                  </div>
                </div>

                <button
                  type="submit"
                  disabled={loading}
                  className="w-full py-3 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white font-medium rounded-lg text-sm transition-all"
                >
                  {loading ? 'Menyimpan...' : 'Konfirmasi Pemesanan'}
                </button>
              </div>
            </div>
          </form>
        )}

        {/* 3. TRACKING TAB */}
        {activeTab === 'track' && (
          <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md space-y-6">
            {/* Search Input */}
            <form onSubmit={handleTrackSearch} className="flex gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-3 h-5 w-5 text-slate-500" />
                <input
                  type="text"
                  required
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Masukkan Nomor Resi (contoh: PPN-123456...)"
                  className="pl-10 w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                />
              </div>
              <button
                type="submit"
                className="bg-blue-600 hover:bg-blue-500 text-white px-6 py-2.5 rounded-lg text-sm font-semibold transition-all"
              >
                Cari Resi
              </button>
            </form>

            {/* Error Alert */}
            {trackingError && <div className="p-4 bg-red-950/50 border border-red-900 text-red-400 text-sm rounded-xl">{trackingError}</div>}

            {/* Tracking Result Visualizer */}
            {trackingData && (
              <div className="border border-slate-800 rounded-xl p-6 bg-slate-950">
                <div className="flex flex-col md:flex-row justify-between border-b border-slate-800 pb-4 mb-6">
                  <div>
                    <span className="text-xs text-slate-500 font-semibold tracking-wider uppercase">Nomor Resi</span>
                    <h3 className="text-lg font-mono font-bold text-blue-400">{trackingData.trackingNumber}</h3>
                  </div>
                  <div className="mt-2 md:mt-0">
                    <span className="text-xs text-slate-500 font-semibold tracking-wider uppercase">Status Terkini</span>
                    <div>
                      <span className="inline-flex px-2.5 py-0.5 rounded-full text-xs font-semibold bg-blue-950 text-blue-400 border border-blue-900">
                        {trackingData.currentStatus}
                      </span>
                    </div>
                  </div>
                </div>

                {/* Timeline Stepper */}
                <div className="relative border-l-2 border-slate-800 ml-3 space-y-8 py-2">
                  {trackingData.events?.map((evt, idx) => (
                    <div key={evt.id} className="relative pl-8">
                      {/* Timeline Dot */}
                      <span className={`absolute -left-[9px] top-1.5 h-4 w-4 rounded-full flex items-center justify-center border-2 ${
                        idx === trackingData.events.length - 1 ? 'bg-blue-500 border-blue-400 animate-ping' : 'bg-slate-950 border-slate-700'
                      }`}>
                        <span className={`h-1.5 w-1.5 rounded-full ${idx === trackingData.events.length - 1 ? 'bg-white' : 'bg-slate-700'}`}></span>
                      </span>

                      {/* Content Card */}
                      <div className="bg-slate-900 border border-slate-800 p-4 rounded-xl">
                        <div className="flex flex-col md:flex-row justify-between gap-1 mb-2">
                          <h4 className="font-bold text-white text-sm">{evt.status}</h4>
                          <span className="text-[11px] text-slate-500 flex items-center gap-1 font-mono">
                            <Clock className="h-3 w-3" />
                            {new Date(evt.createdAt).toLocaleString('id-ID')}
                          </span>
                        </div>
                        <p className="text-xs text-slate-400 font-medium">{evt.description}</p>
                        {evt.location && (
                          <span className="inline-flex mt-2 items-center gap-1 text-[10px] bg-slate-950 text-slate-500 border border-slate-800 px-2 py-0.5 rounded-full font-semibold">
                            <MapPin className="h-2.5 w-2.5" />
                            {evt.location}
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
  );
};

export default CustomerDashboard;
