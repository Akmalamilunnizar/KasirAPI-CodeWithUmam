<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { kontakApi } from '../services/api'
import { Plus, Edit2, Trash2, X, AlertCircle, Phone, MapPin, Building2, User, Filter, Search } from 'lucide-vue-next'

const loading = ref(false)
const kontaks = ref([])
const activeTipe = ref('')
const searchQuery = ref('')
const isModalOpen = ref(false)
const isEditing = ref(false)
const modalError = ref('')
const modalLoading = ref(false)

const form = reactive({
  id: null,
  tipe: 'CUSTOMER',
  jenis_entitas: 'INDIVIDU',
  kode_kontak: '',
  nama_kontak: '',
  no_telp: '',
  alamat: ''
})

const fetchKontaks = async () => {
  loading.value = true
  try {
    const params = {}
    if (activeTipe.value) params.tipe = activeTipe.value
    if (searchQuery.value) params.search = searchQuery.value

    const res = await kontakApi.getAll(params)
    kontaks.value = res.data.data || []
  } catch (err) {
    console.error('Failed to fetch kontaks', err)
  } finally {
    loading.value = false
  }
}

const openCreateModal = (defaultTipe = 'CUSTOMER') => {
  isEditing.value = false
  modalError.value = ''
  form.id = null
  form.tipe = defaultTipe
  form.jenis_entitas = defaultTipe === 'SUPPLIER' ? 'BADAN_USAHA' : 'INDIVIDU'
  form.kode_kontak = `${defaultTipe === 'SUPPLIER' ? 'SPL' : 'CST'}-${Math.floor(100 + Math.random() * 900)}`
  form.nama_kontak = ''
  form.no_telp = ''
  form.alamat = ''
  isModalOpen.value = true
}

const openEditModal = (item) => {
  isEditing.value = true
  modalError.value = ''
  form.id = item.id
  form.tipe = item.tipe
  form.jenis_entitas = item.jenis_entitas
  form.kode_kontak = item.kode_kontak
  form.nama_kontak = item.nama_kontak
  form.no_telp = item.no_telp
  form.alamat = item.alamat
  isModalOpen.value = true
}

const saveKontak = async () => {
  if (!form.nama_kontak.trim()) {
    modalError.value = 'Nama kontak wajib diisi'
    return
  }

  modalLoading.value = true
  modalError.value = ''
  try {
    const payload = {
      tipe: form.tipe,
      jenis_entitas: form.jenis_entitas,
      kode_kontak: form.kode_kontak,
      nama_kontak: form.nama_kontak,
      no_telp: form.no_telp,
      alamat: form.alamat
    }

    if (isEditing.value) {
      await kontakApi.update(form.id, payload)
    } else {
      await kontakApi.create(payload)
    }

    isModalOpen.value = false
    fetchKontaks()
  } catch (err) {
    modalError.value = err.response?.data?.error || 'Gagal menyimpan kontak'
  } finally {
    modalLoading.value = false
  }
}

const deleteKontak = async (id, nama) => {
  if (!confirm(`Hapus kontak "${nama}"?`)) return
  try {
    await kontakApi.delete(id)
    fetchKontaks()
  } catch (err) {
    alert(err.response?.data?.error || 'Gagal menghapus kontak')
  }
}

