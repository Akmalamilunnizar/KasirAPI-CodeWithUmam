<script setup>
import { ref, reactive, onMounted } from 'vue'
import { supplierApi } from '../services/api'
import { Plus, Edit2, Trash2, X, AlertCircle, Phone, MapPin, Building2 } from 'lucide-vue-next'

const loading = ref(false)
const suppliers = ref([])
const isModalOpen = ref(false)
const isEditing = ref(false)
const modalError = ref('')
const modalLoading = ref(false)

const form = reactive({
  id: null,
  kode_supplier: '',
  nama_supplier: '',
  no_telp: '',
  alamat: ''
})

const fetchSuppliers = async () => {
  loading.value = true
  try {
    const res = await supplierApi.getAll()
    suppliers.value = res.data.data || []
  } catch (err) {
    console.error('Failed to fetch suppliers', err)
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  isEditing.value = false
  modalError.value = ''
  form.id = null
  form.kode_supplier = `SPL-${Math.floor(100 + Math.random() * 900)}`
  form.nama_supplier = ''
  form.no_telp = ''
  form.alamat = ''
  isModalOpen.value = true
}

const openEditModal = (item) => {
  isEditing.value = true
  modalError.value = ''
  form.id = item.id
  form.kode_supplier = item.kode_supplier
  form.nama_supplier = item.nama_supplier
  form.no_telp = item.no_telp
  form.alamat = item.alamat
  isModalOpen.value = true
}

const saveSupplier = async () => {
  if (!form.nama_supplier.trim()) {
    modalError.value = 'Nama supplier wajib diisi'
    return
  }

  modalLoading.value = true
  modalError.value = ''
  try {
    const payload = {
      kode_supplier: form.kode_supplier,
      nama_supplier: form.nama_supplier,
      no_telp: form.no_telp,
      alamat: form.alamat
    }

    if (isEditing.value) {
      await supplierApi.update(form.id, payload)
    } else {
      await supplierApi.create(payload)
    }

    isModalOpen.value = false
    fetchSuppliers()
  } catch (err) {
    modalError.value = err.response?.data?.error || 'Gagal menyimpan supplier'
  } finally {
    modalLoading.value = false
  }
}

const deleteSupplier = async (id, nama) => {
  if (!confirm(`Hapus supplier "${nama}"?`)) return
  try {
    await supplierApi.delete(id)
    fetchSuppliers()
  } catch (err) {
    alert(err.response?.data?.error || 'Gagal menghapus supplier')
  }
}

onMounted(() => {
  fetchSuppliers()
})
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
      <div>
        <h2 class="text-base font-semibold text-slate-900">Master Data Supplier & Distributor</h2>
        <p class="text-xs text-slate-500">Daftar produsen pestisida, pabrik pupuk, dan distributor benih</p>
      </div>

      <button
        @click="openCreateModal"
        class="inline-flex items-center gap-2 px-3.5 py-2 text-xs font-medium text-white bg-emerald-600 hover:bg-emerald-700 rounded-lg shadow-sm transition"
      >
        <Plus class="w-4 h-4" />
        <span>Tambah Supplier</span>
      </button>
    </div>

    <!-- Data Table -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold">
              <th class="py-3 px-4">Kode</th>
              <th class="py-3 px-4">Nama Supplier</th>
              <th class="py-3 px-4">Kontak / Telepon</th>
              <th class="py-3 px-4">Alamat Kantor / Gudang</th>
              <th class="py-3 px-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="5" class="py-12 text-center text-slate-400">Memuat data supplier...</td>
            </tr>
            <tr v-else-if="suppliers.length === 0">
              <td colspan="5" class="py-12 text-center text-slate-400">Belum ada data supplier.</td>
            </tr>
            <tr v-else v-for="item in suppliers" :key="item.id" class="hover:bg-slate-50/70 transition">
              <td class="py-3 px-4 font-mono font-medium text-slate-700">{{ item.kode_supplier }}</td>
              <td class="py-3 px-4 font-semibold text-slate-900">
                <div class="flex items-center gap-2">
                  <Building2 class="w-4 h-4 text-slate-400" />
                  <span>{{ item.nama_supplier }}</span>
                </div>
              </td>
              <td class="py-3 px-4 text-slate-600">
                <div class="flex items-center gap-1.5">
                  <Phone class="w-3.5 h-3.5 text-slate-400" />
                  <span>{{ item.no_telp || '-' }}</span>
                </div>
              </td>
              <td class="py-3 px-4 text-slate-600">
                <div class="flex items-center gap-1.5 truncate max-w-xs">
                  <MapPin class="w-3.5 h-3.5 text-slate-400 shrink-0" />
                  <span class="truncate">{{ item.alamat || '-' }}</span>
                </div>
              </td>
              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <button
                    @click="openEditModal(item)"
                    class="p-1.5 text-slate-600 hover:text-blue-700 hover:bg-blue-50 rounded"
                    title="Edit Supplier"
                  >
                    <Edit2 class="w-3.5 h-3.5" />
                  </button>
                  <button
                    @click="deleteSupplier(item.id, item.nama_supplier)"
                    class="p-1.5 text-slate-600 hover:text-rose-700 hover:bg-rose-50 rounded"
                    title="Hapus Supplier"
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
            {{ isEditing ? 'Edit Supplier' : 'Tambah Supplier Baru' }}
          </h3>
          <button @click="isModalOpen = false" class="p-1 text-slate-400 hover:text-slate-600 rounded-lg">
            <X class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveSupplier" class="p-6 space-y-3.5">
          <div v-if="modalError" class="p-3 bg-rose-50 border border-rose-200 text-rose-700 text-xs rounded-lg flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{{ modalError }}</span>
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-700 mb-1">Kode Supplier</label>
            <input
              type="text"
              v-model="form.kode_supplier"
              required
              class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-700 mb-1">Nama Supplier / PT</label>
            <input
              type="text"
              v-model="form.nama_supplier"
              required
              placeholder="Contoh: PT Syngenta Indonesia"
              class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-700 mb-1">Nomor Telepon / Kontak</label>
            <input
              type="text"
              v-model="form.no_telp"
              placeholder="021-xxxxxxx"
              class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-700 mb-1">Alamat Kantor / Gudang</label>
            <textarea
              v-model="form.alamat"
              rows="2"
              placeholder="Alamat lengkap supplier"
              class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
            ></textarea>
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
              {{ modalLoading ? 'Menyimpan...' : 'Simpan Supplier' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
