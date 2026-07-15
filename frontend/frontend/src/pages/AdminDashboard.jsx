import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Package, Shield, MapPin, Truck, PlusCircle, Compass, Info } from 'lucide-react';
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

      const unpaidOrWaiting = orderRes.data.filter(o => o.status === 'COMPLETED');
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
      const assignmentRes = await api.shipping.post('/shippings/assignments', {
        orderId: assignForm.orderId,
        courierId: assignForm.courierId,
      });

      const shipment = assignmentRes.data;

      const order = orders.find(o => o.id === assignForm.orderId);
      const trackingNumber = order.trackingNumber;

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
    <div className="min-h-screen bg-[#F0F0F0] flex flex-col font-sans text-[#333333]">
      {/* Brand Header Banner */}
      <header className="bg-[#4D148C] py-4 px-6 border-b-4 border-[#FF6600] flex justify-between items-center shadow-md">
        <div className="flex items-center gap-2">
          <h1 className="text-3xl font-black tracking-tighter">
            <span className="text-white">Papiton</span>
            <span className="text-[#FF6600] bg-white px-2 ml-1 rounded-sm">Admin</span>
          </h1>
        </div>
        <div className="flex items-center gap-4 text-xs text-white">
          <span className="font-mono">Global Management Suite</span>
          <button
            onClick={handleLogout}
            className="bg-[#FF6600] hover:bg-[#E05300] text-white px-3 py-1 font-bold rounded-sm border border-[#B34700]"
          >
            Log Out
          </button>
        </div>
      </header>

      {/* Main Container Layout */}
      <div className="flex-1 flex flex-col md:flex-row">
        {/* Sidebar */}
        <aside className="w-full md:w-64 bg-[#EAEAEA] border-r border-[#CCCCCC] p-4 flex flex-col justify-between">
          <div className="space-y-6">
            <div className="bg-white border border-[#CCCCCC] p-3 rounded-sm">
              <span className="text-[10px] text-slate-500 font-bold block uppercase mb-1">System Admin</span>
              <h4 className="font-bold text-[#4D148C] text-sm truncate">{user.name}</h4>
              <span className="inline-block mt-2 text-[9px] bg-red-950 text-red-400 border border-red-900 px-1.5 py-0.5 rounded-sm font-bold uppercase">
                {user.role}
              </span>
            </div>

            {/* Side Tabs navigation */}
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
                KELOLA PENGIRIMAN
              </button>
              <button
                onClick={() => setActiveTab('warehouses')}
                className={`w-full text-left py-2 px-3 text-xs font-bold rounded-sm border transition-all flex items-center gap-2 ${
                  activeTab === 'warehouses' 
                    ? 'bg-white border-[#4D148C] text-[#4D148C] shadow-sm' 
                    : 'bg-[#DFDFDF] border-[#CCCCCC] text-[#333333] hover:bg-white hover:border-[#CCCCCC]'
                }`}
              >
                <MapPin className="h-4 w-4 text-[#4D148C]" />
                KELOLA GUDANG (HUB)
              </button>
              <button
                onClick={() => setActiveTab('couriers')}
                className={`w-full text-left py-2 px-3 text-xs font-bold rounded-sm border transition-all flex items-center gap-2 ${
                  activeTab === 'couriers' 
                    ? 'bg-white border-[#4D148C] text-[#4D148C] shadow-sm' 
                    : 'bg-[#DFDFDF] border-[#CCCCCC] text-[#333333] hover:bg-white hover:border-[#CCCCCC]'
                }`}
              >
                <Truck className="h-4 w-4 text-[#4D148C]" />
                DAFTAR KURIR AKTIF
              </button>
              <button
                onClick={() => setActiveTab('routes')}
                className={`w-full text-left py-2 px-3 text-xs font-bold rounded-sm border transition-all flex items-center gap-2 ${
                  activeTab === 'routes' 
                    ? 'bg-white border-[#4D148C] text-[#4D148C] shadow-sm' 
                    : 'bg-[#DFDFDF] border-[#CCCCCC] text-[#333333] hover:bg-white hover:border-[#CCCCCC]'
                }`}
              >
                <Compass className="h-4 w-4 text-[#4D148C]" />
                RUTE LOGISTIK TRANSIT
              </button>
            </div>
          </div>
        </aside>

        {/* Content Workspace Panel */}
        <main className="flex-1 p-6 space-y-6">
          <div className="border-b border-[#CCCCCC] pb-3">
            <h2 className="text-xl font-black text-[#4D148C] tracking-tight">
              {activeTab === 'orders' && 'KELOLA PENGIRIMAN & PENUGASAN KURIR'}
              {activeTab === 'warehouses' && 'KELOLA HUB LOGISTIK (WAREHOUSES)'}
              {activeTab === 'couriers' && 'REGISTRASI & STATUS KURIR'}
              {activeTab === 'routes' && 'PENGATURAN RUTE TRANSIT EKSPEDISI'}
            </h2>
            <p className="text-xs text-slate-500 font-sans mt-0.5">Global admin console for Papiton Express network operations.</p>
          </div>

          {error && <div className="p-2.5 bg-[#FFF2F2] border border-[#FF9999] text-[#990000] text-xs font-semibold rounded-sm">{error}</div>}
          {success && <div className="p-2.5 bg-[#F2FFF2] border border-[#99FF99] text-[#006600] text-xs font-semibold rounded-sm">{success}</div>}

          {/* TAB 1: ORDERS & ASSIGNMENTS */}
          {activeTab === 'orders' && (
            <div className="space-y-6">
              {/* Courier Assignment form */}
              <div className="bg-white border border-[#CCCCCC] p-5 rounded-sm shadow-sm space-y-4">
                <h3 className="text-xs font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 flex items-center gap-1.5 uppercase">
                  <Truck className="h-4 w-4 text-[#FF6600]" />
                  Tugaskan Kurir Pengiriman
                </h3>
                <form onSubmit={handleAssignCourier} className="grid grid-cols-1 md:grid-cols-3 gap-4 items-end">
                  <div>
                    <label className="block text-[10px] font-bold text-slate-500 mb-1">Pilih Pesanan (Status PAID)</label>
                    <select
                      value={assignForm.orderId}
                      onChange={(e) => setAssignForm({ ...assignForm, orderId: e.target.value })}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    >
                      <option value="">-- Pilih order PAID --</option>
                      {orders.filter(o => o.status === 'COMPLETED').map(o => (
                        <option key={o.id} value={o.id}>
                          {o.trackingNumber} - Rp {o.totalPrice?.toLocaleString()}
                        </option>
                      ))}
                    </select>
                  </div>
                  <div>
                    <label className="block text-[10px] font-bold text-slate-500 mb-1">Pilih Kurir Aktif</label>
                    <select
                      value={assignForm.courierId}
                      onChange={(e) => setAssignForm({ ...assignForm, courierId: e.target.value })}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    >
                      <option value="">-- Pilih kurir AVAILABLE --</option>
                      {couriers.filter(c => c.status === 'AVAILABLE').map(c => (
                        <option key={c.id} value={c.id}>{c.name}</option>
                      ))}
                    </select>
                  </div>
                  <button
                    type="submit"
                    disabled={loading}
                    className="py-2 px-4 bg-[#FF6600] hover:bg-[#E05300] text-white text-xs font-bold rounded-sm border-b-2 border-[#B34700] shadow-sm transition-all"
                  >
                    Tugaskan Kurir
                  </button>
                </form>
              </div>

              {/* Order Lists table */}
              <div className="bg-white border border-[#CCCCCC] p-5 rounded-sm shadow-sm space-y-4">
                <h3 className="text-xs font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 uppercase">Monitoring Seluruh Pengiriman</h3>
                <div className="overflow-x-auto">
                  <table className="w-full text-left text-xs border-collapse">
                    <thead>
                      <tr className="bg-[#EAEAEA] border border-[#CCCCCC] text-[#4D148C]">
                        <th className="p-2 border border-[#CCCCCC]">Order ID</th>
                        <th className="p-2 border border-[#CCCCCC]">No. Resi</th>
                        <th className="p-2 border border-[#CCCCCC]">Rute Logistik</th>
                        <th className="p-2 border border-[#CCCCCC]">Total Biaya</th>
                        <th className="p-2 border border-[#CCCCCC]">Status Order</th>
                      </tr>
                    </thead>
                    <tbody>
                      {orders.map(o => (
                        <tr key={o.id} className="hover:bg-[#F9F9F9] border-b border-[#CCCCCC] text-slate-700">
                          <td className="p-2 border border-[#CCCCCC] font-mono text-[10px]">{o.id}</td>
                          <td className="p-2 border border-[#CCCCCC] font-mono font-bold text-blue-700 select-all">{o.trackingNumber}</td>
                          <td className="p-2 border border-[#CCCCCC]">
                            {o.addresses?.find(a => a.type === 'SENDER')?.city} ➔ {o.addresses?.find(a => a.type === 'RECEIVER')?.city}
                          </td>
                          <td className="p-2 border border-[#CCCCCC] font-semibold">Rp {o.totalPrice?.toLocaleString()}</td>
                          <td className="p-2 border border-[#CCCCCC]">
                            <span className={`inline-block px-2.5 py-0.5 rounded-sm text-[10px] font-bold ${
                              o.status === 'COMPLETED' ? 'bg-[#F2FFF2] text-[#006600] border border-[#99FF99]' : 'bg-[#EAEAEA] text-slate-600 border border-[#CCCCCC]'
                            }`}>
                              {o.status === 'COMPLETED' ? 'PAID' : o.status}
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

          {/* TAB 2: MANAGE WAREHOUSES */}
          {activeTab === 'warehouses' && (
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="lg:col-span-1 bg-white border border-[#CCCCCC] p-5 rounded-sm shadow-sm h-fit space-y-4">
                <h3 className="text-xs font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 flex items-center gap-1.5 uppercase">
                  <PlusCircle className="h-4 w-4 text-[#FF6600]" />
                  Daftarkan Hub Baru
                </h3>
                <form onSubmit={handleCreateWarehouse} className="space-y-4">
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Nama Gudang Hub</label>
                    <input
                      type="text"
                      required
                      value={warehouseForm.name}
                      onChange={(e) => setWarehouseForm({ ...warehouseForm, name: e.target.value })}
                      placeholder="Contoh: Surabaya Hub"
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Kota Lokasi</label>
                    <select
                      value={warehouseForm.city}
                      onChange={(e) => setWarehouseForm({ ...warehouseForm, city: e.target.value })}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    >
                      <option value="Jakarta">Jakarta</option>
                      <option value="Bandung">Bandung</option>
                      <option value="Surabaya">Surabaya</option>
                    </select>
                  </div>
                  <button
                    type="submit"
                    disabled={loading}
                    className="w-full py-2 bg-[#4D148C] hover:bg-[#390F66] text-white text-xs font-bold rounded-sm border-b-2 border-[#330D5C]"
                  >
                    DAFTARKAN GUDANG
                  </button>
                </form>
              </div>

              <div className="lg:col-span-2 bg-white border border-[#CCCCCC] p-5 rounded-sm shadow-sm space-y-4">
                <h3 className="text-xs font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 uppercase font-black">Daftar Hub Logistik</h3>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  {warehouses.map(w => (
                    <div key={w.id} className="p-3 bg-[#F9F9F9] border border-[#CCCCCC] rounded-sm space-y-1">
                      <h4 className="font-bold text-sm text-[#4D148C]">{w.name}</h4>
                      <p className="text-xs text-slate-500 flex items-center gap-1">
                        <MapPin className="h-3.5 w-3.5 text-[#FF6600]" />
                        Lokasi: {w.city}
                      </p>
                      <p className="text-[9px] text-slate-400 font-mono pt-1 border-t border-slate-200 mt-2 truncate">ID: {w.id}</p>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}

          {/* TAB 3: MANAGE COURIERS */}
          {activeTab === 'couriers' && (
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="lg:col-span-1 bg-white border border-[#CCCCCC] p-5 rounded-sm shadow-sm h-fit space-y-4">
                <h3 className="text-xs font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 flex items-center gap-1.5 uppercase">
                  <PlusCircle className="h-4 w-4 text-[#FF6600]" />
                  Daftar Kurir Baru
                </h3>
                <form onSubmit={handleCreateCourier} className="space-y-4">
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Nama Kurir</label>
                    <input
                      type="text"
                      required
                      value={courierForm.name}
                      onChange={(e) => setCourierForm({ ...courierForm, name: e.target.value })}
                      placeholder="Nama Kurir"
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">No. Handphone</label>
                    <input
                      type="text"
                      required
                      value={courierForm.phone}
                      onChange={(e) => setCourierForm({ ...courierForm, phone: e.target.value })}
                      placeholder="Contoh: 08123456789"
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    />
                  </div>
                  <button
                    type="submit"
                    disabled={loading}
                    className="w-full py-2 bg-[#4D148C] hover:bg-[#390F66] text-white text-xs font-bold rounded-sm border-b-2 border-[#330D5C]"
                  >
                    DAFTARKAN KURIR
                  </button>
                </form>
              </div>

              <div className="lg:col-span-2 bg-white border border-[#CCCCCC] p-5 rounded-sm shadow-sm space-y-4">
                <h3 className="text-xs font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 uppercase">Daftar Personel Kurir Aktif</h3>
                <div className="overflow-x-auto">
                  <table className="w-full text-left text-xs border-collapse">
                    <thead>
                      <tr className="bg-[#EAEAEA] border border-[#CCCCCC] text-[#4D148C]">
                        <th className="p-2 border border-[#CCCCCC]">Nama Kurir</th>
                        <th className="p-2 border border-[#CCCCCC]">No. Telepon</th>
                        <th className="p-2 border border-[#CCCCCC]">Status Ketersediaan</th>
                      </tr>
                    </thead>
                    <tbody>
                      {couriers.map(c => (
                        <tr key={c.id} className="hover:bg-[#F9F9F9] border-b border-[#CCCCCC]">
                          <td className="p-2 border border-[#CCCCCC] font-bold text-slate-700">{c.name}</td>
                          <td className="p-2 border border-[#CCCCCC]">{c.phone}</td>
                          <td className="p-2 border border-[#CCCCCC]">
                            <span className={`inline-block px-2.5 py-0.5 rounded-sm text-[10px] font-bold ${
                              c.status === 'AVAILABLE' ? 'bg-[#F2FFF2] text-[#006600] border border-[#99FF99]' : 'bg-[#FFF2F2] text-[#990000] border border-[#FF9999]'
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

          {/* TAB 4: ROUTING RULES */}
          {activeTab === 'routes' && (
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="lg:col-span-1 bg-white border border-[#CCCCCC] p-5 rounded-sm shadow-sm h-fit space-y-4">
                <h3 className="text-xs font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 flex items-center gap-1.5 uppercase">
                  <PlusCircle className="h-4 w-4 text-[#FF6600]" />
                  Atur Rute Transit
                </h3>
                <form onSubmit={handleCreateRoute} className="space-y-4">
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Gudang Asal (Origin Hub)</label>
                    <select
                      value={routeForm.originWarehouseId}
                      onChange={(e) => setRouteForm({ ...routeForm, originWarehouseId: e.target.value })}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    >
                      {warehouses.map(w => (
                        <option key={w.id} value={w.id}>{w.name} ({w.city})</option>
                      ))}
                    </select>
                  </div>
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Kota Tujuan Akhir Paket</label>
                    <select
                      value={routeForm.destinationCity}
                      onChange={(e) => setRouteForm({ ...routeForm, destinationCity: e.target.value })}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    >
                      <option value="Jakarta">Jakarta</option>
                      <option value="Bandung">Bandung</option>
                      <option value="Surabaya">Surabaya</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Transit Selanjutnya (Next Hub)</label>
                    <select
                      value={routeForm.nextWarehouseId}
                      onChange={(e) => setRouteForm({ ...routeForm, nextWarehouseId: e.target.value })}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
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
                    className="w-full py-2 bg-[#4D148C] hover:bg-[#390F66] text-white text-xs font-bold rounded-sm border-b-2 border-[#330D5C]"
                  >
                    SIMPAN ATURAN RUTE
                  </button>
                </form>
              </div>

              <div className="lg:col-span-2 bg-white border border-[#CCCCCC] p-5 rounded-sm shadow-sm space-y-4">
                <div className="flex items-start gap-2 p-3.5 bg-[#FFFDEB] border border-[#FFE7A3] text-slate-600 text-xs rounded-sm leading-relaxed">
                  <Info className="h-4 w-4 text-[#FF6600] flex-shrink-0 mt-0.5" />
                  <p>
                    <strong>Instruksi Alur Rute:</strong> Aturan rute ini digunakan oleh Hub Operator saat paket masuk ke Gudang Asal. Sistem akan menganalisis kota tujuan akhir paket untuk menentukan ke hub transit mana paket tersebut harus dikirimkan berikutnya.
                  </p>
                </div>
                <h3 className="text-xs font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 uppercase font-black">Aturan Rute Transit Aktif</h3>
                <p className="text-xs text-slate-500">Anda dapat mengatur aturan rute logistik pada panel kiri.</p>
              </div>
            </div>
          )}
        </main>
      </div>

      {/* Footer */}
      <footer className="bg-[#EAEAEA] border-t border-[#CCCCCC] py-3 text-center text-xs text-[#666666] font-sans">
        © 2003 Papiton Express Inc. All rights reserved. Master Admin Dashboard.
      </footer>
    </div>
  );
};

export default AdminDashboard;
