<script setup>
import { ref, reactive, watch, onMounted, computed } from 'vue'
import { barangApi, kontakApi, transaksiApi } from '../services/api'
import { ArrowDownRight, ArrowUpRight, AlertCircle, CheckCircle2, X, Building2, User } from 'lucide-vue-next'

const props = defineProps({
  isOpen: Boolean,
  preselectedBarangId: Number
})

const emit = defineEmits(['close', 'success'])

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const barangs = ref([])
const kontaks = ref([])

const form = reactive({
  barang_id: '',
  jenis_transaksi: 'OUT',
  qty: 1,
  kontak_id: '',
  expired_date: '',
  keterangan: ''
})

const selectedBarang = computed(() => {
  return barangs.value.find(b => b.id === Number(form.barang_id))
})

const filteredKontaks = computed(() => {
  const targetTipe = form.jenis_transaksi === 'IN' ? 'SUPPLIER' : 'CUSTOMER'
  return kontaks.value.filter(k => k.tipe === targetTipe)
})

const boxCalculation = computed(() => {
  if (!selectedBarang.value || !form.qty) return null
  const pcsPerBox = selectedBarang.value.pcs_per_box || 12
  const boxes = Math.floor(form.qty / pcsPerBox)
  const remainingPcs = form.qty % pcsPerBox
  return { boxes, remainingPcs, pcsPerBox }
})

const fetchDependencies = async () => {
  try {
    const [resBarang, resKontak] = await Promise.all([
      barangApi.getAll(),
      kontakApi.getAll()
    ])
    barangs.value = resBarang.data.data || []
    kontaks.value = resKontak.data.data || []
  } catch (err) {
    console.error('Failed to fetch dependencies', err)
  }
}

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    errorMessage.value = ''
    successMessage.value = ''
    form.jenis_transaksi = 'OUT'
    form.qty = 1
    form.kontak_id = ''
    form.expired_date = ''
    form.keterangan = ''
    form.barang_id = props.preselectedBarangId || (barangs.value[0]?.id || '')
    fetchDependencies()
  }
})

watch(() => form.jenis_transaksi, () => {
  form.kontak_id = ''
})

