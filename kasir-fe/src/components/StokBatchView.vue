<script setup>
import { ref, onMounted, computed } from 'vue'
import { reportApi } from '../services/api'
import { Layers, Calendar, AlertTriangle, CheckCircle, Search, Building2, User } from 'lucide-vue-next'

const loading = ref(false)
const batches = ref([])
const searchQuery = ref('')
const filterNearED = ref(false)

const fetchBatches = async () => {
  loading.value = true
  try {
    const res = await reportApi.getStokPerBarang()
    batches.value = res.data.data || []
  } catch (err) {
    console.error('Failed to fetch batch reports', err)
  } finally {
    loading.value = false
  }
}

const filteredBatches = computed(() => {
  return batches.value.filter(b => {
    const matchesSearch = !searchQuery.value ||
      b.nama_barang.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      b.kode_barang.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      b.nama_supplier.toLowerCase().includes(searchQuery.value.toLowerCase())

    const matchesNearED = !filterNearED.value || b.is_near_ed

    return matchesSearch && matchesNearED
  })
})

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

onMounted(() => {
  fetchBatches()
})

defineExpose({
  refresh: fetchBatches
})
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
      <div>
        <h2 class="text-base font-semibold text-slate-900">Laporan Stok per Batch (Rotasi FEFO)</h2>
        <p class="text-xs text-slate-500">Monitoring tanggal kadaluarsa (First Expired First Out), Supplier Asal, dan konversi Box / Eceran</p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="filterNearED = !filterNearED"
          :class="[
            'inline-flex items-center gap-2 px-3.5 py-2 text-xs font-medium rounded-lg border transition',
            filterNearED
              ? 'bg-rose-50 border-rose-300 text-rose-700 shadow-sm font-semibold'
              : 'bg-white border-slate-200 text-slate-600 hover:bg-slate-50'
          ]"
        >
          <AlertTriangle class="w-4 h-4 text-rose-600" />
          <span>Hanya Mendekati ED (&lt; 90 Hari)</span>
        </button>
      </div>
    </div>

    <!-- Search & Filter Bar -->
    <div class="flex flex-col sm:flex-row items-center gap-3 bg-white p-3 rounded-xl border border-slate-200 shadow-sm">
      <div class="relative w-full sm:w-80">
        <Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
        <input
          type="text"
          v-model="searchQuery"
          placeholder="Cari barang atau supplier..."
          class="w-full pl-9 pr-3 py-1.5 text-xs rounded-lg border border-slate-300 text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-1 focus:ring-emerald-600"
        />
      </div>

      <div class="ml-auto text-xs text-slate-500 font-medium">
        Menampilkan {{ filteredBatches.length }} dari {{ batches.length }} batch aktif
      </div>
    </div>

    <!-- Table Report -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold">
              <th class="py-3 px-4">Batch ID</th>
              <th class="py-3 px-4">Nama Produk</th>
              <th class="py-3 px-4">Kategori</th>
              <th class="py-3 px-4">Supplier Asal</th>
              <th class="py-3 px-4">Expired Date (FEFO)</th>
              <th class="py-3 px-4 text-center">Total Qty (Pcs)</th>
              <th class="py-3 px-4 text-center">Konversi Satuan</th>
              <th class="py-3 px-4 text-right">Status ED</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="8" class="py-12 text-center text-slate-400">Memuat data batch stok...</td>
            </tr>
            <tr v-else-if="filteredBatches.length === 0">
              <td colspan="8" class="py-12 text-center text-slate-400">Tidak ada batch yang sesuai filter.</td>
            </tr>
            <tr v-else v-for="b in filteredBatches" :key="b.batch_id" class="hover:bg-slate-50/70 transition">
              <td class="py-3 px-4 font-mono font-medium text-slate-500">#BATCH-{{ b.batch_id }}</td>
              <td class="py-3 px-4">
                <div class="font-semibold text-slate-900">{{ b.nama_barang }}</div>
                <div class="text-[10px] text-slate-400">{{ b.kode_barang }}</div>
              </td>
              <td class="py-3 px-4">
                <span class="inline-block px-2 py-0.5 rounded bg-slate-100 text-slate-700 text-[11px] font-medium">
                  {{ b.kategori }}
                </span>
              </td>
              <td class="py-3 px-4">
                <div class="flex items-center gap-1.5">
                  <Building2 v-if="b.jenis_entitas === 'BADAN_USAHA'" class="w-3.5 h-3.5 text-blue-600" />
                  <User v-else class="w-3.5 h-3.5 text-emerald-600" />
                  <span class="text-slate-800 font-medium">{{ b.nama_supplier }}</span>
                </div>
              </td>
              <td class="py-3 px-4 font-medium">
                <div class="flex items-center gap-1.5">
                  <Calendar class="w-3.5 h-3.5 text-slate-400" />
                  <span :class="b.is_near_ed ? 'text-rose-600 font-bold' : 'text-slate-800'">
                    {{ formatDate(b.expired_date) }}
                  </span>
                </div>
                <div v-if="b.days_until_ed > 0" class="text-[10px] text-slate-400">
                  ({{ b.days_until_ed }} hari lagi)
                </div>
                <div v-else class="text-[10px] text-rose-600 font-bold">
                  (SUDAH KADALUARSA)
                </div>
              </td>
              <td class="py-3 px-4 text-center font-bold text-slate-900 text-sm">
                {{ b.total_qty_pcs }}
              </td>
              <td class="py-3 px-4 text-center">
                <span class="inline-block px-2.5 py-1 bg-emerald-50 text-emerald-800 rounded font-medium text-[11px]">
                  {{ b.box }} Box + {{ b.pcs }} Pcs
                </span>
                <div class="text-[10px] text-slate-400 mt-0.5">@{{ b.pcs_per_box }} pcs/box</div>
              </td>
              <td class="py-3 px-4 text-right">
                <span
                  v-if="b.is_near_ed"
                  class="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] font-bold bg-rose-100 text-rose-700"
                >
                  <AlertTriangle class="w-3 h-3" />
                  <span>PRIORITAS JUAL (&lt; 90 Hari)</span>
                </span>
                <span
                  v-else
                  class="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] font-medium bg-slate-100 text-slate-600"
                >
                  <CheckCircle class="w-3 h-3 text-emerald-600" />
                  <span>AMAN</span>
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
