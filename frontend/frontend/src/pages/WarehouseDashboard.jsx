import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Package, MapPin, Compass, ArrowRightLeft, LogOut, CheckCircle, ArrowRight } from 'lucide-react';
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
      // 1. Record Movement in Warehouse Service
      const response = await api.warehouse.post('/warehouses/movements', {
        orderId: movementForm.orderId,
        warehouseId: selectedWarehouseId,
        type: movementForm.type,
        destinationCity: movementForm.destinationCity,
        notes: movementForm.notes,
      });

      const movementResult = response.data;
      setRoutingInstruction(movementResult);

      // 2. Fetch order info from Order Service to get trackingNumber
      const orderRes = await api.order.get(`/orders/${movementForm.orderId}`);
      const trackingNumber = orderRes.data.trackingNumber;

      // 3. Post a new status tracking event in Tracking Service
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
    <div className="min-h-screen bg-slate-950 flex flex-col md:flex-row">
      {/* Sidebar */}
      <aside className="w-full md:w-64 bg-slate-900 border-r border-slate-800 p-6 flex flex-col justify-between">
        <div>
          <div className="flex items-center gap-2 mb-8">
            <ArrowRightLeft className="h-6 w-6 text-blue-500" />
            <span className="font-extrabold text-xl bg-gradient-to-r from-blue-500 to-indigo-400 bg-clip-text text-transparent">
              Papiton Warehouse
            </span>
          </div>

          <div className="mb-6 p-4 bg-slate-950 rounded-xl border border-slate-800 space-y-2">
            <p className="text-xs text-slate-500 font-medium">Petugas Gudang:</p>
            <h4 className="font-semibold text-white truncate">{user.name}</h4>
            <span className="inline-block text-[10px] bg-blue-950 text-blue-400 border border-blue-900 px-2 py-0.5 rounded-full font-bold uppercase">
              {user.role}
            </span>
          </div>

          <div className="p-4 bg-slate-950 rounded-xl border border-slate-800 space-y-2 text-[11px] text-slate-500 leading-relaxed">
            <h5 className="font-bold text-slate-400 uppercase">Fungsi Halaman:</h5>
            <p>Digunakan untuk mensimulasikan scan masuk (*check-in*) dan keluar (*check-out*) barang dari titik-titik gudang transit.</p>
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
          <h1 className="text-2xl md:text-3xl font-extrabold text-white">Panel Logistik Gudang</h1>
          <p className="text-slate-400 text-sm">Rekam paket masuk/keluar, transit gudang, dan rute logistik ekspedisi.</p>
        </header>

        {error && <div className="mb-6 p-4 bg-red-950/50 border border-red-900 text-red-400 text-sm rounded-xl">{error}</div>}
        {success && <div className="mb-6 p-4 bg-green-950/50 border border-green-900 text-green-400 text-sm rounded-xl">{success}</div>}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Action Form */}
          <div className="lg:col-span-2 bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md">
            <form onSubmit={handleRecordMovement} className="space-y-4">
              <h3 className="text-lg font-bold text-white mb-2 flex items-center gap-2 border-b border-slate-800 pb-3">
                <Package className="h-5 w-5 text-blue-500" />
                Scan Aktivitas Paket
              </h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Pilih Gudang Aktif</label>
                  <select
                    value={selectedWarehouseId}
                    onChange={(e) => setSelectedWarehouseId(e.target.value)}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                  >
                    {warehouses.map(w => (
                      <option key={w.id} value={w.id}>{w.name} ({w.city})</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Aktivitas Scan</label>
                  <select
                    name="type"
                    value={movementForm.type}
                    onChange={handleFormChange}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                  >
                    <option value="WAREHOUSE_IN">Masuk Gudang (WAREHOUSE_IN)</option>
                    <option value="WAREHOUSE_OUT">Keluar Gudang (WAREHOUSE_OUT)</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Order ID (UUID)</label>
                <input
                  type="text"
                  name="orderId"
                  required
                  value={movementForm.orderId}
                  onChange={handleFormChange}
                  placeholder="Masukkan UUID Order..."
                  className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                />
              </div>

              {movementForm.type === 'WAREHOUSE_IN' && (
                <div>
                  <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Kota Tujuan Akhir (Di Label Paket)</label>
                  <select
                    name="destinationCity"
                    value={movementForm.destinationCity}
                    onChange={handleFormChange}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                  >
                    <option value="Jakarta">Jakarta</option>
                    <option value="Bandung">Bandung</option>
                    <option value="Surabaya">Surabaya</option>
                  </select>
                </div>
              )}

              <div>
                <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Catatan</label>
                <input
                  type="text"
                  name="notes"
                  value={movementForm.notes}
                  onChange={handleFormChange}
                  placeholder="Masukkan catatan tambahan..."
                  className="w-full bg-slate-950 border border-slate-800 rounded-lg py-2.5 px-4 text-sm text-white focus:outline-none focus:border-blue-500"
                />
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full py-3 bg-blue-600 hover:bg-blue-500 text-white font-medium rounded-lg text-sm transition-all"
              >
                {loading ? 'Menyimpan...' : 'Simpan Pergerakan Paket'}
              </button>
            </form>
          </div>

          {/* Routing Instruction Result */}
          <div className="lg:col-span-1 bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-md h-fit space-y-4">
            <h3 className="text-lg font-bold text-white mb-2 flex items-center gap-2">
              <Compass className="h-5 w-5 text-indigo-500" />
              Instruksi Rute Transit
            </h3>

            {!routingInstruction && (
              <div className="text-center py-10 text-slate-500 text-sm">
                Scan paket masuk untuk melihat rekomendasi rute selanjutnya.
              </div>
            )}

            {routingInstruction && (
              <div className="space-y-4">
                <div className="p-4 bg-slate-950 border border-slate-800 rounded-xl space-y-2">
                  <span className="text-[10px] bg-indigo-950 text-indigo-400 border border-indigo-900 px-2 py-0.5 rounded-full font-bold uppercase">
                    Hasil Analisis Rute
                  </span>
                  <p className="text-xs text-slate-400 font-semibold">{routingInstruction.message}</p>
                </div>

                {routingInstruction.nextWarehouseId && (
                  <div className="p-4 bg-slate-950 border border-slate-800 rounded-xl">
                    <span className="text-xs text-slate-500 font-semibold uppercase tracking-wider block mb-1">Kirim Ke Warehouse ID</span>
                    <p className="text-xs font-mono text-white bg-slate-900 p-2 rounded border border-slate-800 select-all font-bold">
                      {routingInstruction.nextWarehouseId}
                    </p>
                    <div className="mt-3 flex items-center gap-2 text-xs text-indigo-400 font-semibold">
                      <span>Proses Lanjutan:</span>
                      <span className="flex items-center gap-1">
                        Pindahkan Ke Rak transit
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
  );
};

export default WarehouseDashboard;
