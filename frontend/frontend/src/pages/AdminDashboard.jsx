import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Package, Shield, MapPin, Truck, PlusCircle, Compass, LogOut, Info } from 'lucide-react';
import { api } from '../services/api';

const AdminDashboard = () => {
  const navigate = useNavigate();
  const user = JSON.parse(localStorage.getItem('user') || '{}');

  // Tab State
  const [activeTab, setActiveTab] = useState('orders'); // 'orders', 'warehouses', 'couriers', 'routes'

  // Data State
  const [orders, setOrders] = useState([]);
  const [warehouses, setWarehouses] = useState([]);
  const [couriers, setCouriers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // Form States
  const [warehouseForm, setWarehouseForm] = useState({ name: '', city: 'Jakarta' });
  const [courierForm, setCourierForm] = useState({ name: '', phone: '' });
  const [routeForm, setRouteForm] = useState({ originWarehouseId: '', destinationCity: 'Bandung', nextWarehouseId: '' });
  
  // Assignment Form State
  const [assignForm, setAssignForm] = useState({ orderId: '', courierId: '' });

  const fetchData = async () => {
    setLoading(true);
    try {
      const [orderRes, whRes, courierRes] = await Promise.all([
        api.order.get('/orders'),
        api.warehouse.get('/warehouses'),
        api.shipping.get('/shippings/couriers')
      ]);

      setOrders(orderRes.data);
      setWarehouses(whRes.data);
      setCouriers(courierRes.data);

      if (whRes.data.length > 0) {
        setRouteForm(prev => ({
          ...prev,
          originWarehouseId: whRes.data[0].id,
          nextWarehouseId: whRes.data.length > 1 ? whRes.data[1].id : ''
        }));
      }

      // Populate assignment defaults
      const unpaidOrWaiting = orderRes.data.filter(o => o.status === 'COMPLETED'); // PAID orders are marked COMPLETED in order service
      const availableCouriers = courierRes.data.filter(c => c.status === 'AVAILABLE');
      
      setAssignForm({
        orderId: unpaidOrWaiting.length > 0 ? unpaidOrWaiting[0].id : '',
        courierId: availableCouriers.length > 0 ? availableCouriers[0].id : ''
      });

    } catch (err) {
      setError('Gagal mengambil data sistem');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleLogout = () => {
    localStorage.clear();
    navigate('/login');
  };

  const handleCreateWarehouse = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');
    try {
      await api.warehouse.post('/warehouses', warehouseForm);
      setSuccess(`Warehouse ${warehouseForm.name} sukses dibuat!`);
      setWarehouseForm({ name: '', city: 'Jakarta' });
      fetchData();
    } catch (err) {
      setError('Gagal membuat gudang baru');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateCourier = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');
    try {
      await api.shipping.post('/shippings/couriers', courierForm);
      setSuccess(`Kurir ${courierForm.name} sukses didaftarkan!`);
      setCourierForm({ name: '', phone: '' });
      fetchData();
    } catch (err) {
      setError('Gagal mendaftarkan kurir baru');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateRoute = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');
    try {
      await api.warehouse.post('/warehouses/routes', {
        originWarehouseId: routeForm.originWarehouseId,
        destinationCity: routeForm.destinationCity,
        nextWarehouseId: routeForm.nextWarehouseId || null
      });
      setSuccess('Aturan rute logistik transit sukses disimpan!');
      fetchData();
    } catch (err) {
      setError('Gagal membuat rute baru. Pastikan rute asal dan tujuan unik.');
    } finally {
      setLoading(false);
    }
  };

  const handleAssignCourier = async (e) => {
    e.preventDefault();
    if (!assignForm.orderId || !assignForm.courierId) {
      setError('Pilih pesanan dan kurir terlebih dahulu');
      return;
    }
    setLoading(true);
    setError('');
    setSuccess('');

    try {
      // 1. Assign Courier in Shipping Service (creates shipment in PENDING_PICKUP, updates courier status to ON_DELIVERY)
      const assignmentRes = await api.shipping.post('/shippings/assignments', {
        orderId: assignForm.orderId,
        courierId: assignForm.courierId,
      });

      const shipment = assignmentRes.data;

      // 2. Fetch order to get trackingNumber
      const order = orders.find(o => o.id === assignForm.orderId);
      const trackingNumber = order.trackingNumber;

      // 3. Post a new status tracking event in Tracking Service
      const courier = couriers.find(c => c.id === assignForm.courierId);
      await api.tracking.post(`/trackings/${trackingNumber}/events`, {
        status: 'COURIER_ASSIGNED',
        location: 'Admin Panel',
        description: `Shipment assigned to Courier ${courier.name}. Shipment ID: ${shipment.id}`,
      });

      setSuccess(`Tugas pengiriman berhasil dibuat! Shipment ID: ${shipment.id}`);
      fetchData();
    } catch (err) {
      setError(err.response?.data?.error || 'Gagal menugaskan kurir ke order');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 flex flex-col md:flex-row">
      {/* Sidebar */}
      <aside className="w-full md:w-64 bg-slate-900 border-r border-slate-800 p-6 flex flex-col justify-between">
        <div>
          <div className="flex items-center gap-2 mb-8">
            <Shield className="h-6 w-6 text-blue-500" />
            <span className="font-extrabold text-xl bg-gradient-to-r from-blue-500 to-indigo-400 bg-clip-text text-transparent">
              Papiton Admin
            </span>
          </div>

          <div className="mb-6 p-4 bg-slate-950 rounded-xl border border-slate-800 space-y-2">
            <p className="text-xs text-slate-500 font-medium">Administrator:</p>
            <h4 className="font-semibold text-white truncate">{user.name}</h4>
            <span className="inline-block text-[10px] bg-red-950 text-red-400 border border-red-900 px-2 py-0.5 rounded-full font-bold uppercase">
              {user.role}
            </span>
          </div>

          <nav className="space-y-1">
            <button
              onClick={() => setActiveTab('orders')}
              className={`w-full text-left py-2.5 px-4 rounded-lg text-sm font-medium transition-all flex items-center gap-2 ${activeTab === 'orders' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'}`}
            >
              <Package className="h-4 w-4" />
              Kelola Pengiriman
            </button>
            <button
              onClick={() => setActiveTab('warehouses')}
              className={`w-full text-left py-2.5 px-4 rounded-lg text-sm font-medium transition-all flex items-center gap-2 ${activeTab === 'warehouses' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'}`}
            >
              <MapPin className="h-4 w-4" />
              Kelola Gudang (Warehouse)
            </button>
            <button
              onClick={() => setActiveTab('couriers')}
              className={`w-full text-left py-2.5 px-4 rounded-lg text-sm font-medium transition-all flex items-center gap-2 ${activeTab === 'couriers' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'}`}
            >
              <Truck className="h-4 w-4" />
              Daftar Kurir
            </button>
            <button
              onClick={() => setActiveTab('routes')}
              className={`w-full text-left py-2.5 px-4 rounded-lg text-sm font-medium transition-all flex items-center gap-2 ${activeTab === 'routes' ? 'bg-blue-600 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'}`}
            >
              <Compass className="h-4 w-4" />
              Rute Transit Logistik
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
        <header className="mb-8">
          <h1 className="text-2xl md:text-3xl font-extrabold text-white">Dashboard Master Sistem</h1>
          <p className="text-slate-400 text-sm">Gudang, Penugasan Kurir, dan Monitoring Status Logistik.</p>
        </header>

        {error && <div className="mb-6 p-4 bg-red-950/50 border border-red-900 text-red-400 text-sm rounded-xl">{error}</div>}
        {success && <div className="mb-6 p-4 bg-green-950/50 border border-green-900 text-green-400 text-sm rounded-xl">{success}</div>}

        {/* 1. ORDERS MONITORING & ASSIGNMENT TAB */}
        {activeTab === 'orders' && (
          <div className="space-y-8">
            {/* Courier Assignment Form */}
            <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md">
              <h3 className="text-lg font-bold text-white mb-4 flex items-center gap-2">
                <Truck className="h-5 w-5 text-blue-500" />
                Penugasan Kurir Pengiriman
              </h3>
              <form onSubmit={handleAssignCourier} className="grid grid-cols-1 md:grid-cols-3 gap-4 items-end">
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Pilih Pesanan (Sudah Bayar)</label>
                  <select
                    value={assignForm.orderId}
                    onChange={(e) => setAssignForm({ ...assignForm, orderId: e.target.value })}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                  >
                    <option value="">-- Tidak ada order berbayar --</option>
                    {orders.filter(o => o.status === 'COMPLETED').map(o => (
                      <option key={o.id} value={o.id}>
                        {o.trackingNumber} - Rp {o.totalPrice?.toLocaleString()}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Pilih Kurir Aktif</label>
                  <select
                    value={assignForm.courierId}
                    onChange={(e) => setAssignForm({ ...assignForm, courierId: e.target.value })}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                  >
                    <option value="">-- Tidak ada kurir siap --</option>
                    {couriers.filter(c => c.status === 'AVAILABLE').map(c => (
                      <option key={c.id} value={c.id}>{c.name}</option>
                    ))}
                  </select>
                </div>
                <button
                  type="submit"
                  disabled={loading}
                  className="bg-blue-600 hover:bg-blue-500 text-white font-semibold py-2.5 rounded-lg text-sm transition-all"
                >
                  Tugaskan Kurir
                </button>
              </form>
            </div>

            {/* System-wide Order Monitoring */}
            <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md">
              <h3 className="text-lg font-bold text-white mb-4">Semua Pesanan Pengiriman</h3>
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead>
                    <tr className="border-b border-slate-800 text-slate-400 font-semibold">
                      <th className="pb-3">Order ID</th>
                      <th className="pb-3">No. Resi</th>
                      <th className="pb-3">Rute</th>
                      <th className="pb-3">Harga</th>
                      <th className="pb-3">Status</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-800 text-slate-300">
                    {orders.map(o => (
                      <tr key={o.id}>
                        <td className="py-3 font-mono text-xs">{o.id}</td>
                        <td className="py-3 font-mono font-bold text-blue-400">{o.trackingNumber}</td>
                        <td className="py-3">
                          {o.addresses?.find(a => a.type === 'SENDER')?.city} ➔ {o.addresses?.find(a => a.type === 'RECEIVER')?.city}
                        </td>
                        <td className="py-3">Rp {o.totalPrice?.toLocaleString()}</td>
                        <td className="py-3">
                          <span className={`inline-flex px-2.5 py-0.5 rounded-full text-xs font-semibold ${
                            o.status === 'COMPLETED' ? 'bg-green-950 text-green-400 border border-green-900' : 'bg-slate-950 text-slate-400 border border-slate-800'
                          }`}>
                            {o.status === 'COMPLETED' ? 'PAID / COMPLETED' : o.status}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        )}

        {/* 2. MANAGE WAREHOUSES TAB */}
        {activeTab === 'warehouses' && (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-1 bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md h-fit space-y-4">
              <h3 className="text-lg font-bold text-white mb-2 flex items-center gap-2">
                <PlusCircle className="h-5 w-5 text-blue-500" />
                Daftarkan Gudang Baru
              </h3>
              <form onSubmit={handleCreateWarehouse} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Nama Gudang</label>
                  <input
                    type="text"
                    required
                    value={warehouseForm.name}
                    onChange={(e) => setWarehouseForm({ ...warehouseForm, name: e.target.value })}
                    placeholder="Contoh: Jakarta Hub"
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Kota Lokasi</label>
                  <select
                    value={warehouseForm.city}
                    onChange={(e) => setWarehouseForm({ ...warehouseForm, city: e.target.value })}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                  >
                    <option value="Jakarta">Jakarta</option>
                    <option value="Bandung">Bandung</option>
                    <option value="Surabaya">Surabaya</option>
                  </select>
                </div>
                <button
                  type="submit"
                  disabled={loading}
                  className="w-full py-2.5 bg-blue-600 hover:bg-blue-500 text-white font-medium rounded-lg text-sm transition-all"
                >
                  Daftarkan Gudang
                </button>
              </form>
            </div>

            <div className="lg:col-span-2 bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md">
              <h3 className="text-lg font-bold text-white mb-4">Daftar Titik Gudang</h3>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                {warehouses.map(w => (
                  <div key={w.id} className="p-4 bg-slate-950 border border-slate-800 rounded-xl">
                    <h4 className="font-bold text-white text-sm">{w.name}</h4>
                    <p className="text-xs text-slate-500 mt-1 flex items-center gap-1">
                      <MapPin className="h-3.5 w-3.5" />
                      Kota: {w.city}
                    </p>
                    <p className="text-[10px] text-slate-600 font-mono mt-2 truncate">ID: {w.id}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* 3. MANAGE COURIERS TAB */}
        {activeTab === 'couriers' && (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-1 bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md h-fit space-y-4">
              <h3 className="text-lg font-bold text-white mb-2 flex items-center gap-2">
                <PlusCircle className="h-5 w-5 text-blue-500" />
                Registrasi Kurir Baru
              </h3>
              <form onSubmit={handleCreateCourier} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Nama Kurir</label>
                  <input
                    type="text"
                    required
                    value={courierForm.name}
                    onChange={(e) => setCourierForm({ ...courierForm, name: e.target.value })}
                    placeholder="Nama kurir logistik"
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">No. Telepon</label>
                  <input
                    type="text"
                    required
                    value={courierForm.phone}
                    onChange={(e) => setCourierForm({ ...courierForm, phone: e.target.value })}
                    placeholder="Contoh: 081234..."
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                  />
                </div>
                <button
                  type="submit"
                  disabled={loading}
                  className="w-full py-2.5 bg-blue-600 hover:bg-blue-500 text-white font-medium rounded-lg text-sm transition-all"
                >
                  Registrasikan Kurir
                </button>
              </form>
            </div>

            <div className="lg:col-span-2 bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md">
              <h3 className="text-lg font-bold text-white mb-4">Daftar Kurir Aktif</h3>
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm">
                  <thead>
                    <tr className="border-b border-slate-800 text-slate-400 font-semibold">
                      <th className="pb-3">Nama</th>
                      <th className="pb-3">Telepon</th>
                      <th className="pb-3">Status</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-800 text-slate-300">
                    {couriers.map(c => (
                      <tr key={c.id}>
                        <td className="py-3 font-semibold">{c.name}</td>
                        <td className="py-3">{c.phone}</td>
                        <td className="py-3">
                          <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-semibold ${
                            c.status === 'AVAILABLE' ? 'bg-green-950 text-green-400 border border-green-900' : 'bg-yellow-950 text-yellow-400 border border-yellow-900'
                          }`}>
                            {c.status}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        )}

        {/* 4. MANAGE STATIC ROUTES TAB */}
        {activeTab === 'routes' && (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-1 bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md h-fit space-y-4">
              <h3 className="text-lg font-bold text-white mb-2 flex items-center gap-2">
                <PlusCircle className="h-5 w-5 text-blue-500" />
                Rute Baru
              </h3>
              <form onSubmit={handleCreateRoute} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Gudang Asal (Origin)</label>
                  <select
                    value={routeForm.originWarehouseId}
                    onChange={(e) => setRouteForm({ ...routeForm, originWarehouseId: e.target.value })}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                  >
                    {warehouses.map(w => (
                      <option key={w.id} value={w.id}>{w.name} ({w.city})</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Kota Tujuan Akhir</label>
                  <select
                    value={routeForm.destinationCity}
                    onChange={(e) => setRouteForm({ ...routeForm, destinationCity: e.target.value })}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                  >
                    <option value="Jakarta">Jakarta</option>
                    <option value="Bandung">Bandung</option>
                    <option value="Surabaya">Surabaya</option>
                  </select>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Transit Selanjutnya (Kosongkan jika Kirim Langsung)</label>
                  <select
                    value={routeForm.nextWarehouseId}
                    onChange={(e) => setRouteForm({ ...routeForm, nextWarehouseId: e.target.value })}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none"
                  >
                    <option value="">-- Kirim Langsung (Bukan Transit) --</option>
                    {warehouses.map(w => (
                      <option key={w.id} value={w.id}>{w.name} ({w.city})</option>
                    ))}
                  </select>
                </div>
                <button
                  type="submit"
                  disabled={loading}
                  className="w-full py-2.5 bg-blue-600 hover:bg-blue-500 text-white font-medium rounded-lg text-sm transition-all"
                >
                  Simpan Aturan Rute
                </button>
              </form>
            </div>

            <div className="lg:col-span-2 bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md">
              <div className="flex items-center gap-2 p-4 bg-slate-950 border border-slate-800 text-slate-400 text-xs rounded-xl mb-4 leading-relaxed">
                <Info className="h-5 w-5 text-indigo-500 flex-shrink-0" />
                Aturan rute logistik mendikte pergerakan paket. Jika paket masuk (*WAREHOUSE_IN*) ke gudang Asal dengan kota Tujuan akhir tertentu, sistem otomatis merekomendasikan untuk mengirim paket tersebut ke Warehouse Transit Selanjutnya.
              </div>
              <h3 className="text-lg font-bold text-white mb-4">Daftar Aturan Rute Logistik</h3>
              <p className="text-xs text-slate-500">Anda dapat mengatur aturan rute logistik pada panel kiri.</p>
            </div>
          </div>
        )}
      </main>
    </div>
  );
};

export default AdminDashboard;
