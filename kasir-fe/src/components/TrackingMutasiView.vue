<script setup>
import { ref, onMounted, computed } from 'vue'
import { reportApi } from '../services/api'
import {
  ArrowDownRight,
  ArrowUpRight,
  History,
  Calendar,
  Building2,
  User,
  Search
} from 'lucide-vue-next'

const loading = ref(false)
const logs = ref([])
const filterJenis = ref('')
const searchQuery = ref('')

const fetchLogs = async () => {
  loading.value = true
  try {
    const params = {}
    if (filterJenis.value) params.jenis_transaksi = filterJenis.value

    const res = await reportApi.getTrackingStok(params)
    logs.value = res.data.data || []
  } catch (err) {
    console.error('Failed to fetch tracking logs', err)
  } finally {
    loading.value = false
  }
}

const filteredLogs = computed(() => {
  return logs.value.filter(item => {
    if (!searchQuery.value) return true
    const q = searchQuery.value.toLowerCase()
    const matchBarang = item.barang?.nama_barang?.toLowerCase().includes(q) || item.barang?.kode_barang?.toLowerCase().includes(q)
    const matchKet = item.keterangan?.toLowerCase().includes(q)
    const matchKontak = item.kontak?.nama_kontak?.toLowerCase().includes(q)
    return matchBarang || matchKet || matchKontak
  })
})

const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchLogs()
})

defineExpose({
  refresh: fetchLogs
})
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
      <div>
        <h2 class="text-base font-semibold text-slate-900">Riwayat Audit Mutasi Stok (Party Ledger)</h2>
        <p class="text-xs text-slate-500">Track record terpadu mutasi barang masuk (Supplier) & keluar (Customer) tanpa data NULL</p>
      </div>

      <!-- Quick Filter Buttons -->
      <div class="flex items-center gap-1.5 bg-slate-100 p-1 rounded-lg">
        <button
          @click="filterJenis = ''; fetchLogs()"
          :class="[
            'px-3 py-1.5 text-xs font-medium rounded-md transition',
            filterJenis === '' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'
          ]"
        >
          Semua
        </button>
        <button
          @click="filterJenis = 'IN'; fetchLogs()"
          :class="[
            'px-3 py-1.5 text-xs font-medium rounded-md transition',
            filterJenis === 'IN' ? 'bg-white text-emerald-700 shadow-sm font-semibold' : 'text-slate-600 hover:text-slate-900'
          ]"
        >
          Masuk (Supplier)
        </button>
        <button
          @click="filterJenis = 'OUT'; fetchLogs()"
          :class="[
            'px-3 py-1.5 text-xs font-medium rounded-md transition',
            filterJenis === 'OUT' ? 'bg-white text-rose-700 shadow-sm font-semibold' : 'text-slate-600 hover:text-slate-900'
          ]"
        >
          Keluar (Customer)
        </button>
      </div>
    </div>

    <!-- Search Bar -->
    <div class="flex flex-col sm:flex-row items-center gap-3 bg-white p-3 rounded-xl border border-slate-200 shadow-sm">
      <div class="relative w-full sm:w-80">
        <Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
        <input
          type="text"
          v-model="searchQuery"
          placeholder="Cari produk, customer, PT supplier, atau nota..."
          class="w-full pl-9 pr-3 py-1.5 text-xs rounded-lg border border-slate-300 text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-1 focus:ring-emerald-600"
        />
      </div>

      <div class="ml-auto text-xs text-slate-500 font-medium">
        Total {{ filteredLogs.length }} catatan mutasi
      </div>
    </div>

    <!-- Audit Log Table -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold">
              <th class="py-3 px-4">Waktu Transaksi</th>
              <th class="py-3 px-4">Tipe Mutasi</th>
              <th class="py-3 px-4">Produk Pertanian</th>
              <th class="py-3 px-4 text-center">Jumlah (Qty)</th>
              <th class="py-3 px-4">Pihak Kontak Terkait</th>
              <th class="py-3 px-4">Keterangan / Referensi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="6" class="py-12 text-center text-slate-400">Memuat riwayat mutasi...</td>
            </tr>
            <tr v-else-if="filteredLogs.length === 0">
              <td colspan="6" class="py-12 text-center text-slate-400">Belum ada transaksi mutasi.</td>
            </tr>
            <tr v-else v-for="log in filteredLogs" :key="log.id" class="hover:bg-slate-50/70 transition">
              <td class="py-3 px-4 text-slate-600 font-mono text-[11px]">
                <div class="flex items-center gap-1.5">
                  <Calendar class="w-3.5 h-3.5 text-slate-400" />
                  <span>{{ formatDateTime(log.tanggal_transaksi) }}</span>
                </div>
              </td>
              <td class="py-3 px-4">
                <span
                  :class="[
                    'inline-flex items-center gap-1 px-2.5 py-1 rounded text-[11px] font-bold',
                    log.jenis_transaksi === 'IN'
                      ? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                      : 'bg-rose-50 text-rose-700 border border-rose-200'
                  ]"
                >
                  <ArrowUpRight v-if="log.jenis_transaksi === 'IN'" class="w-3.5 h-3.5" />
                  <ArrowDownRight v-else class="w-3.5 h-3.5" />
                  <span>{{ log.jenis_transaksi === 'IN' ? 'STOK MASUK' : 'STOK KELUAR' }}</span>
                </span>
              </td>
              <td class="py-3 px-4">
                <div class="font-semibold text-slate-900">{{ log.barang?.nama_barang || 'Barang #' + log.barang_id }}</div>
                <div class="text-[10px] text-slate-400">{{ log.barang?.kode_barang }} • {{ log.barang?.kategori }}</div>
              </td>
              <td class="py-3 px-4 text-center">
                <span
                  :class="[
                    'font-bold text-sm',
                    log.jenis_transaksi === 'IN' ? 'text-emerald-700' : 'text-rose-700'
                  ]"
                >
                  {{ log.jenis_transaksi === 'IN' ? '+' : '-' }}{{ log.qty }}
                </span>
                <span class="text-[10px] text-slate-500 ml-1">pcs</span>
              </td>

              <!-- Kolom Kontak dengan Pembeda PT/CV vs Individu Petani -->
              <td class="py-3 px-4">
                <div v-if="log.kontak" class="flex flex-col gap-0.5">
                  <div class="flex items-center gap-1.5">
                    <Building2 v-if="log.kontak.jenis_entitas === 'BADAN_USAHA'" class="w-3.5 h-3.5 text-blue-600 shrink-0" />
                    <User v-else class="w-3.5 h-3.5 text-emerald-600 shrink-0" />
                    <span class="font-semibold text-slate-900">{{ log.kontak.nama_kontak }}</span>
                  </div>
                  <div class="flex items-center gap-1">
                    <span
                      :class="log.kontak.tipe === 'SUPPLIER' ? 'text-blue-700 bg-blue-50' : 'text-emerald-700 bg-emerald-50'"
                      class="text-[9px] px-1 rounded font-bold"
                    >
                      {{ log.kontak.tipe }}
                    </span>
                    <span
                      :class="log.kontak.jenis_entitas === 'BADAN_USAHA' ? 'text-indigo-700 bg-indigo-50' : 'text-slate-600 bg-slate-100'"
                      class="text-[9px] px-1 rounded font-medium"
                    >
                      {{ log.kontak.jenis_entitas === 'BADAN_USAHA' ? '🏢 PT/CV' : '👤 Petani' }}
                    </span>
                  </div>
                </div>
                <span v-else class="text-slate-400">-</span>
              </td>

              <td class="py-3 px-4 text-slate-600">
                {{ log.keterangan || '-' }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