onMounted(() => {
  fetchKontaks()
})
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
      <div>
        <h2 class="text-base font-semibold text-slate-900">Master Data Kontak (Suppliers & Customers)</h2>
        <p class="text-xs text-slate-500">Kelola distributor/pabrik (Supplier) dan petani/sub-agen (Customer) dengan pembeda entitas</p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="openCreateModal('SUPPLIER')"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-lg transition"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>+ Supplier</span>
        </button>
        <button
          @click="openCreateModal('CUSTOMER')"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-white bg-emerald-600 hover:bg-emerald-700 rounded-lg shadow-sm transition"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>+ Customer / Petani</span>
        </button>
      </div>
    </div>

    <!-- Filter Tab & Search Bar -->
    <div class="flex flex-col sm:flex-row items-center gap-3 bg-white p-3 rounded-xl border border-slate-200 shadow-sm">
      <div class="flex items-center gap-1 bg-slate-100 p-1 rounded-lg">
        <button
          @click="activeTipe = ''; fetchKontaks()"
          :class="[
            'px-3 py-1 text-xs font-medium rounded-md transition',
            activeTipe === '' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'
          ]"
        >
          Semua ({{ kontaks.length }})
        </button>
        <button
          @click="activeTipe = 'SUPPLIER'; fetchKontaks()"
          :class="[
            'px-3 py-1 text-xs font-medium rounded-md transition',
            activeTipe === 'SUPPLIER' ? 'bg-white text-blue-700 shadow-sm' : 'text-slate-600 hover:text-slate-900'
          ]"
        >
          Suppliers (Distributor)
        </button>
        <button
          @click="activeTipe = 'CUSTOMER'; fetchKontaks()"
          :class="[
            'px-3 py-1 text-xs font-medium rounded-md transition',
            activeTipe === 'CUSTOMER' ? 'bg-white text-emerald-700 shadow-sm' : 'text-slate-600 hover:text-slate-900'
          ]"
        >
          Customers (Petani/Agen)
        </button>
      </div>

      <div class="relative w-full sm:w-72 sm:ml-auto">
        <Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
        <input
          type="text"
          v-model="searchQuery"
          @input="fetchKontaks"
          placeholder="Cari kode, nama, atau kota..."
          class="w-full pl-9 pr-3 py-1.5 text-xs rounded-lg border border-slate-300 text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-1 focus:ring-emerald-600"
        />
      </div>
    </div>

    <!-- Data Table -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-slate-600 font-semibold">
              <th class="py-3 px-4">Kode</th>
              <th class="py-3 px-4">Nama Kontak</th>
              <th class="py-3 px-4">Tipe & Entitas</th>
              <th class="py-3 px-4">Kontak / Telepon</th>
              <th class="py-3 px-4">Alamat</th>
              <th class="py-3 px-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-if="loading">
              <td colspan="6" class="py-12 text-center text-slate-400">Memuat data kontak...</td>
            </tr>
            <tr v-else-if="kontaks.length === 0">
              <td colspan="6" class="py-12 text-center text-slate-400">Belum ada data kontak.</td>
            </tr>
            <tr v-else v-for="item in kontaks" :key="item.id" class="hover:bg-slate-50/70 transition">
              <td class="py-3 px-4 font-mono font-medium text-slate-700">{{ item.kode_kontak }}</td>
              <td class="py-3 px-4 font-semibold text-slate-900">
                <div class="flex items-center gap-2">
                  <Building2 v-if="item.jenis_entitas === 'BADAN_USAHA'" class="w-4 h-4 text-blue-600 shrink-0" />
                  <User v-else class="w-4 h-4 text-emerald-600 shrink-0" />
                  <span>{{ item.nama_kontak }}</span>
                </div>
              </td>
              <td class="py-3 px-4">
                <div class="flex items-center gap-1.5">
                  <span
                    :class="item.tipe === 'SUPPLIER' ? 'bg-blue-50 text-blue-700' : 'bg-emerald-50 text-emerald-700'"
                    class="px-2 py-0.5 rounded text-[10px] font-bold"
                  >
                    {{ item.tipe }}
                  </span>
                  <span
                    :class="item.jenis_entitas === 'BADAN_USAHA' ? 'bg-indigo-50 text-indigo-700 border border-indigo-200' : 'bg-slate-100 text-slate-600'"
                    class="px-1.5 py-0.5 rounded text-[10px] font-medium"
                  >
                    {{ item.jenis_entitas === 'BADAN_USAHA' ? '🏢 PT/CV' : '👤 Individu' }}
                  </span>
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
                    title="Edit Kontak"
                  >
                    <Edit2 class="w-3.5 h-3.5" />
                  </button>
                  <button
                    @click="deleteKontak(item.id, item.nama_kontak)"
                    class="p-1.5 text-slate-600 hover:text-rose-700 hover:bg-rose-50 rounded"
                    title="Hapus Kontak"
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
            {{ isEditing ? 'Edit Kontak' : 'Tambah Kontak Baru' }}
          </h3>
          <button @click="isModalOpen = false" class="p-1 text-slate-400 hover:text-slate-600 rounded-lg">
            <X class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveKontak" class="p-6 space-y-3.5">
          <div v-if="modalError" class="p-3 bg-rose-50 border border-rose-200 text-rose-700 text-xs rounded-lg flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{{ modalError }}</span>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-medium text-slate-700 mb-1">Peran Kontak</label>
              <select
                v-model="form.tipe"
                class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 bg-white text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
              >
                <option value="CUSTOMER">CUSTOMER (Pembeli)</option>
                <option value="SUPPLIER">SUPPLIER (Distributor)</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-medium text-slate-700 mb-1">Bentuk Entitas</label>
              <select
                v-model="form.jenis_entitas"
                class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 bg-white text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
              >
                <option value="INDIVIDU">👤 Petani / Individu</option>
                <option value="BADAN_USAHA">🏢 PT / CV / Koperasi</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-700 mb-1">Kode Kontak</label>
            <input
              type="text"
              v-model="form.kode_kontak"
              required
              class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-700 mb-1">Nama Kontak / Instansi</label>
            <input
              type="text"
              v-model="form.nama_kontak"
              required
              placeholder="Contoh: PT Syngenta / Pak Slamet"
              class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-700 mb-1">Nomor Telepon</label>
            <input
              type="text"
              v-model="form.no_telp"
              placeholder="0812-xxxxxxx"
              class="w-full text-xs rounded-lg border border-slate-300 px-3 py-2 text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-600"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-700 mb-1">Alamat Lengkap</label>
            <textarea
              v-model="form.alamat"
              rows="2"
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
              {{ modalLoading ? 'Menyimpan...' : 'Simpan Kontak' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
