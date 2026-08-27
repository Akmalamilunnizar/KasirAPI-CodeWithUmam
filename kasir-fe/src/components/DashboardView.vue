<script setup>
import { ref, onMounted } from 'vue'
import { reportApi, barangApi } from '../services/api'
import {
  TrendingUp,
  AlertTriangle,
  Package,
  Layers,
  CalendarClock,
  ArrowRight,
  PlusCircle,
  CheckCircle,
  AlertCircle
} from 'lucide-vue-next'

const emit = defineEmits(['openMutasi', 'navigate'])

const loading = ref(true)
const topBarangs = ref([])
const lowStocks = ref([])
const stokBatches = ref([])
const totalBarangsCount = ref(0)
const nearEDCount = ref(0)

const fetchDashboardData = async () => {
  loading.value = true
  try {
    const [resTop, resAlert, resBatch, resBarang] = await Promise.all([
      reportApi.getTopBarang(),
      reportApi.getLowStockAlerts(),
      reportApi.getStokPerBarang(),
      barangApi.getAll()
    ])

    topBarangs.value = resTop.data.data || []
    lowStocks.value = resAlert.data.data || []
    stokBatches.value = resBatch.data.data || []
    totalBarangsCount.value = resBarang.data.total || 0

    nearEDCount.value = stokBatches.value.filter(b => b.is_near_ed).length
  } catch (err) {
    console.error('Failed to load dashboard', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDashboardData()
})

defineExpose({
  refresh: fetchDashboardData
})
</script>

<template>
  <div class="space-y-6">
    <!-- Quick Actions Bar -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
      <div>
        <h2 class="text-lg font-semibold text-slate-900">Dashboard Overview</h2>
        <p class="text-xs text-slate-500">Ringkasan stok pertanian, rotasi batch FEFO, dan status penjualan</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="$emit('openMutasi')"
          class="inline-flex items-center gap-2 px-3.5 py-2 text-xs font-medium text-white bg-emerald-600 hover:bg-emerald-700 rounded-lg shadow-sm transition"
        >
          <PlusCircle class="w-4 h-4" />
          <span>Catat Mutasi Stok</span>
        </button>
      </div>
    </div>

    <!-- 4 Stats Cards Grid -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-500">Total Master Produk</span>
          <div class="p-2 bg-slate-100 rounded-lg text-slate-600">
            <Package class="w-4 h-4" />
          </div>
        </div>
        <div class="mt-2 text-2xl font-bold text-slate-900">{{ totalBarangsCount }}</div>
        <div class="mt-1 text-[11px] text-slate-500">Katalog barang aktif di toko</div>
      </div>

      <div class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-500">Total Batch Aktif</span>
          <div class="p-2 bg-emerald-50 rounded-lg text-emerald-600">
            <Layers class="w-4 h-4" />
          </div>
        </div>
        <div class="mt-2 text-2xl font-bold text-emerald-700">{{ stokBatches.length }}</div>
        <div class="mt-1 text-[11px] text-slate-500">Batch tersimpan dalam FEFO</div>
      </div>

      <div class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-500">Peringatan Stok Menipis</span>
          <div class="p-2 bg-amber-50 rounded-lg text-amber-600">
            <AlertTriangle class="w-4 h-4" />
          </div>
        </div>
        <div class="mt-2 text-2xl font-bold text-amber-600">{{ lowStocks.length }}</div>
        <div class="mt-1 text-[11px] text-slate-500">Produk $\le$ ambang batas min stok</div>
      </div>

      <div class="bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-500">Mendekati ED (&lt; 90 Hari)</span>
          <div class="p-2 bg-rose-50 rounded-lg text-rose-600">
            <CalendarClock class="w-4 h-4" />
          </div>
        </div>
        <div class="mt-2 text-2xl font-bold text-rose-600">{{ nearEDCount }}</div>
        <div class="mt-1 text-[11px] text-slate-500">Batch butuh prioritas penjualan</div>
      </div>
    </div>

    <!-- Main Grid: Top 3 Selling Items & Low Stock Alerts -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- Top 3 Best Selling Agricultural Items (5 Cols) -->
      <div class="lg:col-span-5 bg-white p-5 rounded-xl border border-slate-200 shadow-sm flex flex-col justify-between">
        <div>
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div class="flex items-center gap-2">
              <TrendingUp class="w-4 h-4 text-emerald-600" />
              <h3 class="text-sm font-semibold text-slate-900">Top 3 Produk Terlaris</h3>
            </div>
            <span class="text-[11px] px-2 py-0.5 bg-emerald-50 text-emerald-700 font-medium rounded-full">Favorit Petani</span>
          </div>

          <div v-if="loading" class="py-12 text-center text-xs text-slate-400">Memuat data...</div>
          <div v-else-if="topBarangs.length === 0" class="py-12 text-center text-xs text-slate-400">Belum ada data transaksi penjualan</div>
          <div v-else class="mt-4 space-y-4">
            <div
              v-for="(item, idx) in topBarangs"
              :key="item.barang_id"
              class="p-3.5 rounded-lg border border-slate-100 bg-slate-50/60 hover:bg-slate-50 transition"
            >
              <div class="flex items-center justify-between mb-1.5">
                <div class="flex items-center gap-2.5">
                  <span
                    :class="[
                      'w-5 h-5 rounded flex items-center justify-center text-[11px] font-bold',
                      idx === 0 ? 'bg-amber-400 text-amber-950' : idx === 1 ? 'bg-slate-300 text-slate-900' : 'bg-amber-700/30 text-amber-900'
                    ]"
                  >
                    {{ idx + 1 }}
                  </span>
                  <div>
                    <h4 class="text-xs font-semibold text-slate-900 leading-tight">{{ item.nama_barang }}</h4>
                    <span class="text-[10px] text-slate-500">{{ item.kode_barang }} • {{ item.kategori }}</span>
                  </div>
                </div>
                <div class="text-right">
                  <span class="text-xs font-bold text-emerald-700">{{ item.total_terjual }}</span>
                  <span class="text-[10px] text-slate-500 ml-1">pcs</span>
                </div>
              </div>

              <!-- Progress bar representation -->
              <div class="w-full bg-slate-200 h-1.5 rounded-full overflow-hidden">
                <div
                  class="bg-emerald-600 h-full rounded-full transition-all duration-500"
                  :style="{ width: `${Math.min(100, (item.total_terjual / (topBarangs[0]?.total_terjual || 1)) * 100)}%` }"
                ></div>
              </div>
            </div>
          </div>
        </div>

        <button
          @click="$emit('navigate', 'tracking')"
          class="mt-4 pt-3 border-t border-slate-100 flex items-center justify-center gap-1.5 text-xs font-medium text-emerald-700 hover:text-emerald-800"
        >
          <span>Lihat Seluruh Riwayat Mutasi</span>
          <ArrowRight class="w-3.5 h-3.5" />
        </button>
      </div>

      <!-- Low Stock Alerts Table (7 Cols) -->
      <div class="lg:col-span-7 bg-white p-5 rounded-xl border border-slate-200 shadow-sm flex flex-col justify-between">
        <div>
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div class="flex items-center gap-2">
              <AlertCircle class="w-4 h-4 text-amber-600" />
              <h3 class="text-sm font-semibold text-slate-900">Peringatan Ambang Batas Stok (Low Stock)</h3>
            </div>
            <span class="text-[11px] text-slate-500">{{ lowStocks.length }} item butuh restock</span>
          </div>

          <div v-if="loading" class="py-12 text-center text-xs text-slate-400">Memuat data...</div>
          <div v-else-if="lowStocks.length === 0" class="py-12 text-center text-xs text-emerald-600 flex flex-col items-center gap-1">
            <CheckCircle class="w-6 h-6" />
            <span>Semua stok barang berada di atas ambang minimum aman.</span>
          </div>
          <div v-else class="mt-3 overflow-x-auto">
            <table class="w-full text-left text-xs">
              <thead>
                <tr class="text-[11px] font-semibold text-slate-500 border-b border-slate-100 pb-2">
                  <th class="py-2">Produk</th>
                  <th class="py-2">Kategori</th>
                  <th class="py-2 text-center">Stok Saat Ini</th>
                  <th class="py-2 text-center">Min. Stok</th>
                  <th class="py-2 text-right">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr v-for="alert in lowStocks" :key="alert.barang_id" class="hover:bg-slate-50/60">
                  <td class="py-2.5 pr-2">
                    <div class="font-medium text-slate-900">{{ alert.nama_barang }}</div>
                    <div class="text-[10px] text-slate-400">{{ alert.kode_barang }}</div>
                  </td>
                  <td class="py-2.5 text-slate-600 text-[11px]">{{ alert.kategori }}</td>
                  <td class="py-2.5 text-center font-semibold text-slate-900">
                    <span :class="alert.current_stok === 0 ? 'text-rose-600' : 'text-amber-600'">
                      {{ alert.current_stok }} pcs
                    </span>
                  </td>
                  <td class="py-2.5 text-center text-slate-500">{{ alert.min_stok }} pcs</td>
                  <td class="py-2.5 text-right">
                    <span
                      :class="[
                        'inline-block px-2 py-0.5 rounded text-[10px] font-medium',
                        alert.status === 'OUT_OF_STOCK'
                          ? 'bg-rose-100 text-rose-700'
                          : alert.status === 'CRITICAL'
                          ? 'bg-amber-100 text-amber-800'
                          : 'bg-yellow-50 text-yellow-700 border border-yellow-200'
                      ]"
                    >
                      {{ alert.status }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <button
          @click="$emit('navigate', 'barang')"
          class="mt-4 pt-3 border-t border-slate-100 flex items-center justify-center gap-1.5 text-xs font-medium text-slate-600 hover:text-slate-900"
        >
          <span>Kelola Master Barang</span>
          <ArrowRight class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>
  </div>
</template>
