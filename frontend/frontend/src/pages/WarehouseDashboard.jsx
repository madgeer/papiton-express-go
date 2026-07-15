import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Package, Compass, ArrowRightLeft, LogOut, ArrowRight } from 'lucide-react';
import { api } from '../services/api';

const WarehouseDashboard = () => {
  const navigate = useNavigate();
  const user = JSON.parse(localStorage.getItem('user') || '{}');

  // State
  const [warehouses, setWarehouses] = useState([]);
  const [selectedWarehouseId, setSelectedWarehouseId] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [routingInstruction, setRoutingInstruction] = useState(null);

  // Form State
  const [movementForm, setMovementForm] = useState({
    orderId: '',
    type: 'WAREHOUSE_IN', // 'WAREHOUSE_IN', 'WAREHOUSE_OUT'
    destinationCity: 'Bandung',
    notes: '',
  });

  const fetchWarehouses = async () => {
    try {
      const response = await api.warehouse.get('/warehouses');
      setWarehouses(response.data);
      if (response.data.length > 0) {
        setSelectedWarehouseId(response.data[0].id);
      }
    } catch (err) {
      setError('Gagal memuat daftar gudang');
    }
  };

  useEffect(() => {
    fetchWarehouses();
  }, []);

  const handleLogout = () => {
    localStorage.clear();
    navigate('/login');
  };

  const handleFormChange = (e) => {
    setMovementForm({ ...movementForm, [e.target.name]: e.target.value });
  };

  const handleRecordMovement = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');
    setRoutingInstruction(null);

    const activeWarehouse = warehouses.find(w => w.id === selectedWarehouseId);
    if (!activeWarehouse) return;

    try {
      const response = await api.warehouse.post('/warehouses/movements', {
        orderId: movementForm.orderId,
        warehouseId: selectedWarehouseId,
        type: movementForm.type,
        destinationCity: movementForm.destinationCity,
        notes: movementForm.notes,
      });

      const movementResult = response.data;
      setRoutingInstruction(movementResult);

      const orderRes = await api.order.get(`/orders/${movementForm.orderId}`);
      const trackingNumber = orderRes.data.trackingNumber;

      const trackingEventStatus = movementForm.type === 'WAREHOUSE_IN' ? 'ARRIVED_AT_WAREHOUSE' : 'DEPARTED_FROM_WAREHOUSE';
      const trackingEventDesc = movementForm.type === 'WAREHOUSE_IN' 
        ? `Paket telah masuk dan diproses di gudang ${activeWarehouse.name} (${activeWarehouse.city}). ${movementResult.message}`
        : `Paket telah meninggalkan gudang ${activeWarehouse.name} (${activeWarehouse.city}).`;

      await api.tracking.post(`/trackings/${trackingNumber}/events`, {
        status: trackingEventStatus,
        location: `${activeWarehouse.name} (${activeWarehouse.city})`,
        description: trackingEventDesc,
      });

      setSuccess(`Aktivitas paket masuk/keluar gudang berhasil direkam!`);
    } catch (err) {
      setError(err.response?.data?.error || 'Gagal merekam data pergerakan paket. Pastikan Order ID valid.');
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
            <span className="text-[#FF6600] bg-white px-2 ml-1 rounded-sm">Warehouse</span>
          </h1>
        </div>
        <div className="flex items-center gap-4 text-xs text-white">
          <span className="font-mono">Logistics & Hub Operations</span>
          <button
            onClick={handleLogout}
            className="bg-[#FF6600] hover:bg-[#E05300] text-white px-3 py-1 font-bold rounded-sm border border-[#B34700]"
          >
            Log Out
          </button>
        </div>
      </header>

      {/* Main Container */}
      <div className="flex-1 flex flex-col md:flex-row">
        {/* Sidebar */}
        <aside className="w-full md:w-64 bg-[#EAEAEA] border-r border-[#CCCCCC] p-4 flex flex-col justify-between">
          <div className="space-y-6">
            <div className="bg-white border border-[#CCCCCC] p-3 rounded-sm">
              <span className="text-[10px] text-slate-500 font-bold block uppercase mb-1">Hub Operator</span>
              <h4 className="font-bold text-[#4D148C] text-sm truncate">{user.name}</h4>
              <span className="inline-block mt-2 text-[9px] bg-slate-200 text-[#4D148C] font-bold px-1.5 py-0.5 rounded-sm border border-slate-300">
                {user.role}
              </span>
            </div>

            <div className="p-3 bg-white border border-[#CCCCCC] rounded-sm text-xs text-slate-500 space-y-1.5">
              <h5 className="font-bold text-[#333333] uppercase">Deskripsi Tugas:</h5>
              <p className="leading-relaxed">
                Gunakan panel ini untuk mencatat paket transit masuk (*Check-in*) atau transit keluar (*Check-out*) dari gudang lokasi Anda saat ini.
              </p>
            </div>
          </div>
        </aside>

        {/* Content Panel */}
        <main className="flex-1 p-6 space-y-6">
          <div className="border-b border-[#CCCCCC] pb-3">
            <h2 className="text-xl font-black text-[#4D148C] tracking-tight">PENCATATAN TRANSIT & BARCODE SCAN SIMULATION</h2>
            <p className="text-xs text-slate-500 font-sans mt-0.5">Pantau status pergudangan dan routing rute logistik ekspedisi.</p>
          </div>

          {error && <div className="p-2.5 bg-[#FFF2F2] border border-[#FF9999] text-[#990000] text-xs font-semibold rounded-sm">{error}</div>}
          {success && <div className="p-2.5 bg-[#F2FFF2] border border-[#99FF99] text-[#006600] text-xs font-semibold rounded-sm">{success}</div>}

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Record Scan Form */}
            <div className="lg:col-span-2 bg-white border border-[#CCCCCC] p-6 rounded-sm shadow-sm">
              <form onSubmit={handleRecordMovement} className="space-y-4">
                <h3 className="text-xs font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 flex items-center gap-1.5 uppercase">
                  <ArrowRightLeft className="h-4 w-4 text-[#FF6600]" />
                  Form Scan Pergerakan Barang
                </h3>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Gudang Lokasi Anda</label>
                    <select
                      value={selectedWarehouseId}
                      onChange={(e) => setSelectedWarehouseId(e.target.value)}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    >
                      {warehouses.map(w => (
                        <option key={w.id} value={w.id}>{w.name} ({w.city})</option>
                      ))}
                    </select>
                  </div>
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Jenis Aktivitas</label>
                    <select
                      name="type"
                      value={movementForm.type}
                      onChange={handleFormChange}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    >
                      <option value="WAREHOUSE_IN">Masuk Gudang (WAREHOUSE_IN)</option>
                      <option value="WAREHOUSE_OUT">Keluar Gudang (WAREHOUSE_OUT)</option>
                    </select>
                  </div>
                </div>

                <div>
                  <label className="block text-[11px] font-bold text-[#333333] mb-1">Order ID (UUID)</label>
                  <input
                    type="text"
                    name="orderId"
                    required
                    value={movementForm.orderId}
                    onChange={handleFormChange}
                    placeholder="Masukkan UUID Order..."
                    className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none focus:border-[#4D148C]"
                  />
                </div>

                {movementForm.type === 'WAREHOUSE_IN' && (
                  <div>
                    <label className="block text-[11px] font-bold text-[#333333] mb-1">Kota Tujuan Penerima (Sesuai Label)</label>
                    <select
                      name="destinationCity"
                      value={movementForm.destinationCity}
                      onChange={handleFormChange}
                      className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                    >
                      <option value="Jakarta">Jakarta</option>
                      <option value="Bandung">Bandung</option>
                      <option value="Surabaya">Surabaya</option>
                    </select>
                  </div>
                )}

                <div>
                  <label className="block text-[11px] font-bold text-[#333333] mb-1">Catatan Gerakan Paket</label>
                  <input
                    type="text"
                    name="notes"
                    value={movementForm.notes}
                    onChange={handleFormChange}
                    placeholder="Masukkan catatan logistik..."
                    className="w-full bg-[#FFFFFF] border border-[#999999] rounded-sm py-1.5 px-3 text-xs focus:outline-none"
                  />
                </div>

                <button
                  type="submit"
                  disabled={loading}
                  className="w-full py-2.5 bg-[#4D148C] hover:bg-[#390F66] text-white text-xs font-bold rounded-sm border-b-2 border-[#330D5C]"
                >
                  {loading ? 'Menyimpan...' : 'SIMPAN DATA TRANSIT'}
                </button>
              </form>
            </div>

            {/* Routing Output */}
            <div className="lg:col-span-1 bg-white border border-[#CCCCCC] p-4 rounded-sm shadow-sm h-fit space-y-4">
              <h3 className="text-xs font-bold text-[#4D148C] border-b border-[#CCCCCC] pb-2 flex items-center gap-1.5 uppercase">
                <Compass className="h-4 w-4 text-[#FF6600]" />
                Petunjuk Rute Selanjutnya
              </h3>

              {!routingInstruction && (
                <div className="text-center py-8 text-slate-500 text-xs">
                  Lakukan scan masuk untuk mendapatkan petunjuk transit rute logistik selanjutnya dari server.
                </div>
              )}

              {routingInstruction && (
                <div className="space-y-4">
                  <div className="p-3 bg-[#F9F9F9] border border-[#CCCCCC] rounded-sm space-y-1.5">
                    <span className="inline-block text-[9px] font-bold px-1.5 py-0.5 rounded-sm bg-[#E9E1F5] text-[#4D148C] border border-[#D5C2EB]">
                      Status Analisis Rute
                    </span>
                    <p className="text-xs text-slate-700 font-semibold">{routingInstruction.message}</p>
                  </div>

                  {routingInstruction.nextWarehouseId && (
                    <div className="p-3 bg-[#F9F9F9] border border-[#CCCCCC] rounded-sm space-y-2">
                      <span className="text-[10px] text-slate-500 font-bold block uppercase">Kirim Paket Ke Warehouse ID</span>
                      <p className="text-xs font-mono font-bold text-slate-800 bg-white p-2 border border-[#CCCCCC] select-all">
                        {routingInstruction.nextWarehouseId}
                      </p>
                      <div className="text-xs text-[#FF6600] font-bold flex items-center gap-1">
                        <span>Aksi Gudang:</span>
                        <span className="flex items-center gap-1">
                          Masukkan ke rak transit
                          <ArrowRight className="h-3.5 w-3.5" />
                        </span>
                      </div>
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        </main>
      </div>

      {/* Footer */}
      <footer className="bg-[#EAEAEA] border-t border-[#CCCCCC] py-3 text-center text-xs text-[#666666] font-sans">
        © 2003 Papiton Express Inc. All rights reserved. Warehouse & Hub System.
      </footer>
    </div>
  );
};

export default WarehouseDashboard;