const handleSubmit = async () => {
  errorMessage.value = ''
  successMessage.value = ''

  if (!form.barang_id) {
    errorMessage.value = 'Silakan pilih produk pertanian'
    return
  }
  if (!form.kontak_id) {
    errorMessage.value = form.jenis_transaksi === 'IN' ? 'Silakan pilih supplier' : 'Silakan pilih customer/petani pembeli'
    return
  }
  if (!form.qty || form.qty <= 0) {
    errorMessage.value = 'Jumlah qty harus lebih dari 0'
    return
  }
  if (form.jenis_transaksi === 'IN' && !form.expired_date) {
    errorMessage.value = 'Expired Date wajib diisi untuk stok masuk'
    return
  }

  if (form.jenis_transaksi === 'OUT' && selectedBarang.value) {
    if (form.qty > selectedBarang.value.total_stok) {
      errorMessage.value = `Stok fisik tidak mencukupi! Stok saat ini: ${selectedBarang.value.total_stok} pcs`
      return
    }
  }

  loading.value = true
  try {
    const payload = {
      barang_id: Number(form.barang_id),
      kontak_id: Number(form.kontak_id),
      jenis_transaksi: form.jenis_transaksi,
      qty: Number(form.qty),
      expired_date: form.jenis_transaksi === 'IN' ? form.expired_date : '',
      keterangan: form.keterangan
    }

    const res = await transaksiApi.mutasi(payload)
    successMessage.value = res.data.message || 'Mutasi stok berhasil'
    setTimeout(() => {
      emit('success')
      emit('close')
    }, 1000)
  } catch (err) {
    errorMessage.value = err.response?.data?.message || err.response?.data?.error || 'Gagal memproses transaksi'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDependencies()
})
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
    <div class="relative w-full max-w-lg bg-white rounded-xl shadow-2xl border border-slate-200 overflow-hidden">
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-slate-100 bg-slate-50/50">
        <div>
          <h3 class="text-base font-semibold text-slate-900">Transaksi Mutasi Stok</h3>
          <p class="text-xs text-slate-500 mt-0.5">Catat barang masuk (Supplier) atau keluar (Customer/Petani)</p>
        </div>
        <button @click="$emit('close')" class="p-1 rounded-lg text-slate-400 hover:text-slate-600 hover:bg-slate-100">
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Form Body -->
      <form @submit.prevent="handleSubmit" class="p-6 space-y-4">
        <div v-if="errorMessage" class="flex items-start gap-2.5 p-3 rounded-lg bg-rose-50 border border-rose-200 text-rose-700 text-xs">
          <AlertCircle class="w-4 h-4 mt-0.5 shrink-0" />
          <span>{{ errorMessage }}</span>
        </div>
        <div v-if="successMessage" class="flex items-start gap-2.5 p-3 rounded-lg bg-emerald-50 border border-emerald-200 text-emerald-700 text-xs">
          <CheckCircle2 class="w-4 h-4 mt-0.5 shrink-0" />
          <span>{{ successMessage }}</span>
        </div>

        <!-- Switcher IN / OUT -->
        <div>
          <label class="block text-xs font-medium text-slate-700 mb-1.5">Jenis Transaksi</label>
          <div class="grid grid-cols-2 gap-2">
            <button
              type="button"
              @click="form.jenis_transaksi = 'OUT'"
              :class="[
                'flex items-center justify-center gap-2 py-2.5 px-3 rounded-lg text-xs font-medium border transition-all',
                form.jenis_transaksi === 'OUT'
                  ? 'bg-rose-50 border-rose-300 text-rose-700 shadow-sm font-semibold'
                  : 'bg-white border-slate-200 text-slate-600 hover:bg-slate-50'
              ]"
            >
              <ArrowDownRight class="w-4 h-4 text-rose-600" />
              <span>Stok Keluar / Jual (OUT)</span>
            </button>
            <button
              type="button"
              @click="form.jenis_transaksi = 'IN'"
              :class="[
                'flex items-center justify-center gap-2 py-2.5 px-3 rounded-lg text-xs font-medium border transition-all',
                form.jenis_transaksi === 'IN'
                  ? 'bg-emerald-50 border-emerald-300 text-emerald-700 shadow-sm font-semibold'
                  : 'bg-white border-slate-200 text-slate-600 hover:bg-slate-50'
              ]"
            >
              <ArrowUpRight class="w-4 h-4 text-emerald-600" />
              <span>Stok Masuk / PO (IN)</span>
            </button>
          </div>
        </div>

        <!-- Pilih Barang -->
        <div>
          <label class="block text-xs font-medium text-slate-700 mb-1">Pilih Produk Pertanian</label>
          <select
            v-model="form.barang_id"
            class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
          >
            <option value="" disabled>-- Pilih Barang --</option>
            <option v-for="b in barangs" :key="b.id" :value="b.id">
              {{ b.kode_barang }} - {{ b.nama_barang }} (Stok: {{ b.total_stok }} pcs)
            </option>
          </select>
        </div>

        <!-- Pihak Terkait (Supplier jika IN, Customer jika OUT) -->
        <div>
          <label class="block text-xs font-medium text-slate-700 mb-1">
            {{ form.jenis_transaksi === 'IN' ? 'Supplier Pengirim' : 'Customer / Petani Pembeli' }}
            <span class="text-rose-500">*</span>
          </label>
          <select
            v-model="form.kontak_id"
            required
            class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 bg-white text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
          >
            <option value="" disabled>
              -- Pilih {{ form.jenis_transaksi === 'IN' ? 'Supplier' : 'Customer/Petani' }} --
            </option>
            <option v-for="k in filteredKontaks" :key="k.id" :value="k.id">
              [{{ k.jenis_entitas === 'BADAN_USAHA' ? 'PT/CV' : 'Petani' }}] {{ k.nama_kontak }} ({{ k.kode_kontak }})
            </option>
          </select>
        </div>

        <!-- Field Expired Date (Khusus IN) -->
        <div v-if="form.jenis_transaksi === 'IN'" class="p-3 bg-emerald-50/50 rounded-lg border border-emerald-100">
          <label class="block text-xs font-medium text-slate-700 mb-1">
            Expired Date Batch <span class="text-rose-500">*</span>
          </label>
          <input
            type="date"
            v-model="form.expired_date"
            required
            class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 bg-white text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
          />
        </div>

        <!-- Live Stock & FEFO Info -->
        <div v-if="selectedBarang && form.jenis_transaksi === 'OUT'" class="p-3 rounded-lg bg-slate-50 border border-slate-200 text-xs space-y-1">
          <div class="flex justify-between text-slate-600">
            <span>Stok Fisik Aktif:</span>
            <span class="font-semibold text-slate-900">{{ selectedBarang.total_stok }} pcs ({{ selectedBarang.total_box }} box + {{ selectedBarang.sisa_pcs }} pcs)</span>
          </div>
          <div class="text-[11px] text-amber-700 pt-1 border-t border-slate-200">
            ⚡ <b>FEFO Otomatis:</b> Batch tanggal kadaluarsa terdekat dipotong lebih dulu.
          </div>
        </div>

        <!-- Qty Input -->
        <div>
          <div class="flex justify-between items-center mb-1">
            <label class="text-xs font-medium text-slate-700">Jumlah (Qty Pcs)</label>
            <span v-if="boxCalculation" class="text-[11px] text-slate-500">
              ≈ {{ boxCalculation.boxes }} Box + {{ boxCalculation.remainingPcs }} Pcs (@{{ boxCalculation.pcsPerBox }}/box)
            </span>
          </div>
          <input
            type="number"
            min="1"
            v-model.number="form.qty"
            required
            class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
          />
        </div>

        <!-- Keterangan -->
        <div>
          <label class="block text-xs font-medium text-slate-700 mb-1">Keterangan Transaksi</label>
          <input
            type="text"
            v-model="form.keterangan"
            class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
            placeholder="Contoh: Nota Penjualan / Faktur PO"
          />
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-100">
          <button
            type="button"
            @click="$emit('close')"
            class="px-4 py-2 text-xs font-medium text-slate-600 hover:text-slate-800 hover:bg-slate-100 rounded-lg transition"
          >
            Batal
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="px-4 py-2 text-xs font-medium text-white bg-emerald-600 hover:bg-emerald-700 disabled:opacity-50 rounded-lg shadow-sm transition"
          >
            {{ loading ? 'Memproses...' : 'Simpan Transaksi' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
