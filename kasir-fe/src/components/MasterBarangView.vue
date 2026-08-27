<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { barangApi } from '../services/api'
import {
  Package,
  Plus,
  Search,
  Filter,
  Edit2,
  Trash2,
  AlertCircle,
  X,
  Boxes,
  ArrowUpDown
} from 'lucide-vue-next'

const emit = defineEmits(['openMutasiFor'])

const loading = ref(false)
const barangs = ref([])
const searchQuery = ref('')
const selectedKategori = ref('')
const isModalOpen = ref(false)
const isEditing = ref(false)
const modalError = ref('')
const modalLoading = ref(false)

const kategoriList = [
  'BENIH_BIBIT',
  'PUPUK_PADAT',
  'PUPUK_CAIR',
  'PESTISIDA',
  'ALAT_PERTANIAN'
]

const form = reactive({
  id: null,
  kode_barang: '',
  nama_barang: '',
  kategori: 'BENIH_BIBIT',
  min_stok: 20,
  pcs_per_box: 12
})

const fetchBarangs = async () => {
  loading.value = true
  try {
    const params = {}
    if (selectedKategori.value) params.kategori = selectedKategori.value
    if (searchQuery.value) params.search = searchQuery.value

    const res = await barangApi.getAll(params)
    barangs.value = res.data.data || []
  } catch (err) {
    console.error('Failed to fetch barangs', err)
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  isEditing.value = false
  modalError.value = ''
  form.id = null
  form.kode_barang = `AGR-${Math.floor(1000 + Math.random() * 9000)}`
  form.nama_barang = ''
  form.kategori = 'BENIH_BIBIT'
  form.min_stok = 20
  form.pcs_per_box = 12
  isModalOpen.value = true
}

const openEditModal = (item) => {
  isEditing.value = true
  modalError.value = ''
  form.id = item.id
  form.kode_barang = item.kode_barang
  form.nama_barang = item.nama_barang
  form.kategori = item.kategori
  form.min_stok = item.min_stok
  form.pcs_per_box = item.pcs_per_box
  isModalOpen.value = true
}

const saveBarang = async () => {
  if (!form.nama_barang.trim()) {
    modalError.value = 'Nama barang wajib diisi'
    return
  }

  modalLoading.value = true
  modalError.value = ''
  try {
    const payload = {
      kode_barang: form.kode_barang,
      nama_barang: form.nama_barang,
      kategori: form.kategori,
      min_stok: Number(form.min_stok),
      pcs_per_box: Number(form.pcs_per_box)
    }

    if (isEditing.value) {
      await barangApi.update(form.id, payload)
    } else {
      await barangApi.create(payload)
    }

    isModalOpen.value = false
    fetchBarangs()
  } catch (err) {
    modalError.value = err.response?.data?.error || 'Gagal menyimpan data barang'
  } finally {
    modalLoading.value = false
  }
}

const deleteBarang = async (id, nama) => {
  if (!confirm(`Hapus barang "${nama}"? (Data akan di-soft delete)`)) return
  try {
    await barangApi.delete(id)
    fetchBarangs()
  } catch (err) {
    alert(err.response?.data?.error || 'Gagal menghapus barang')
  }
}

onMounted(() => {
  fetchBarangs()
})

defineExpose({
  refresh: fetchBarangs
})
</script>

<template>
  <div class="space-y-4">
    <!-- Header & Controls -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
      <div>
        <h2 class="text-base font-semibold text-slate-900">Master Data Produk Pertanian</h2>
        <p class="text-xs text-slate-500">Kelola katalog barang, kategori, ambang batas stok, dan rasio box</p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="openCreateModal"
          class="inline-flex items-center gap-2 px-3.5 py-2 text-xs font-medium text-white bg-emerald-600 hover:bg-emerald-700 rounded-lg shadow-sm transition"
        >
          <Plus class="w-4 h-4" />
          <span>Tambah Produk</span>
        </button>
      </div>
    </div>

    <!-- Filter & Search Bar -->
    <div class="flex flex-col sm:flex-row items-center gap-3 bg-white p-3 rounded-xl border border-slate-200 shadow-sm">
      <div class="relative w-full sm:w-80">
        <Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
        <input
          type="text"
          v-model="searchQuery"
          @input="fetchBarangs"
          placeholder="Cari kode atau nama barang..."
          class="w-full pl-9 pr-3 py-1.5 text-xs rounded-lg border border-slate-300 text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-1 focus:ring-emerald-600"
        />
      </div>

      <div class="w-full sm:w-56">
        <select
          v-model="selectedKategori"
          @change="fetchBarangs"
          class="w-full py-1.5 px-3 text-xs rounded-lg border border-slate-300 bg-white text-slate-800 focus:outline-none focus:ring-1 focus:ring-emerald-600"
        >
          <option value="">Semua Kategori</option>
          <option v-for="kat in kategoriList" :key="kat" :value="kat">{{ kat }}</option>
        </select>
      </div>

      <div class="ml-auto text-xs text-slate-500 font-medium">
        Total: {{ barangs.length }} item
      </div>
    </div>

    <!-- Data Table -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold">
              <th class="py-3 px-4">Kode</th>
              <th class="py-3 px-4">Nama Produk</th>
              <th class="py-3 px-4">Kategori</th>
              <th class="py-3 px-4 text-center">Rasio Box</th>
              <th class="py-3 px-4 text-center">Min. Stok</th>
              <th class="py-3 px-4 text-center">Total Stok Aktif</th>
              <th class="py-3 px-4 text-center">Status</th>
              <th class="py-3 px-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="8" class="py-12 text-center text-slate-400">Memuat katalog barang...</td>
            </tr>
            <tr v-else-if="barangs.length === 0">
              <td colspan="8" class="py-12 text-center text-slate-400">Tidak ada barang yang sesuai filter.</td>
            </tr>
            <tr v-else v-for="item in barangs" :key="item.id" class="hover:bg-slate-50/70 transition">
              <td class="py-3 px-4 font-mono font-medium text-slate-700">{{ item.kode_barang }}</td>
              <td class="py-3 px-4">
                <div class="font-semibold text-slate-900">{{ item.nama_barang }}</div>
                <div class="text-[10px] text-slate-400">ID: #{{ item.id }}</div>
              </td>
              <td class="py-3 px-4">
                <span class="inline-block px-2 py-0.5 rounded bg-slate-100 text-slate-700 text-[11px] font-medium">
                  {{ item.kategori }}
                </span>
              </td>
              <td class="py-3 px-4 text-center text-slate-600">
                1 Box = {{ item.pcs_per_box }} Pcs
              </td>
              <td class="py-3 px-4 text-center font-medium text-slate-700">
                {{ item.min_stok }} pcs
              </td>
              <td class="py-3 px-4 text-center">
                <div class="font-bold text-slate-900">{{ item.total_stok }} pcs</div>
                <div class="text-[10px] text-slate-500">({{ item.total_box }} box + {{ item.sisa_pcs }} pcs)</div>
              </td>
              <td class="py-3 px-4 text-center">
                <span
                  :class="[
                    'inline-block px-2 py-0.5 rounded text-[10px] font-medium',
                    item.stok_status === 'OUT_OF_STOCK'
                      ? 'bg-rose-100 text-rose-700'
                      : item.stok_status === 'CRITICAL'
                      ? 'bg-amber-100 text-amber-800'
                      : item.stok_status === 'WARNING'
                      ? 'bg-yellow-50 text-yellow-700 border border-yellow-200'
                      : 'bg-emerald-50 text-emerald-700'
                  ]"
                >
                  {{ item.stok_status }}
                </span>
              </td>
              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <button
                    @click="$emit('openMutasiFor', item.id)"
                    class="p-1.5 text-slate-600 hover:text-emerald-700 hover:bg-emerald-50 rounded"
                    title="Catat Mutasi"
                  >
                    <ArrowUpDown class="w-3.5 h-3.5" />
                  </button>
                  <button
                    @click="openEditModal(item)"
                    class="p-1.5 text-slate-600 hover:text-blue-700 hover:bg-blue-50 rounded"
                    title="Edit Barang"
                  >
                    <Edit2 class="w-3.5 h-3.5" />
                  </button>
                  <button
                    @click="deleteBarang(item.id, item.nama_barang)"
                    class="p-1.5 text-slate-600 hover:text-rose-700 hover:bg-rose-50 rounded"
                    title="Hapus Barang"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create / Edit Modal -->
    <div v-if="isModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
      <div class="relative w-full max-w-md bg-white rounded-xl shadow-2xl border border-slate-200 overflow-hidden">
        <div class="flex items-center justify-between px-6 py-4 border-b border-slate-100 bg-slate-50">
          <h3 class="text-sm font-semibold text-slate-900">
            {{ isEditing ? 'Edit Produk Pertanian' : 'Tambah Produk Baru' }}
          </h3>
          <button @click="isModalOpen = false" class="p-1 text-slate-400 hover:text-slate-600 rounded-lg">
            <X class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveBarang" class="p-6 space-y-3.5">
          <div v-if="modalError" class="p-3 bg-rose-50 border border-rose-200 text-rose-700 text-xs rounded-lg flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{{ modalError }}</span>
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-700 mb-1">Kode Barang</label>
            <input
              type="text"
              v-model="form.kode_barang"
              required
              class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-700 mb-1">Nama Produk Pertanian</label>
            <input
              type="text"
              v-model="form.nama_barang"
              required
              placeholder="Contoh: Benih Jagung Hibrida Pioneer P27"
              class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-700 mb-1">Kategori Produk</label>
            <select
              v-model="form.kategori"
              class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 bg-white text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
            >
              <option v-for="kat in kategoriList" :key="kat" :value="kat">{{ kat }}</option>
            </select>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-medium text-slate-700 mb-1">Min. Stok (Pcs)</label>
              <input
                type="number"
                min="0"
                v-model.number="form.min_stok"
                required
                class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
              />
            </div>
            <div>
              <label class="block text-xs font-medium text-slate-700 mb-1">Pcs per Box</label>
              <input
                type="number"
                min="1"
                v-model.number="form.pcs_per_box"
                required
                class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
              />
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
            <button
              type="button"
              @click="isModalOpen = false"
              class="px-4 py-2 text-xs font-medium text-slate-600 hover:text-slate-800 rounded-lg"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="modalLoading"
              class="px-4 py-2 text-xs font-medium text-white bg-emerald-600 hover:bg-emerald-700 disabled:opacity-50 rounded-lg shadow-sm"
            >
              {{ modalLoading ? 'Menyimpan...' : 'Simpan Produk' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
