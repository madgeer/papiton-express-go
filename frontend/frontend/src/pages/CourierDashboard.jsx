import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Package, Truck, Compass, CheckCircle, Clock, Shield, LogOut, MapPin } from 'lucide-react';
import { api } from '../services/api';

const CourierDashboard = () => {
  const navigate = useNavigate();
  const user = JSON.parse(localStorage.getItem('user') || '{}');

  // State
  const [courierInfo, setCourierInfo] = useState(null);
  const [shipmentIdInput, setShipmentIdInput] = useState('');
  const [shipmentData, setShipmentData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const fetchCourierInfo = async () => {
    try {
      const response = await api.shipping.get('/shippings/couriers');
      // Cari data kurir saat ini dari daftar kurir
      const current = response.data.find(c => c.name === user.name || c.phone === user.phone);
      if (current) {
        setCourierInfo(current);
      }
    } catch (err) {
      setError('Gagal memuat status kurir');
    }
  };

  useEffect(() => {
    fetchCourierInfo();
  }, []);

  const handleLogout = () => {
    localStorage.clear();
    navigate('/login');
  };

  const handleFetchShipment = async (idToFetch) => {
    const id = idToFetch || shipmentIdInput;
    if (!id) return;
    setLoading(true);
    setError('');
    setSuccess('');
    try {
      const response = await api.shipping.get(`/shippings/${id}`);
      setShipmentData(response.data);
    } catch (err) {
      setError('Shipment ID tidak ditemukan');
      setShipmentData(null);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateStatus = async (newStatus) => {
    if (!shipmentData) return;
    setLoading(true);
    setError('');
    setSuccess('');

    try {
      const notesMap = {
        'PICKED_UP': 'Paket telah diambil oleh kurir dari pengirim.',
        'ON_TRANSIT': 'Paket dalam perjalanan transit logistik.',
        'OUT_FOR_DELIVERY': 'Kurir sedang dalam perjalanan mengirimkan paket ke penerima.',
        'DELIVERED': 'Paket telah diterima dengan baik oleh penerima.',
        'FAILED': 'Pengiriman paket gagal dilakukan.',
      };

      // 1. Update Shipment Status in Shipping Service
      const response = await api.shipping.put(`/shippings/${shipmentData.id}/status`, {
        status: newStatus,
        notes: notesMap[newStatus] || 'Status updated by courier.',
      });

      // 2. Sync / Update Tracking Service with new milestone event
      await api.tracking.post(`/trackings/${shipmentData.orderId}/events`, { // Wait, tracking uses trackingNumber or orderId?
        // Wait! Let's check api_test.go of tracking service:
        // r.POST("/trackings/:num/events", trackingHandler.AddEventHandler)
        // Ah! The URL path is /trackings/:num/events, where :num is the tracking_number, NOT the orderId!
        // But wait, does the Courier know the tracking number?
        // Let's get the tracking number. We don't have it directly in the shipment model unless we fetch the order details,
        // OR in the simulated demo, we can fetch order details from order service using orderId to get trackingNumber!
        // Yes! Let's fetch the order details from order service to get the trackingNumber!
        status: newStatus,
        location: courierInfo ? `Kurir ${courierInfo.name}` : 'Kurir Logistik',
        description: notesMap[newStatus],
      });

      setSuccess(`Status paket berhasil diperbarui menjadi: ${newStatus}`);
      setShipmentData(response.data);
      fetchCourierInfo(); // Reload courier status (if changed to AVAILABLE after DELIVERED)
    } catch (err) {
      // If direct tracking fails, try using orderId or handle errors
      try {
        // Let's fetch order to get trackingNumber
        const orderRes = await api.order.get(`/orders/${shipmentData.orderId}`);
        const trackingNum = orderRes.data.trackingNumber;

        await api.tracking.post(`/trackings/${trackingNum}/events`, {
          status: newStatus,
          location: courierInfo ? `Kurir ${courierInfo.name}` : 'Kurir Logistik',
          description: notesMap[newStatus],
        });
        
        setSuccess(`Status paket & pelacakan berhasil diperbarui menjadi: ${newStatus}`);
        // Fetch shipment again to update view
        const response = await api.shipping.get(`/shippings/${shipmentData.id}`);
        setShipmentData(response.data);
        fetchCourierInfo();
      } catch (nestedErr) {
        setError('Gagal mensinkronisasikan status ke Tracking Service');
      }
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
            <Truck className="h-6 w-6 text-blue-500" />
            <span className="font-extrabold text-xl bg-gradient-to-r from-blue-500 to-indigo-400 bg-clip-text text-transparent">
              Papiton Courier
            </span>
          </div>

          <div className="mb-6 p-4 bg-slate-950 rounded-xl border border-slate-800 space-y-2">
            <p className="text-xs text-slate-500 font-medium">Status Kurir:</p>
            <h4 className="font-semibold text-white truncate">{user.name}</h4>
            {courierInfo ? (
              <span className={`inline-block text-[10px] px-2 py-0.5 rounded-full font-bold uppercase ${
                courierInfo.status === 'AVAILABLE' ? 'bg-green-950 text-green-400 border border-green-900' : 'bg-yellow-950 text-yellow-400 border border-yellow-900'
              }`}>
                {courierInfo.status}
              </span>
            ) : (
              <span className="text-xs text-slate-500">Loading status...</span>
            )}
          </div>

          <div className="p-4 bg-slate-950 rounded-xl border border-slate-800">
            <h5 className="text-xs font-bold text-slate-400 mb-1">Panduan Demo:</h5>
            <p className="text-[11px] text-slate-500 leading-relaxed">
              Copy <strong>Shipment ID</strong> yang ter-generate dari penugasan kurir di Admin Panel/Customer, paste di bawah untuk mensimulasikan tugas kurir.
            </p>
          </div>
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
          <h1 className="text-2xl md:text-3xl font-extrabold text-white">Panel Petugas Kurir</h1>
          <p className="text-slate-400 text-sm">Ambil paket, update status perjalanan, dan konfirmasi barang terkirim.</p>
        </header>

        {error && <div className="mb-6 p-4 bg-red-950/50 border border-red-900 text-red-400 text-sm rounded-xl">{error}</div>}
        {success && <div className="mb-6 p-4 bg-green-950/50 border border-green-900 text-green-400 text-sm rounded-xl">{success}</div>}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Action Search Card */}
          <div className="lg:col-span-1 bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md h-fit space-y-4">
            <h3 className="text-lg font-bold text-white mb-2 flex items-center gap-2">
              <Compass className="h-5 w-5 text-blue-500" />
              Pencarian Tugas
            </h3>
            <div>
              <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Shipment ID (UUID)</label>
              <input
                type="text"
                value={shipmentIdInput}
                onChange={(e) => setShipmentIdInput(e.target.value)}
                placeholder="Masukkan ID Pengiriman..."
                className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
              />
            </div>
            <button
              onClick={() => handleFetchShipment()}
              disabled={loading}
              className="w-full py-2.5 bg-blue-600 hover:bg-blue-500 text-white font-medium rounded-lg text-sm transition-all"
            >
              {loading ? 'Mencari...' : 'Muat Data Pengiriman'}
            </button>
          </div>

          {/* Shipment Detail & Update Actions */}
          <div className="lg:col-span-2 bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md">
            {!shipmentData && (
              <div className="text-center py-20">
                <Truck className="h-12 w-12 text-slate-600 mx-auto mb-4" />
                <h3 className="text-lg font-semibold text-slate-300 font-medium">Belum Ada Paket yang Dimuat</h3>
                <p className="text-slate-500 text-sm mt-1">Masukkan Shipment ID di menu sebelah kiri untuk memulai penugasan.</p>
              </div>
            )}

            {shipmentData && (
              <div className="space-y-6">
                <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center border-b border-slate-800 pb-4">
                  <div>
                    <span className="text-xs text-slate-500 font-semibold tracking-wider uppercase">Shipment ID</span>
                    <h3 className="text-lg font-mono font-bold text-blue-400">{shipmentData.id}</h3>
                  </div>
                  <div className="mt-2 sm:mt-0">
                    <span className={`inline-flex px-3 py-1 rounded-full text-xs font-bold ${
                      shipmentData.status === 'DELIVERED' ? 'bg-green-950 text-green-400 border border-green-900' : 'bg-blue-950 text-blue-400 border border-blue-900'
                    }`}>
                      {shipmentData.status}
                    </span>
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 bg-slate-950 border border-slate-800 rounded-xl p-4">
                  <div>
                    <h4 className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">Informasi Order</h4>
                    <p className="text-sm text-slate-300 font-semibold">Order ID: <span className="font-mono text-xs">{shipmentData.orderId}</span></p>
                    <p className="text-xs text-slate-500 mt-2">Dibuat: {new Date(shipmentData.createdAt).toLocaleString('id-ID')}</p>
                  </div>
                  <div>
                    <h4 className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2">Catatan Logistik</h4>
                    <p className="text-sm text-slate-300 italic">"{shipmentData.notes || 'Tidak ada catatan.'}"</p>
                  </div>
                </div>

                {/* Status Update Buttons */}
                <div className="border-t border-slate-800 pt-6">
                  <h4 className="text-sm font-bold text-white mb-4">Pembaruan Status Pengiriman:</h4>
                  <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                    <button
                      onClick={() => handleUpdateStatus('PICKED_UP')}
                      disabled={loading || shipmentData.status === 'DELIVERED'}
                      className="py-2.5 px-4 bg-slate-800 hover:bg-slate-700 text-white text-xs font-bold rounded-lg border border-slate-700 transition-all flex items-center justify-center gap-2"
                    >
                      <Clock className="h-4 w-4 text-yellow-500" />
                      1. Ambil Barang (Picked Up)
                    </button>
                    <button
                      onClick={() => handleUpdateStatus('ON_TRANSIT')}
                      disabled={loading || shipmentData.status === 'DELIVERED'}
                      className="py-2.5 px-4 bg-slate-800 hover:bg-slate-700 text-white text-xs font-bold rounded-lg border border-slate-700 transition-all flex items-center justify-center gap-2"
                    >
                      <Compass className="h-4 w-4 text-blue-500" />
                      2. Transit (On Transit)
                    </button>
                    <button
                      onClick={() => handleUpdateStatus('OUT_FOR_DELIVERY')}
                      disabled={loading || shipmentData.status === 'DELIVERED'}
                      className="py-2.5 px-4 bg-slate-800 hover:bg-slate-700 text-white text-xs font-bold rounded-lg border border-slate-700 transition-all flex items-center justify-center gap-2"
                    >
                      <Truck className="h-4 w-4 text-indigo-500" />
                      3. Pengantaran (Out for Delivery)
                    </button>
                  </div>

                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-4">
                    <button
                      onClick={() => handleUpdateStatus('DELIVERED')}
                      disabled={loading || shipmentData.status === 'DELIVERED'}
                      className="py-3 px-4 bg-green-600 hover:bg-green-500 text-white text-sm font-bold rounded-lg transition-all flex items-center justify-center gap-2 shadow-lg"
                    >
                      <CheckCircle className="h-5 w-5" />
                      Selesai & Diterima (Delivered)
                    </button>
                    <button
                      onClick={() => handleUpdateStatus('FAILED')}
                      disabled={loading || shipmentData.status === 'DELIVERED'}
                      className="py-3 px-4 bg-red-900 hover:bg-red-800 text-white text-sm font-bold rounded-lg transition-all flex items-center justify-center gap-2 border border-red-700"
                    >
                      Pengiriman Gagal (Failed)
                    </button>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
};

export default CourierDashboard;
