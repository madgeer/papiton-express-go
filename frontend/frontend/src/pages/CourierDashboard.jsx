import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Truck, Compass, CheckCircle, Clock, LogOut } from 'lucide-react';
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

  const handleFetchShipment = async () => {
    if (!shipmentIdInput) return;
    setLoading(true);
    setError('');
    setSuccess('');
    try {
      const response = await api.shipping.get(`/shippings/${shipmentIdInput}`);
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

      const response = await api.shipping.put(`/shippings/${shipmentData.id}/status`, {
        status: newStatus,
        notes: notesMap[newStatus] || 'Status updated by courier.',
      });

      const orderRes = await api.order.get(`/orders/${shipmentData.orderId}`);
      const trackingNum = orderRes.data.trackingNumber;

      await api.tracking.post(`/trackings/${trackingNum}/events`, {
        status: newStatus,
        location: courierInfo ? `Kurir ${courierInfo.name}` : 'Kurir Logistik',
        description: notesMap[newStatus],
      });
      
      setSuccess(`Status paket & pelacakan berhasil diperbarui menjadi: ${newStatus}`);
      setShipmentData(response.data);
      fetchCourierInfo();
    } catch (err) {
      setError('Gagal memperbarui status pengiriman');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#F0F0F0] flex flex-col font-sans text-[#333333]">
      {/* Brand Header */}
      <header className="bg-[#4D148C] py-4 px-6 border-b-4 border-[#FF6600] flex justify-between items-center shadow-md">
        <div className="flex items-center gap-2">
          <h1 className="text-3xl font-black tracking-tighter">
            <span className="text-white">Papiton</span>
            <span className="text-[#FF6600] bg-white px-2 ml-1 rounded-sm">Courier</span>
          </h1>
        </div>
        <div className="flex items-center gap-4 text-xs text-white">
          <span className="font-mono">Courier Dispatch Portal</span>
          <button
            onClick={handleLogout}
            className="bg-[#FF6600] hover:bg-[#E05300] text-white px-3 py-1 font-bold rounded-sm border border-[#B34700]"
          >
            Log Out
          </button>
        </div>
      </header>

      {/* Main Area */}
      <div className="flex-1 flex flex-col md:flex-row">
        {/* Sidebar */}
        <aside className="w-full md:w-64 bg-[#EAEAEA] border-r border-[#CCCCCC] p-4 flex flex-col justify-between">
          <div className="space-y-6">
            <div className="bg-white border border-[#CCCCCC] p-3 rounded-sm">
              <span className="text-[10px] text-slate-500 font-bold block uppercase mb-1">Active Courier</span>
              <h4 className="font-bold text-[#4D148C] text-sm truncate">{user.name}</h4>
              <p className="text-xs text-slate-600 mt-1">Status Penugasan:</p>
              {courierInfo ? (
                <span className={`inline-block mt-1 text-[10px] px-2 py-0.5 rounded-sm font-bold uppercase ${
                  courierInfo.status === 'AVAILABLE' ? 'bg-[#F2FFF2] text-[#006600] border border-[#99FF99]' : 'bg-[#FFF2F2] text-[#990000] border border-[#FF9999]'
                }`}>
                  {courierInfo.status}
                </span>
              ) : (
                <span className="text-xs text-slate-500">Loading status...</span>
              )}
            </div>

            <div className="p-3 bg-white border border-[#CCCCCC] rounded-sm text-xs text-slate-500 space-y-1.5">
              <h5 className="font-bold text-[#333333] uppercase">Petunjuk Kurir:</h5>
              <p className="leading-relaxed">
                Salin <strong>Shipment ID</strong> dari panel admin, paste di menu pencarian untuk memuat dan memproses tugas Anda.
              </p>
            </div>
          </div>
        </aside>

        {/* Content Workspace */}
        <main className="flex-1 p-6 space-y-6">
          <div className="border-b border-[#CCCCCC] pb-3">
            <h2 className="text-xl font-black text-[#4D148C] tracking-tight">KONTROL DISPATCH & STATUS KIRIMAN</h2>
            <p className="text-xs text-slate-500 font-sans mt-0.5">Pantau, ambil, dan selesaikan pengantaran barang logistik.</p>
          </div>

          {error && <div className="p-2.5 bg-[#FFF2F2] border border-[#FF9999] text-[#990000] text-xs font-semibold rounded-sm">{error}</div>}
          {success && <div className="p-2.5 bg-[#F2FFF2] border border-[#99FF99] text-[#006600] text-xs font-semibold rounded-sm">{success}</div>}

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Search Input Box */}
            <div className="lg:col-span-1 bg-white border border-[#CCCCCC] p-4 rounded-sm shadow-sm space-y-4 h-fit">
              <h3 className="text-xs font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 flex items-center gap-1.5">
                <Compass className="h-4 w-4 text-[#FF6600]" />
                MUAT SHIPMENT BARU
              </h3>
              <div>
                <label className="block text-[10px] font-bold text-slate-500 mb-1">Shipment ID (UUID)</label>
                <input
                  type="text"
                  value={shipmentIdInput}
                  onChange={(e) => setShipmentIdInput(e.target.value)}
                  placeholder="Paste Shipment ID..."
                  className="w-full bg-white border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none focus:border-[#4D148C]"
                />
              </div>
              <button
                onClick={handleFetchShipment}
                disabled={loading}
                className="w-full py-2 bg-[#4D148C] hover:bg-[#390F66] text-white text-xs font-bold rounded-sm border-b-2 border-[#330D5C]"
              >
                {loading ? 'Memuat...' : 'CARI PENGIRIMAN'}
              </button>
            </div>

            {/* Shipment Workspace details */}
            <div className="lg:col-span-2 bg-white border border-[#CCCCCC] p-6 rounded-sm shadow-sm">
              {!shipmentData && (
                <div className="text-center py-16">
                  <Truck className="h-10 w-10 text-slate-400 mx-auto mb-2" />
                  <p className="text-xs text-slate-500 font-bold">Belum Ada Shipment yang Dimuat</p>
                  <p className="text-[10px] text-slate-400 mt-1">Masukkan Shipment ID di panel kiri untuk memproses tugas kurir Anda.</p>
                </div>
              )}

              {shipmentData && (
                <div className="space-y-6">
                  <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center border-b border-[#CCCCCC] pb-3">
                    <div>
                      <span className="text-[10px] text-slate-500 font-bold block uppercase">Shipment ID</span>
                      <h3 className="text-sm font-mono font-bold text-blue-700">{shipmentData.id}</h3>
                    </div>
                    <div className="mt-2 sm:mt-0">
                      <span className={`inline-block px-3 py-0.5 rounded-sm text-xs font-bold ${
                        shipmentData.status === 'DELIVERED' ? 'bg-[#F2FFF2] text-[#006600] border border-[#99FF99]' : 'bg-[#E9E1F5] text-[#4D148C] border border-[#D5C2EB]'
                      }`}>
                        {shipmentData.status}
                      </span>
                    </div>
                  </div>

                  <div className="bg-[#F9F9F9] border border-[#CCCCCC] p-4 rounded-sm grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <span className="text-[10px] text-slate-500 font-bold block uppercase">Order ID</span>
                      <p className="text-xs font-mono font-bold text-slate-700">{shipmentData.orderId}</p>
                      <p className="text-[10px] text-slate-400 mt-2">Dibuat: {new Date(shipmentData.createdAt).toLocaleString('id-ID')}</p>
                    </div>
                    <div>
                      <span className="text-[10px] text-slate-500 font-bold block uppercase">Logistics Notes</span>
                      <p className="text-xs text-slate-600 italic">"{shipmentData.notes || 'Tidak ada catatan logistik.'}"</p>
                    </div>
                  </div>

                  {/* Actions buttons */}
                  <div className="border-t border-[#CCCCCC] pt-4 space-y-4">
                    <h4 className="text-xs font-bold text-[#4D148C] uppercase tracking-wider">Langkah Proses Pembaruan Status:</h4>
                    
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                      <button
                        onClick={() => handleUpdateStatus('PICKED_UP')}
                        disabled={loading || shipmentData.status === 'DELIVERED'}
                        className="py-2 px-3 bg-[#EAEAEA] hover:bg-[#DFDFDF] border border-[#CCCCCC] text-[#333333] text-xs font-bold rounded-sm transition-all flex items-center justify-center gap-1"
                      >
                        <Clock className="h-3.5 w-3.5 text-[#FF6600]" />
                        1. Picked Up
                      </button>
                      <button
                        onClick={() => handleUpdateStatus('ON_TRANSIT')}
                        disabled={loading || shipmentData.status === 'DELIVERED'}
                        className="py-2 px-3 bg-[#EAEAEA] hover:bg-[#DFDFDF] border border-[#CCCCCC] text-[#333333] text-xs font-bold rounded-sm transition-all flex items-center justify-center gap-1"
                      >
                        <Compass className="h-3.5 w-3.5 text-blue-600" />
                        2. On Transit
                      </button>
                      <button
                        onClick={() => handleUpdateStatus('OUT_FOR_DELIVERY')}
                        disabled={loading || shipmentData.status === 'DELIVERED'}
                        className="py-2 px-3 bg-[#EAEAEA] hover:bg-[#DFDFDF] border border-[#CCCCCC] text-[#333333] text-xs font-bold rounded-sm transition-all flex items-center justify-center gap-1"
                      >
                        <Truck className="h-3.5 w-3.5 text-indigo-700" />
                        3. Out for Delivery
                      </button>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
                      <button
                        onClick={() => handleUpdateStatus('DELIVERED')}
                        disabled={loading || shipmentData.status === 'DELIVERED'}
                        className="py-2.5 px-4 bg-[#008000] hover:bg-[#006600] text-white text-xs font-bold rounded-sm border-b-2 border-[#004d00] transition-all flex items-center justify-center gap-1.5 shadow-sm"
                      >
                        <CheckCircle className="h-4 w-4" />
                        KONFIRMASI SELESAI (DELIVERED)
                      </button>
                      <button
                        onClick={() => handleUpdateStatus('FAILED')}
                        disabled={loading || shipmentData.status === 'DELIVERED'}
                        className="py-2.5 px-4 bg-[#990000] hover:bg-[#800000] text-white text-xs font-bold rounded-sm border-b-2 border-[#660000] transition-all"
                      >
                        PENGIRIMAN GAGAL (FAILED)
                      </button>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </main>
      </div>

      {/* Footer */}
      <footer className="bg-[#EAEAEA] border-t border-[#CCCCCC] py-3 text-center text-xs text-[#666666] font-sans">
        © 2003 Papiton Express Inc. All rights reserved. Courier Dispatch System.
      </footer>
    </div>
  );
};

export default CourierDashboard;
